package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"backend/config"
	"backend/models"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

var jwks *keyfunc.JWKS

// InitJWKS initializes the JWKS (JSON Web Key Set) from Firebase
func InitJWKS() error {
	// JWKS endpoint for Firebase
	jwksURL := "https://www.googleapis.com/service_accounts/v1/jwk/securetoken@system.gserviceaccount.com"

	options := keyfunc.Options{
		RefreshInterval:  time.Hour,
		RefreshRateLimit: time.Minute * 5,
		RefreshErrorHandler: func(err error) {
			fmt.Printf("There was an error with the jwt.KeyFunc\nError: %s\n", err.Error())
		},
		RefreshUnknownKID: true,
	}

	var err error
	jwks, err = keyfunc.Get(jwksURL, options)
	return err
}

// AuthMiddleware authenticates the JWT token in the Authorization header or cookie
func AuthMiddleware(projectID string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var tokenStr string

			// 1. Try Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				// 2. Try cookie as fallback
				cookie, err := c.Cookie("auth_token")
				if err != nil {
					return c.JSON(http.StatusUnauthorized, "Missing Authorization header or auth_token cookie")
				}
				tokenStr = cookie.Value
			}

			// Parse the JWT token
			token, err := jwt.Parse(tokenStr, GetJWKS().Keyfunc)
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
			}

			// Optional: Extract and log 'kid'
			kid, _ := token.Header["kid"].(string)
			// fmt.Printf("JWT Header kid: %s\n", kid)

			_, err = GetJWKS().Keyfunc(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, fmt.Sprintf("Key ID '%s' not found in JWKS", kid))
			}

			// Extract claims
			claims, ok := token.Claims.(jwt.MapClaims)
			log.Printf("Claims: %+v\n", claims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, "Invalid token claims")
			}

			// // Validate issuer
			// expectedIss := fmt.Sprintf("https://securetoken.google.com/%s", projectID)
			// if iss, ok := claims["iss"].(string); !ok || iss != expectedIss {
			// 	return c.JSON(http.StatusUnauthorized, "Invalid issuer")
			// }

			// Optional: Validate audience
			// if aud, ok := claims["aud"].(string); !ok || aud != projectID {
			// 	return c.JSON(http.StatusUnauthorized, "Invalid audience")
			// }

			// Extract user_id (uid)
			uid, ok := claims["user_id"].(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, "User ID missing in token")
			}

			// Query user from the database directly
			var user models.User
			if err := config.DB.Where("firebase_uid = ?", uid).First(&user).Error; err != nil {
				return c.JSON(http.StatusUnauthorized, "User not found")
			}
			log.Printf("User found: %+v\n", user)

			// Set user ID in context for downstream handlers
			c.Set("uid", user.UserID) // type is uint NOT string
			return next(c)
		}
	}
}

// GetJWKS returns the initialized JWKS
func GetJWKS() *keyfunc.JWKS {
	return jwks
}
