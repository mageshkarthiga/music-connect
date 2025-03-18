package main

import (
	"fmt"
	"gorm.io/driver/postgres" // ORM for Golang
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"github.com/joho/godotenv"
)

// definition of table structure starts here

type User struct {
	UserID       uint   `gorm:"primaryKey;autoIncrement"`
	PhoneNumber  string `gorm:"size:11;unique"`
	EmailAddress string
	Location     string
	UserName     string
	Events       []*Event `gorm:"many2many:user_events"`
	Tracks       []*Track `gorm:"many2many:music_preferences"`
}

type Event struct {
	EventID          uint `gorm:"primaryKey;autoIncrement"`
	EventName        string
	EventDescription string
	EventURL         string
	EventImageURL    string
	Users            []*User       `gorm:"many2many:user_events"`
	Attractions      []*Attraction `gorm:"many2many:event_attractions"`
	Venue            []*Venue      `gorm:"many2many:user_venues"`
}

type Attraction struct {
	AttractionID       uint `gorm:"primaryKey;autoIncrement"`
	AttractionName     string
	AttractionType     string
	AttractionURL      string
	AttractionImageURL string
	Events             []*Event `gorm:"many2many:event_attractions"`
}

type Venue struct {
	VenueID       uint `gorm:"primaryKey;autoIncrement"`
	VenueName     string
	Location      string
	SeatMap       string
	Accessibility string
	Events        []*Event `gorm:"many2many:event_venues"`
}

type Artist struct {
	ArtistID       uint `gorm:"primaryKey;autoIncrement"`
	ArtistName     string
	ArtistImageURL string
	Tracks         []*Track `gorm:"many2many:track_artists"`
}

type Track struct {
	TrackID    uint `gorm:"primaryKey;autoIncrement"`
	TrackTitle string
	ArtistID   uint `gorm:"foreignKey:artist_id"`
	Genre      string
	TrackURL   string
	Artists    []*Artist   `gorm:"many2many:track_artists"`
	Playlists  []*Playlist `gorm:"many2many:playlist_tracks"`
}

type Playlist struct {
	PlaylistID   uint `gorm:"primaryKey;autoIncrement"`
	PlaylistName string
	Tracks       []*Track `gorm:"many2many:playlist_tracks"`
}

// definition of table structure ends here

func main() {
	err := godotenv.Load()
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
	tables := []interface{}{&User{}, &Event{}, &Attraction{}, &Venue{}, &Artist{}, &Track{}, &Playlist{}}
	for _, table := range tables {
		if db.Migrator().HasTable(table) {
			db.Migrator().DropTable(table)
		}
	}

	// Commit tables to DB
	err = db.AutoMigrate(&User{}, &Event{}, &Attraction{}, &Venue{}, &Artist{}, &Track{}, &Playlist{})
	if err != nil {
		fmt.Println("Error during migration ⚠️", err)
		return
	}

}
