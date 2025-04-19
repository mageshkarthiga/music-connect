package main

import (
	"log"

	"github.com/joho/godotenv"
)

const SPOTIFY_BASE_URL = "https://api.spotify.com/v1"
const SPOTIFY_AUTH_TOKEN = "https://accounts.spotify.com/api/token"
const SUPABASE_URL = "https://kzxuobrnlppliqiwwgvu.supabase.co"

var (
	TICKETMASTER_CALLBACK_URL string
	PYTHON_SCRAPER_URL string

	PORT string
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}
	TICKETMASTER_CALLBACK_URL = getEnv("TICKETMASTER_CALLBACK_URL", "http://localhost:3002/ticketmaster-scrape-callback")
	PYTHON_SCRAPER_URL = getEnv("PYTHON_SCRAPER_URL", "http://localhost:3001/scrape/ticketmaster")
	PORT = getEnv("PORT", "3002")
}
