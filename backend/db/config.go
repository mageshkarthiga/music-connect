package main

import (
	"fmt"
	"gorm.io/driver/postgres" // ORM for Golang
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"github.com/joho/godotenv"
	"music-connect/db/models"
)

// definition of table structure starts here
//moved to /models/


// definition of table structure ends here

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		fmt.Println("DB_PASSWORD environment variable is not set")
		return
	}

	dsn := fmt.Sprintf("postgresql://postgres:%s@db.kzxuobrnlppliqiwwgvu.supabase.co:5432/postgres",dbPassword) 
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}) // open DB connection to supabase postgres via gorm
	if err != nil {
		fmt.Println("Error connecting to database ⚠️", err)
		return
	}
	fmt.Println("Successfully connected to the database! ✅", db)

	// Check if tables exist and drop if they do
	manyToManyTables := []string{
		"user_events", "music_preferences", "event_attractions",
		"event_venues", "track_artists", "playlist_tracks",
	}
	for _, tableName := range manyToManyTables {
		if db.Migrator().HasTable(tableName) {
			db.Migrator().DropTable(tableName)
		}
	}

	tables := []interface{}{
		&models.User{}, &models.Event{}, &models.Attraction{},
		&models.Venue{}, &models.Artist{}, &models.Track{}, &models.Playlist{},
	}
	for _, table := range tables {
		if db.Migrator().HasTable(table) {
			db.Migrator().DropTable(table)
		}
	}

	// Commit tables to DB
	err = db.AutoMigrate(&models.User{}, &models.Event{}, &models.Attraction{}, &models.Venue{}, &models.Artist{}, &models.Track{}, &models.Playlist{})
	if err != nil {
		fmt.Println("Error during migration ⚠️", err)
		return
	}

}
