package controllers

import (
	"backend/config"
	"backend/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// SendFriendRequest sends a friend request from one user to another
func SendFriendRequest(c echo.Context) error {
	userID := c.Get("uid").(uint) // Get the current user ID from the context

	// Convert friend_id from string to uint
	friendIDStr := c.Param("friend_id")
	friendID, err := strconv.ParseUint(friendIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid friend ID")
	}

	// Check if the friend ID is the same as the user ID
	if userID == uint(friendID) {
		return c.JSON(http.StatusBadRequest, "You cannot send a friend request to yourself")
	}

	// Check if the friendship already exists
	var existingFriendship models.Friendship
	if err := config.DB.Where("user_id = ? AND friend_id = ?", userID, friendID).First(&existingFriendship).Error; err == nil {
		return c.JSON(http.StatusConflict, "Friendship already exists or request is pending")
	}

	// Create a new friendship with status "pending"
	friendships := models.Friendship{
		UserID:    userID,
		FriendID:  uint(friendID),
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	if err := config.DB.Create(&friendships).Error; err != nil {
		log.Printf("Error sending friend request: %v", err)
		return c.JSON(http.StatusInternalServerError, "Failed to send friend request")
	}

	return c.JSON(http.StatusOK, "Friend request sent")
}

// AcceptFriendRequest accepts a pending friend request
func AcceptFriendRequest(c echo.Context) error {
	userID := c.Get("uid").(uint) // Get the current user ID

	// Convert friend_id from string to uint
	friendIDStr := c.Param("friend_id")
	friendID, err := strconv.ParseUint(friendIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid friend ID")
	}

	// Find the friendships with status "pending"
	var friendships models.Friendship
	if err := config.DB.Where("user_id = ? AND friend_id = ? AND status = ?", friendID, userID, "pending").First(&friendships).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Friend request not found or already accepted")
	}

	// Update the friendships status to "accepted"
	friendships.Status = "accepted"
	if err := config.DB.Save(&friendships).Error; err != nil {
		log.Printf("Error accepting friend request: %v", err)
		return c.JSON(http.StatusInternalServerError, "Failed to accept friend request")
	}

	// Optionally, create a reciprocal friendships if needed
	reciprocalfriendships := models.Friendship{
		UserID:    userID,
		FriendID:  uint(friendID),
		Status:    "accepted",
		CreatedAt: time.Now(),
	}
	if err := config.DB.Create(&reciprocalfriendships).Error; err != nil {
		log.Printf("Error creating reciprocal friendships: %v", err)
		return c.JSON(http.StatusInternalServerError, "Failed to create reciprocal friendships")
	}

	return c.JSON(http.StatusOK, "Friend request accepted")
}

// RejectFriendRequest rejects a pending friend request
func RejectFriendRequest(c echo.Context) error {
	userID := c.Get("uid").(uint) // Get the current user ID

	// Convert friend_id from string to uint
	friendIDStr := c.Param("friend_id")
	friendID, err := strconv.ParseUint(friendIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid friend ID")
	}

	// Find the friendships with status "pending"
	var friendships models.Friendship
	if err := config.DB.Where("user_id = ? AND friend_id = ? AND status = ?", friendID, userID, "pending").First(&friendships).Error; err != nil {
		return c.JSON(http.StatusNotFound, "Friend request not found")
	}

	// Delete the pending friend request
	if err := config.DB.Delete(&friendships).Error; err != nil {
		log.Printf("Error rejecting friend request: %v", err)
		return c.JSON(http.StatusInternalServerError, "Failed to reject friend request")
	}

	return c.JSON(http.StatusOK, "Friend request rejected")
}

// RemoveFriend removes a friend from the user's friends list
func RemoveFriend(c echo.Context) error {
	userID := c.Get("uid").(uint) // Get the current user ID

	// Convert friend_id from string to uint
	friendIDStr := c.Param("friend_id")
	friendID, err := strconv.ParseUint(friendIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid friend ID")
	}

	// Prepare the raw SQL queries to remove friendships in both directions
	// First direction: userID -> friendID
	query1 := "DELETE FROM friendships WHERE user_id = ? AND friend_id = ?"
	if err := config.DB.Exec(query1, userID, friendID).Error; err != nil {
		log.Printf("Error removing friendship (userID -> friendID): %v", err)
		return c.JSON(http.StatusInternalServerError, "Failed to remove friendship (userID -> friendID)")
	}

	// Second direction: friendID -> userID
	query2 := "DELETE FROM friendships WHERE user_id = ? AND friend_id = ?"
	if err := config.DB.Exec(query2, friendID, userID).Error; err != nil {
		log.Printf("Error removing reciprocal friendship (friendID -> userID): %v", err)
		return c.JSON(http.StatusInternalServerError, "Failed to remove reciprocal friendship (friendID -> userID)")
	}

	return c.JSON(http.StatusOK, "Friend removed successfully")
}

// GetFriends returns a list of friends for the current user
func GetFriends(c echo.Context) error {
	userID := c.Get("uid").(uint) // Get the current user ID

	var friends []models.User
	// Fetch friends of the current user
	if err := config.DB.Table("users").Joins("join friendships on friendships.friend_id = users.user_id").
		Where("friendships.user_id = ? AND friendships.status = ?", userID, "accepted").
		Find(&friends).Error; err != nil {
		log.Printf("Error fetching friends: %v", err)
		return c.JSON(http.StatusInternalServerError, "Failed to fetch friends")
	}

	return c.JSON(http.StatusOK, friends)
}

// GetPendingFriendRequests retrieves a list of pending friend requests for the current user
func GetPendingFriendRequests(c echo.Context) error {
	userID := c.Get("uid").(uint) // Get the current user ID

	var pendingRequests []models.User
	// Fetch pending friend requests where the current user is the friend_id
	if err := config.DB.Table("users").Joins("join friendships on friendships.user_id = users.user_id").
		Where("friendships.friend_id = ? AND friendships.status = ?", userID, "pending").
		Find(&pendingRequests).Error; err != nil {
		log.Printf("Error fetching pending friend requests: %v", err)
		return c.JSON(http.StatusInternalServerError, "Failed to fetch pending friend requests")
	}

	// Calculate the count of pending requests
	count := len(pendingRequests)

	// Return the pending requests along with the count
	return c.JSON(http.StatusOK, map[string]interface{}{
		"count": count,
		"users": pendingRequests,
	})
}

// Get status of friendship between two users
func GetFriendshipStatus(c echo.Context) error {
	userID := c.Get("uid").(uint) // Get the current user ID

	// Convert friend_id from string to uint64
	friendIDStr := c.Param("friend_id")
	friendID, err := strconv.ParseUint(friendIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid friend ID")
	}

	var friendship models.Friendship

	// 1. Check if current user sent the friend request
	if err := config.DB.Where("user_id = ? AND friend_id = ?", userID, friendID).First(&friendship).Error; err == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":     friendship.Status,
			"created_at": friendship.CreatedAt,
			"user_id":    friendship.UserID,
			"friend_id":  friendship.FriendID,
		})
	}

	// 2. Check if friend sent the request to the current user
	if err := config.DB.Where("user_id = ? AND friend_id = ?", friendID, userID).First(&friendship).Error; err == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":     friendship.Status,
			"created_at": friendship.CreatedAt,
			"user_id":    friendship.UserID,
			"friend_id":  friendship.FriendID,
		})
	}

	// 3. Not found in either direction
	return c.JSON(http.StatusNotFound, "Friendship not found")
}
