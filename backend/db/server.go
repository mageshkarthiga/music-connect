// backend/db/server.go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"music-connect/db/controllers"
	"github.com/labstack/echo/v4/middleware"

)




func main() {
	_ = godotenv.Load("../../.env")

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("DB_PASSWORD environment variable is not set")
	}

	dsn := fmt.Sprintf("postgresql://postgres.kzxuobrnlppliqiwwgvu:%s@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres?statement_cache_mode=describe", dbPassword)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),		
	})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
	}))
	//TODO: Change this before deployment

	userController := controllers.NewUserController(db)
	e.POST("/users", userController.CreateUser)
	e.GET("/users", userController.GetAllUsers)
	e.GET("/users/:id", userController.GetUser)
	e.PUT("/users/:id", userController.UpdateUser)
	e.DELETE("/users/:id", userController.DeleteUser)
	e.GET("/users/firebase/:uid", userController.GetUserByFirebaseUID)

	e.Logger.Fatal(e.Start(":1323"))
}
