package utils

import (
    "github.com/labstack/echo/v4"
    "net/http"
)

func RespondJSON(c echo.Context, code int, message string) error {
    return c.JSON(code, map[string]interface{}{
        "message": message,
    })
}
