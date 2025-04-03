package tests

import (
    "testing"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
    "my-go-backend/controllers"
)

func TestCreateUser(t *testing.T) {
    e := echo.New()
    req := httptest.NewRequest(http.MethodPost, "/api/users", nil)
    rec := httptest.NewRecorder()

    c := e.NewContext(req, rec)

    if assert.NoError(t, controllers.CreateUser(c)) {
        assert.Equal(t, http.StatusCreated, rec.Code)
    }
}
