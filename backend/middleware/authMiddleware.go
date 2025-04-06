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

func InitJWKS() error {
    // JWKS endpoint for Firebase
    jwksURL := "https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com"

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

            token, err := jwt.Parse(tokenStr, jwks.Keyfunc)
            if err != nil {
                return c.JSON(http.StatusUnauthorized, fmt.Sprintf("Invalid token: %v", err))
            }

            if !token.Valid {
                return c.JSON(http.StatusUnauthorized, "Invalid token")
            }

            claims := token.Claims.(jwt.MapClaims)

            // Verify important claims
            if claims["aud"] != projectID {
                return c.JSON(http.StatusUnauthorized, "Invalid audience")
            }
            if claims["iss"] != fmt.Sprintf("https://securetoken.google.com/%s", projectID) {
                return c.JSON(http.StatusUnauthorized, "Invalid issuer")
            }

            uid := claims["user_id"].(string)
            c.Set("uid", uid)

            return next(c)
        }
    }
}
