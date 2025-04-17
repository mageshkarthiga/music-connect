package main

import (
	"log"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	if err := godotenv.Load(".env"); 
	err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	c := cron.New()
	// c.AddFunc("0 0 * * 0", func() {
		c.AddFunc("54 13 * * 4", func() {
		log.Println("Running weekly job...")
		go callSpotifyAPI()

		triggerScraper()
		go CallbackServer()
	})
	c.Start()
	select {}
}
