package middlewares

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/HiWay-Media/hwm-go-utils/api/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

//var SECRETRSAPUBLICKEY = []byte(config.CONFIGURATION.KeycloakPublicKey)

// JwtProtected wrap http handler functions for jwt verification
func JwtProtected(publicKey string) fiber.Handler {
	pKey := "-----BEGIN PUBLIC KEY-----\n"+publicKey+"\n-----END PUBLIC KEY-----\n"
	//
	return func(c *fiber.Ctx) error {
		authHeader := strings.Split(c.GetReqHeaders()["Authorization"], "Bearer ")
		if len(authHeader) != 2 {
			log.Println("Malformed token on request: ", c.Request().URI())
			//utils.Bug("Malformed token on request: %s", c.Request().URI())
			return c.Status(http.StatusUnauthorized).JSON(models.ApiDefaultError("malformed token"))
		} else {
			tokenString := authHeader[1]
			// need to fix this metod
			isOk, token, err := verifyJWT_RSA(tokenString, []byte(pKey))
			if err != nil || !isOk {
				return c.Status(http.StatusUnauthorized).JSON(models.ApiDefaultError(fmt.Sprintf("error during verify jwt, err: %s", err.Error())))
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				c.Context().SetUserValue("tokenClaims", claims)
				return c.Next()
			} else {
				return c.Status(http.StatusUnauthorized).JSON(models.ApiDefaultError("Unhautorized, token is not valid"))
			}
			return c.Status(http.StatusUnauthorized).JSON(models.ApiDefaultError("unhautorized"))
		}
	}
}

// Verify a JWT token using an RSA public key
func verifyJWT_RSA(token string, publicKey []byte) (bool, *jwt.Token, error) {

	var parsedToken *jwt.Token

	// parse token
	state, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		// ensure signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unknown signing method")
		}

		parsedToken = token

		// verify
		key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
		if err != nil {
			return nil, fmt.Errorf("AuthKeycloak", err.Error())
		}

		return key, nil
	})

	if err != nil {
		return false, &jwt.Token{}, fmt.Errorf("AuthKeycloak", err.Error())
	}

	if !state.Valid {
		return false, &jwt.Token{}, errors.New("invalid jwt token")
	}

	return true, parsedToken, nil
}

func GetRolesListFromJwt(claims jwt.MapClaims) ([]string, error) {
	//getting roles from jwt claims
	realmaccess := claims["realm_access"].(map[string]interface{})
	realmaccessroles := realmaccess["roles"].([]interface{})

	var roles []string
	for _, v := range realmaccessroles {
		roles = append(roles, v.(string))
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
// if at least one role specified is present in the jwt, it can pass, otherwise early return 401
func RoleCheck(rolesToCheck []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		//getting roles list from jwt and match against role needed for the api endpoint
		rolesFromJwt, _ := GetRolesListFromJwt(c.Context().Value("tokenClaims").(jwt.MapClaims))
		var missingCounter int
		for _, role := range rolesToCheck {
			if !CheckRolePresent(rolesFromJwt, role) {
				missingCounter++
			} // if all requested role missing, error and return
			if missingCounter == len(rolesToCheck) {
				//logrus.Errorf("Missing roles in claims: %s", role)
				return c.Status(http.StatusUnauthorized).JSON(models.ApiDefaultError(fmt.Sprintf("Missing role %s", role)))
			}
		}
		return c.Next()
	}
}

func GetClientIdFromJwt(claims jwt.MapClaims) string {
	return GetCustomerNameFromJwt(claims) + "_client"
}

// get current customer from jwt
func GetCustomerNameFromJwt(claims jwt.MapClaims) string {
	return claims["customer_name"].(string)
}

// get current customer_id from jwt
func GetCustomerIdFromJwt(claims jwt.MapClaims) int {
	//return claims["customer_id"].(int)
	return int(claims["customer_id"].(float64))
}

// get current logged in user email (username) from jwt
func GetCurrentLoggedUserEmailFromJwt(claims jwt.MapClaims) string {
	return claims["email"].(string)
}

// get current logged in user email (username) from jwt
func GetCurrentLoggedUserUUIDFromJwt(claims jwt.MapClaims) string {
	return claims["sub"].(string)
}
