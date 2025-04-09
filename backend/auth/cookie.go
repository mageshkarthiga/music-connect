package auth

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// SetAuthCookie sets the auth_token cookie in the response
func SetAuthCookie(c echo.Context, tokenStr string) error {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    tokenStr,
		HttpOnly: false,
		Secure:   false, 
		Path:     "/",
		SameSite: http.SameSiteLaxMode, 
		MaxAge:   360000000,
	}

	// Set the cookie in the response
	c.SetCookie(cookie)

	return nil
}
