package main

import (
	"fmt"

	"gorm.io/driver/postgres" // ORM for Golang
	"gorm.io/gorm"
)

// definition of table structure starts here

type User struct {
	UserId       uint   `gorm:"primaryKey;autoIncrement"`
	PhoneNumber  string `gorm:"size:11;unique"`
	EmailAddress string
	Location     string
	UserName     string
}

// definition of table structure ends here


func main() {
	dsn := "host=localhost user=postgres password= dbname=music_connect port=5432" // please change your username and password accordingly
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // open DB connection to postgres via gorm
	if err!=nil{
		fmt.Println("Error connecting to database ⚠️", err)
		return
	}
	fmt.Println("Successfully connected to the database! ✅",db)
}
