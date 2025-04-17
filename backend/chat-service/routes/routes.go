package routes

import (
    "chat-service/controllers"
    "chat-service/server"
    "github.com/labstack/echo/v4"
    "fmt"
)

func RegisterRoutes(e *echo.Echo, wsServer *server.WsServer) {
    fmt.Println("Registering routes...")
    // WebSocket endpoint
    e.GET("/ws", func(c echo.Context) error {
        server.ServeWs(wsServer, c.Response().Writer, c.Request()) // Use the WebSocket handler
        return nil
    })

    // Message Retrieval Route
    e.GET("/rooms/:roomName/messages", controllers.GetMessagesForRoom)

    // User Chat History Route
    e.GET("/users/:userID/chat-history", controllers.GetUsersWithChatHistory)
}