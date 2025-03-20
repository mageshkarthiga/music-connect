// backend/db/controllers/user_controller.go
package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"music-connect/db/models" 
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) CreateUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}
	if err := uc.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	return c.JSON(http.StatusCreated, user)
}

func (uc *UserController) GetUser(c echo.Context) error {
	id := c.Param("id")
	var user models.User
	if err := uc.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Not found"})
	}
	return c.JSON(http.StatusOK, user)
}

func (uc *UserController) GetAllUsers(c echo.Context) error {
	var users []models.User
	if err := uc.DB.Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	return c.JSON(http.StatusOK, users)
}

func (uc *UserController) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	var user models.User
	if err := uc.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Not found"})
	}

	updates := new(models.User)
	if err := c.Bind(updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	user.PhoneNumber = updates.PhoneNumber
	user.EmailAddress = updates.EmailAddress
	user.Location = updates.Location
	user.UserName = updates.UserName

	if err := uc.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	return c.JSON(http.StatusOK, user)
}

func (uc *UserController) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	if err := uc.DB.Delete(&models.User{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	return c.NoContent(http.StatusNoContent)
}
