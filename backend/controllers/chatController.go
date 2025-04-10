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