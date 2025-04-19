// controllers/userController.go

package controllers

import (
	"backend/config"
	"backend/models"
	// "fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

// GetUsers fetches all users
func GetUsers(c echo.Context) error {
	var users []models.User
	// Fetch all users from the database
	if err := config.DB.Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to fetch users")
	}
	return c.JSON(http.StatusOK, users)
}

// GetUser fetches current user 
func GetMe(c echo.Context) error {
	userID := c.Get("uid")
	log.Printf("User ID from context: %v\n", userID)

    var user models.User
    if err := config.DB.First(&user, userID).Error; err != nil {
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch user"})
    }

	return c.JSON(http.StatusOK, user)
}



// CreateUser creates a new user
func CreateUser(c echo.Context) error {
	var user models.User

	// Bind the request body to the user struct
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	// Log the user object to check if it's being populated correctly
	log.Printf("User data received: %+v", user)

	// Validate the fields
	if user.PhoneNumber == "" || user.EmailAddress == "" || user.UserName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Missing required fields"})
	}

	// Insert the user into the database
	if err := config.DB.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create user"})
	}

	// Successfully created user
	return c.JSON(http.StatusOK, user)
}


func UpdateUser(c echo.Context) error {
	id := c.Param("UserID")
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to update user")
	}
	return c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a user by ID
func DeleteUser(c echo.Context) error {
	id := c.Param("id")
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to delete user")
	}
	return c.JSON(http.StatusOK, "User deleted successfully")
}

// GetUserByFirebaseUID fetches user by Firebase UID
func GetUserByFirebaseUID(c echo.Context) error {
	uid := c.Param("uid")
	var user models.User
	if err := config.DB.Where("firebase_uid = ?", uid).First(&user).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, user)
}

func GetUserByUserID(c echo.Context) error {
	id := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, "user_id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, user)
}

// AddFriend creates a friendship between two users
func AddFriend(c echo.Context) error {
	var payload struct {
		UserID   uint `json:"user_id"`
		FriendID uint `json:"friend_id"`
	}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	var user models.User
	if err := config.DB.Preload("Friends").First(&user, payload.UserID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	var friend models.User
	if err := config.DB.First(&friend, payload.FriendID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Friend not found")
	}

	// Prevent duplicate friendships
	for _, existing := range user.Friends {
		if existing.UserID == friend.UserID {
			return c.JSON(http.StatusBadRequest, "Already friends")
		}
	}

	// Append friend
	if err := config.DB.Model(&user).Association("Friends").Append(&friend); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to add friend")
	}

	return c.JSON(http.StatusOK, "Friend added successfully")
}


// AddFriend creates a friendship between two users
func AddFriend(c echo.Context) error {
	var payload struct {
		UserID   uint `json:"user_id"`
		FriendID uint `json:"friend_id"`
	}

	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request")
	}

	var user models.User
	if err := config.DB.Preload("Friends").First(&user, payload.UserID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	var friend models.User
	if err := config.DB.First(&friend, payload.FriendID).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Friend not found")
	}

	// Prevent duplicate friendships
	for _, existing := range user.Friends {
		if existing.UserID == friend.UserID {
			return c.JSON(http.StatusBadRequest, "Already friends")
		}
	}

	// Append friend
	if err := config.DB.Model(&user).Association("Friends").Append(&friend); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to add friend")
	}

	return c.JSON(http.StatusOK, "Friend added successfully")
}
