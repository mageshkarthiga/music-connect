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

	// Use '=' to assign to the global DB variable
	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // also disables prep statements at protocol level
	}), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: false,
	})
	if err != nil {
		log.Fatal("Error connecting to database:", err)
		return err
	}

	log.Println("✅ Database connection initialized!")

	// Optional: Perform migration if needed
	// err = DB.AutoMigrate(&models.User{}, &models.Event{})
	// if err != nil {
	// 	log.Fatalf("❌ Migration failed: %v", err)
	// }

	return nil
}
