package auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

// SetAuthCookie sets the auth_token cookie in the response
func SetAuthCookie(c echo.Context, tokenStr string) error {
	appEnv := os.Getenv("APP_ENV")

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    tokenStr,
		HttpOnly: true, 
		Path:     "/",
		MaxAge:   36000000,
	}

	if appEnv == "production" {
		// For deployed environments
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode // Needed for cross-site cookies with HTTPS
	} else {
		// For local development
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode // Lax allows POSTs to work on same-site dev
	}

	c.SetCookie(cookie)
	return nil
}
