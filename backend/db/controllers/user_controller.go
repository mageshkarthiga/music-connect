
package controllers

import (
	"net/http"
"fmt"
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
	id := c.Param("UserID")
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
	id := c.Param("UserID")

	var user models.User
	if err := uc.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": fmt.Sprintf("User with ID %s not found", id),
		})
	}

	var updates models.User
	if err := c.Bind(&updates); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error":      "Failed to parse JSON input",
			"parseError": err.Error(),
		})
	}

	fmt.Printf("Received update body: %+v\n", updates)

	// Only update non-zero/non-empty fields
	if updates.UserName != "" {
		user.UserName = updates.UserName
	}
	if updates.EmailAddress != "" {
		user.EmailAddress = updates.EmailAddress
	}
	if updates.PhoneNumber != "" {
		user.PhoneNumber = updates.PhoneNumber
	}
	if updates.Location != "" {
		user.Location = updates.Location
	}
	if updates.ProfilePhotoUrl != "" {
		user.ProfilePhotoUrl = updates.ProfilePhotoUrl
	}

	if err := uc.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error":         "Database save failed",
			"databaseError": err.Error(),
		})
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

func (uc *UserController) GetUserByFirebaseUID(c echo.Context) error {
	uid := c.Param("uid")
	var user models.User
	if err := uc.DB.Where("firebase_uid = ?", uid).First(&user).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, user)
}
