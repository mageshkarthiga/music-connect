package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/MicahParks/keyfunc"
)

var jwks *keyfunc.JWKS

// InitJWKS initializes the JWKS (JSON Web Key Set) from Firebase
func InitJWKS() error {
	// JWKS endpoint for Firebase
	jwksURL := "https://www.googleapis.com/service_accounts/v1/jwk/securetoken@system.gserviceaccount.com"

	options := keyfunc.Options{
		RefreshInterval:   time.Hour,
		RefreshRateLimit:  time.Minute * 5,
		RefreshErrorHandler: func(err error) {
			fmt.Printf("There was an error with the jwt.KeyFunc\nError: %s\n", err.Error())
		},
		RefreshUnknownKID: true,
	}

	var err error
	jwks, err = keyfunc.Get(jwksURL, options)
	return err
}

// AuthMiddleware authenticates the JWT token in the Authorization header
func AuthMiddleware(projectID string) echo.MiddlewareFunc {


	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, "Missing Authorization header")
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, "Invalid Authorization format")
			}

			tokenStr := parts[1]

			// Parse the JWT token
			token, err := jwt.Parse(tokenStr, GetJWKS().Keyfunc)
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
			}

			// Extract 'kid' from the token header
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, "Missing or invalid 'kid' in token header")
			}

			// Log the 'kid' from the token header for debugging
			fmt.Printf("JWT Header kid: %s\n", kid)

			// Attempt to get the key by kid using jwks.Keyfunc
			key, err := GetJWKS().Keyfunc(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, fmt.Sprintf("Key ID '%s' not found in JWKS", kid))
			}

			// Log the matched key for debugging
			fmt.Printf("Found matching key: %+v\n", key)

			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, "Invalid token claims")
			}

            projectID = "music-connect-608f6" // Set the project ID here

			// // Ensure the audience is correct
			// aud, ok := claims["aud"].(string)
			// if !ok || aud != projectID {
			// 	return c.JSON(http.StatusUnauthorized, "Invalid audience")
			// }

			// Ensure the issuer is correct
			iss, ok := claims["iss"].(string)
			if !ok || iss != fmt.Sprintf("https://securetoken.google.com/%s", projectID) {
				return c.JSON(http.StatusUnauthorized, "Invalid issuer")
			}

			// Extract user_id (uid) from claims
			uid, ok := claims["user_id"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, "User ID missing in token")
			}

			// Set the user_id in the context
			c.Set("uid", uid)

			return next(c)
		}
	}
}

// GetJWKS returns the initialized JWKS
func GetJWKS() *keyfunc.JWKS {
	return jwks
}
