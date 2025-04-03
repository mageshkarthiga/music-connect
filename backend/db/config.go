package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres" // ORM for Golang
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

// definition of table structure starts here

type User struct {
	UserID       uint   `gorm:"primaryKey;autoIncrement"`
	PhoneNumber  string `gorm:"size:11;unique"`
	EmailAddress string
	Location     string
	UserName     string
	Events       []*Event `gorm:"many2many:user_events;joinForeignKey:UserID;joinReferences:EventID"`
	Tracks       []*Track `gorm:"many2many:music_preferences;joinForeignKey:UserID;joinReferences:TrackID"`
	Devices      []Device
}

type Event struct {
	EventID          uint `gorm:"primaryKey;autoIncrement"`
	EventName        string
	EventDescription string
	EventURL         string
	EventImageURL    string
	Users            []*User       `gorm:"many2many:user_events;joinForeignKey:EventID;joinReferences:UserID"`
	Attractions      []*Attraction `gorm:"many2many:event_attractions;joinForeignKey:EventID;joinReferences:AttractionID"`
	Venues           []*Venue      `gorm:"many2many:event_venues;joinForeignKey:EventID;joinReferences:VenueID"`
}

type Attraction struct {
	AttractionID       uint `gorm:"primaryKey;autoIncrement"`
	AttractionName     string
	AttractionType     string
	AttractionURL      string
	AttractionImageURL string
	Events             []*Event `gorm:"many2many:event_attractions;joinForeignKey:AttractionID;joinReferences:EventID"`
}

type Venue struct {
	VenueID       uint `gorm:"primaryKey;autoIncrement"`
	VenueName     string
	Location      string
	SeatMap       string
	Accessibility string
	Events        []*Event `gorm:"many2many:event_venues;joinForeignKey:VenueID;joinReferences:EventID"`
}

type Artist struct {
	ArtistID        uint `gorm:"primaryKey;autoIncrement"`
	ArtistName      string
	ArtistSpotifyID string
	Tracks          []*Track `gorm:"many2many:track_artists;joinForeignKey:ArtistID;joinReferences:TrackID"`
}

type Track struct {
	TrackID    uint `gorm:"primaryKey;autoIncrement"`
	TrackTitle string
	ArtistID   uint `gorm:"foreignKey:artist_id"`
	Genre      string
	TrackURI   string
	Artists    []*Artist   `gorm:"many2many:track_artists;joinForeignKey:TrackID;joinReferences:ArtistID"`
	Playlists  []*Playlist `gorm:"many2many:playlist_tracks;joinForeignKey:TrackID;joinReferences:PlaylistID"`
}

type Playlist struct {
	PlaylistID   uint `gorm:"primaryKey;autoIncrement"`
	PlaylistName string
	Tracks       []*Track `gorm:"many2many:playlist_tracks;joinForeignKey:PlaylistID;joinReferences:TrackID"`
}

type Device struct {
	DeviceID   uint `gorm:"primaryKey;autoIncrement"`
	DeviceName string
	DeviceUUID string `gorm:"uniqueIndex"`
	UserID     uint
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

	dsn := fmt.Sprintf("postgresql://postgres.kzxuobrnlppliqiwwgvu:%s@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres", dbPassword)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logger.Info),
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
		&Device{}, &Playlist{}, &Track{}, &Artist{},
		&Venue{}, &Attraction{}, &Event{}, &User{},
	}
	for _, table := range tables {
		if db.Migrator().HasTable(table) {
			err := db.Migrator().DropTable(table)
			if err != nil {
				fmt.Printf("Error dropping table %T: %v\n", table, err)
				return
			}
		}
	}
	// Commit tables to DB
	err = db.AutoMigrate(&User{}, &Event{}, &Attraction{}, &Venue{}, &Artist{}, &Track{}, &Playlist{}, &Device{})
	session := db.Session(&gorm.Session{PrepareStmt: true})
	if session != nil {
		fmt.Println("Migration successful")
	}

	fmt.Println("Migration completed successfully!")

	// Add dummy data
	artist1 := Artist{
		ArtistName:      "ROSE",
		ArtistSpotifyID: "3eVa5w3URK5duf6eyVDbu9",
	}

	artist2 := Artist{
		ArtistName:      "Kendrick Lamar",
		ArtistSpotifyID: "2YZyLoL8N0Wb9xBt1NhZWg",
	}

	artist3 := Artist{
		ArtistName:      "LE SSERAFIM",
		ArtistSpotifyID: "4SpbR6yFEvexJuaBpgAU5p",
	}

	artist4 := Artist{
		ArtistName:      "Selena Gomez",
		ArtistSpotifyID: "0C8ZW7ezQVs4URX5aX7Kqx",
	}

	artist5 := Artist{
		ArtistName:      "Bruno Mars",
		ArtistSpotifyID: "0du5cEVh5yTK9QJze8zA0C",
	}

	artist6 := Artist{
		ArtistName:      "Sabrina Carpenter",
		ArtistSpotifyID: "74KM79TiuVKeVCqs8QtB0B",
	}

	artist7 := Artist{
		ArtistName:      "Doechii",
		ArtistSpotifyID: "4E2rKHVDssGJm2SCDOMMJB",
	}

	artist8 := Artist{
		ArtistName:      "Benson Boone",
		ArtistSpotifyID: "22wbnEMDvgVIAGdFeek6ET",
	}
	artist9 := Artist{
		ArtistName:      "Adrianna Cinta",
		ArtistSpotifyID: "7aGKWIJ44Gs7eQ7cCKVskG",
	}
	artist10 := Artist{
		ArtistName:      "James Arthur",
		ArtistSpotifyID: "4IWBUUAFIplrNtaOHcJPRM",
	}

	//Insert artists into the database
	err = db.Create(&artist1).Error
	if err != nil {
		fmt.Println("Error inserting artist1:", err)
		return
	}
	err = db.Create(&artist2).Error
	if err != nil {
		fmt.Println("Error inserting artist2:", err)
		return
	}
	err = db.Create(&artist3).Error
	if err != nil {
		fmt.Println("Error inserting artist3:", err)
		return
	}
	err = db.Create(&artist4).Error
	if err != nil {
		fmt.Println("Error inserting artist4:", err)
		return
	}
	err = db.Create(&artist5).Error
	if err != nil {
		fmt.Println("Error inserting artist5:", err)
		return
	}
	err = db.Create(&artist6).Error
	if err != nil {
		fmt.Println("Error inserting artist6:", err)
		return
	}
	err = db.Create(&artist7).Error
	if err != nil {
		fmt.Println("Error inserting artist7:", err)
		return
	}
	err = db.Create(&artist8).Error
	if err != nil {
		fmt.Println("Error inserting artist8:", err)
		return
	}
	err = db.Create(&artist9).Error
	if err != nil {
		fmt.Println("Error inserting artist9:", err)
		return
	}
	err = db.Create(&artist10).Error
	if err != nil {
		fmt.Println("Error inserting artist10:", err)
		return
	}
}
