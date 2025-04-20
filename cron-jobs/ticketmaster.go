package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type RequestBody struct {
	CallbackUrl string `json:"callbackUrl"`
}

type ScrapedData struct {
	EventName        string `json:"eventName"`
	EventDescription string `json:"eventDescription"`
	EventUrl         string `json:"eventUrl"`
	EventImageUrl    string `json:"eventImageUrl"`
	Location         string `json:"location"`
	VenueName        string `json:"venueName"`
}

type Event struct {
	EventID          int    `json:"event_id"`
	EventName        string `json:"event_name"`
	EventDescription string `json:"event_description"`
	EventImageUrl    string `json:"event_image_url"`
	EventUrl         string `json:"event_url"`
}

type Venue struct {
	VenueID   int    `json:"venue_id"`
	VenueName string `json:"venue_name"`
	Location  string `json:"location"`
}

func triggerScraper() {
	requestBody := RequestBody{CallbackUrl: TICKETMASTER_CALLBACK_URL}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("error triggering scraper: %v", err)
	}

	resp, err := http.Post(PYTHON_SCRAPER_URL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("failed to trigger scraper: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Triggered scraper, status code: %d\n", resp.StatusCode)
}
