package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

// Global DB variable to hold the connection
var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() error {
	if os.Getenv("ENV") != "prod" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file:", err)
		}
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return fmt.Errorf("DB_PASSWORD environment variable is not set")
	}

	dsn := fmt.Sprintf("postgresql://postgres.kzxuobrnlppliqiwwgvu:%s@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres?statement_cache_mode=off", dbPassword)

	DB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // also disables prep statements at protocol level
	}), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: false,
	})
	if err != nil {
		log.Fatal("Error connecting to database:", err, DB)
	}

	// Check if the 'users' table exists before running migrations
	// if !tableExists(DB, "users")

	// 	err = DB.AutoMigrate(&models.User{}, &models.Event{}, &models.Attraction{}, &models.Venue{}, &models.Artist{}, &models.Track{}, &models.Playlist{}, &models.Device{})
	// 	if err != nil {
	// 		return fmt.Errorf("Error during migration: %v", err)
	// 	}
	// 	fmt.Println("Migration completed successfully! ✅")
	// } else {
	// 	fmt.Println("Tables already exist, skipping migration.")
	//

	return nil
}

// tableExists checks if a table exists in the database
func tableExists(db *gorm.DB, tableName string) bool {
	var count int64
	err := db.Raw("SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = ?", tableName).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}
