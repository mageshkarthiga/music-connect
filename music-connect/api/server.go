package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New() // new instance of Echo server
	e.GET("/", func (c echo.Context) error { // context of API endpoint
		return c.String(http.StatusOK, "Hello, World")
	})
	e.Logger.Fatal(e.Start(":1323")) // starts the server on the specified port number
} 