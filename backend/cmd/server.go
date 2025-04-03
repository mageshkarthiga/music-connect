package main

import (
    "fmt"
    "log"
    "my-go-backend/config"
    "my-go-backend/routes"
    "github.com/labstack/echo/v4"
)

func main() {
    config.Connect()
    defer config.Close()

    e := echo.New()

    routes.SetupRoutes(e)

    fmt.Println("Server is running on port 8080...")
    log.Fatal(e.Start(":8080"))
}
