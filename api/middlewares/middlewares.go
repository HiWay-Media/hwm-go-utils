package middlewares

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/HiWay-Media/hwm-go-utils/api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

const tokenClaimsKey = "tokenClaims"

type jwtConfig struct {
	issuer   string
	audience string
}

// JwtOption configures optional checks performed by JwtProtected.
type JwtOption func(*jwtConfig)

// WithIssuer requires the token "iss" claim to equal issuer
// (for Keycloak: <server>/realms/<realm>).
func WithIssuer(issuer string) JwtOption {
	return func(c *jwtConfig) { c.issuer = issuer }
}

// WithAudience requires audience to be present in the token "aud" claim.
func WithAudience(audience string) JwtOption {
	return func(c *jwtConfig) { c.audience = audience }
}

// JwtProtected wrap http handler functions for jwt verification.
// publicKey is the realm RSA public key, either bare base64 (as shown by the
// Keycloak console) or a full PEM block.
func JwtProtected(publicKey string, opts ...JwtOption) fiber.Handler {
	cfg := jwtConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}

	key, keyErr := parsePublicKey(publicKey)
	if keyErr != nil {
		log.Printf("JwtProtected: invalid public key, every request will be rejected: %v", keyErr)
	}

	return func(c *fiber.Ctx) error {
		if keyErr != nil {
			return c.Status(http.StatusUnauthorized).JSON(models.ApiDefaultError("unauthorized"))
		}
		tokenString, ok := bearerToken(c.Get(fiber.HeaderAuthorization))
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(models.ApiDefaultError("malformed token"))
		}
		claims, err := verifyJWT(tokenString, key, cfg)
		if err != nil {
			log.Printf("jwt rejected on %s: %v", c.Path(), err)
			return c.Status(http.StatusUnauthorized).JSON(models.ApiDefaultError("invalid token"))
		}
		c.Context().SetUserValue(tokenClaimsKey, claims)
		return c.Next()
	}
}

func bearerToken(header string) (string, bool) {
	const prefix = "Bearer "
	if len(header) <= len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return "", false
	}
	token := strings.TrimSpace(header[len(prefix):])
	return token, token != ""
}

func parsePublicKey(publicKey string) (*rsa.PublicKey, error) {
	pem := strings.TrimSpace(publicKey)
	if !strings.HasPrefix(pem, "-----BEGIN") {
		pem = "-----BEGIN PUBLIC KEY-----\n" + insertNewlines(pem, 64) + "\n-----END PUBLIC KEY-----\n"
	}
	return jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
}

func insertNewlines(input string, every int) string {
	var b strings.Builder
	for i, char := range input {
		b.WriteRune(char)
		if (i+1)%every == 0 && i+1 < len(input) {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// verifyJWT checks signature (RS256/384/512 only), expiry and the optional
// issuer/audience, and returns the token claims.
func verifyJWT(tokenString string, key *rsa.PublicKey, cfg jwtConfig) (jwt.MapClaims, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256", "RS384", "RS512"}))
	claims := jwt.MapClaims{}
	token, err := parser.ParseWithClaims(tokenString, claims, func(*jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid jwt token")
	}
	// MapClaims.Valid only checks exp when present: a token without exp never expires
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		return nil, errors.New("missing or expired exp claim")
	}
	if cfg.issuer != "" && !claims.VerifyIssuer(cfg.issuer, true) {
		return nil, fmt.Errorf("unexpected issuer %v", claims["iss"])
	}
	if cfg.audience != "" && !claims.VerifyAudience(cfg.audience, true) {
		return nil, fmt.Errorf("audience %q not in token", cfg.audience)
	}
	return claims, nil
}

// GetTokenClaims returns the claims stored by JwtProtected.
func GetTokenClaims(c *fiber.Ctx) (jwt.MapClaims, bool) {
	claims, ok := c.Context().UserValue(tokenClaimsKey).(jwt.MapClaims)
	return claims, ok
}

func GetRolesListFromJwt(claims jwt.MapClaims) ([]string, error) {
	//getting roles from jwt claims
	realmaccess, ok := claims["realm_access"].(map[string]interface{})
	if !ok {
		return nil, errors.New("realm_access claim missing")
	}
	realmaccessroles, ok := realmaccess["roles"].([]interface{})
	if !ok {
		return nil, errors.New("realm_access.roles claim missing")
	}

	var roles []string
	for _, v := range realmaccessroles {
		if role, ok := v.(string); ok {
			roles = append(roles, role)
		}
	}

	return roles, nil
}

func CheckRolePresent(roles []string, roleToFind string) bool {
	for _, role := range roles {
		if role == roleToFind {
			return true
		}
	}
	return false
}

// function needed to wrap http handlers for limit the access for specific roles,
// if at least one role specified is present in the jwt, it can pass, otherwise early return 403
// (401 if JwtProtected did not run before it)
func RoleCheck(rolesToCheck []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims, ok := GetTokenClaims(c)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(models.ApiDefaultError("unauthorized"))
		}
		//getting roles list from jwt and match against role needed for the api endpoint
		rolesFromJwt, _ := GetRolesListFromJwt(claims)
		for _, role := range rolesToCheck {
			if CheckRolePresent(rolesFromJwt, role) {
				return c.Next()
			}
		}
		if len(rolesToCheck) == 0 {
			return c.Next()
		}
		return c.Status(http.StatusForbidden).JSON(models.ApiDefaultError(fmt.Sprintf("Missing role %s", strings.Join(rolesToCheck, ", "))))
	}
}

func GetClientIdFromJwt(claims jwt.MapClaims) string {
	return GetCustomerNameFromJwt(claims) + "_client"
}

// get current customer from jwt ("" if the claim is missing)
func GetCustomerNameFromJwt(claims jwt.MapClaims) string {
	s, _ := claims["customer_name"].(string)
	return s
}

// get current customer_id from jwt (0 if the claim is missing)
func GetCustomerIdFromJwt(claims jwt.MapClaims) int {
	f, _ := claims["customer_id"].(float64)
	return int(f)
}

// get current logged in user email (username) from jwt ("" if the claim is missing)
func GetCurrentLoggedUserEmailFromJwt(claims jwt.MapClaims) string {
	s, _ := claims["email"].(string)
	return s
}

// get current logged in user uuid from jwt ("" if the claim is missing)
func GetCurrentLoggedUserUUIDFromJwt(claims jwt.MapClaims) string {
	s, _ := claims["sub"].(string)
	return s
}
