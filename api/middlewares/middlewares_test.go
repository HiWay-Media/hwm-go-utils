package middlewares

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const testIssuer = "https://sso.example/realms/test"

func newKey(t *testing.T) (*rsa.PrivateKey, string, string) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	bare := base64.StdEncoding.EncodeToString(der)
	full := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der}))
	return priv, bare, full
}

func validClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"iss":           testIssuer,
		"aud":           []interface{}{"account", "my-api"},
		"exp":           time.Now().Add(time.Hour).Unix(),
		"sub":           "user-uuid",
		"email":         "user@example.com",
		"customer_name": "acme",
		"customer_id":   float64(42),
		"realm_access":  map[string]interface{}{"roles": []interface{}{"viewer"}},
	}
}

func sign(t *testing.T, priv *rsa.PrivateKey, claims jwt.MapClaims) string {
	t.Helper()
	s, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(priv)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func newApp(handlers ...fiber.Handler) *fiber.App {
	app := fiber.New()
	handlers = append(handlers, func(c *fiber.Ctx) error {
		claims, _ := GetTokenClaims(c)
		return c.SendString(GetCurrentLoggedUserEmailFromJwt(claims))
	})
	app.Get("/", handlers...)
	return app
}

func do(t *testing.T, app *fiber.App, authorization string) int {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	return resp.StatusCode
}

func TestJwtProtected(t *testing.T) {
	priv, bare, full := newKey(t)
	other, _, _ := newKey(t)

	noExp := validClaims()
	delete(noExp, "exp")
	expired := validClaims()
	expired["exp"] = time.Now().Add(-time.Minute).Unix()

	// HS256 signed with the public key bytes as HMAC secret (alg confusion)
	hsToken, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, validClaims()).SignedString([]byte(full))
	noneToken, _ := jwt.NewWithClaims(jwt.SigningMethodNone, validClaims()).SignedString(jwt.UnsafeAllowNoneSignatureType)

	strict := []JwtOption{WithIssuer(testIssuer), WithAudience("my-api")}
	wrongIss := validClaims()
	wrongIss["iss"] = "https://evil.example/realms/test"
	wrongAud := validClaims()
	wrongAud["aud"] = "account"

	cases := []struct {
		name   string
		key    string
		opts   []JwtOption
		header string
		want   int
	}{
		{"valid, bare key", bare, nil, "Bearer " + sign(t, priv, validClaims()), 200},
		{"valid, pem key", full, nil, "Bearer " + sign(t, priv, validClaims()), 200},
		{"valid, lowercase bearer", bare, nil, "bearer " + sign(t, priv, validClaims()), 200},
		{"valid, issuer and audience", bare, strict, "Bearer " + sign(t, priv, validClaims()), 200},
		{"no header", bare, nil, "", 401},
		{"no bearer prefix", bare, nil, sign(t, priv, validClaims()), 401},
		{"empty bearer", bare, nil, "Bearer ", 401},
		{"garbage token", bare, nil, "Bearer not.a.jwt", 401},
		{"signed by other key", bare, nil, "Bearer " + sign(t, other, validClaims()), 401},
		{"alg HS256 confusion", bare, nil, "Bearer " + hsToken, 401},
		{"alg none", bare, nil, "Bearer " + noneToken, 401},
		{"expired", bare, nil, "Bearer " + sign(t, priv, expired), 401},
		{"missing exp", bare, nil, "Bearer " + sign(t, priv, noExp), 401},
		{"wrong issuer", bare, strict, "Bearer " + sign(t, priv, wrongIss), 401},
		{"wrong audience", bare, strict, "Bearer " + sign(t, priv, wrongAud), 401},
		{"invalid public key", "not-a-key", nil, "Bearer " + sign(t, priv, validClaims()), 401},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := do(t, newApp(JwtProtected(c.key, c.opts...)), c.header); got != c.want {
				t.Errorf("status %d, want %d", got, c.want)
			}
		})
	}
}

func TestRoleCheck(t *testing.T) {
	priv, bare, _ := newKey(t)
	noRoles := validClaims()
	delete(noRoles, "realm_access")

	cases := []struct {
		name     string
		handlers []fiber.Handler
		claims   jwt.MapClaims
		want     int
	}{
		{"role present", []fiber.Handler{JwtProtected(bare), RoleCheck([]string{"admin", "viewer"})}, validClaims(), 200},
		{"role missing", []fiber.Handler{JwtProtected(bare), RoleCheck([]string{"admin"})}, validClaims(), 403},
		{"no realm_access claim", []fiber.Handler{JwtProtected(bare), RoleCheck([]string{"admin"})}, noRoles, 403},
		{"no roles required", []fiber.Handler{JwtProtected(bare), RoleCheck(nil)}, noRoles, 200},
		{"without JwtProtected", []fiber.Handler{RoleCheck([]string{"viewer"})}, validClaims(), 401},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := do(t, newApp(c.handlers...), "Bearer "+sign(t, priv, c.claims)); got != c.want {
				t.Errorf("status %d, want %d", got, c.want)
			}
		})
	}
}

func TestClaimGettersDoNotPanicOnMissingClaims(t *testing.T) {
	empty := jwt.MapClaims{}
	if GetCustomerNameFromJwt(empty) != "" || GetCustomerIdFromJwt(empty) != 0 ||
		GetCurrentLoggedUserEmailFromJwt(empty) != "" || GetCurrentLoggedUserUUIDFromJwt(empty) != "" ||
		GetClientIdFromJwt(empty) != "_client" {
		t.Error("expected zero values for missing claims")
	}
	if _, err := GetRolesListFromJwt(empty); err == nil {
		t.Error("expected error for missing realm_access")
	}

	c := validClaims()
	if GetCustomerIdFromJwt(c) != 42 || GetClientIdFromJwt(c) != "acme_client" || GetCurrentLoggedUserUUIDFromJwt(c) != "user-uuid" {
		t.Error("unexpected claim values")
	}
}
