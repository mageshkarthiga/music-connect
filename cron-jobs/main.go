package main

import (
	"log"
	"github.com/joho/godotenv"
	// "github.com/robfig/cron/v3"
)

const SUPABASE_URL = "https://kzxuobrnlppliqiwwgvu.supabase.co";

func main() {
	if err := godotenv.Load("../.env"); 
	err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	callSpotifyAPI()
	// go CallbackServer()

	// c := cron.New()
	// c.AddFunc("0 0 * * 0", func() {
	// 	callSpotifyAPI()
	// 	// callTicketMasterAPI("https://app.ticketmaster.com/discovery/v2/events.json")
	// })
	// c.Start()
	// select {}
}
