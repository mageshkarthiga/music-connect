package controllers

import (
    "backend/services/chat"
    "net/http"
    "github.com/labstack/echo/v4"
)

func GetMessagesForRoom(c echo.Context) error {
	roomName := c.Param("roomName")
	messages, err := chat.GetMessagesForRoom(roomName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, messages)
}

func GetUsersWithChatHistory(c echo.Context) error {
	userID := c.Param("userID")
	users, err := chat.GetUsersWithChatHistory(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}