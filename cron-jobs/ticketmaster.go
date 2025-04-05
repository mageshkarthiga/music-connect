package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
	requestBody := RequestBody{CallbackUrl: "http://localhost:3002/ticketmaster-scrape-callback"}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("error triggering scraper: %v", err)
	}

	resp, err := http.Post("http://localhost:3001/scrape/ticketmaster", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("failed to trigger scraper: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("Triggered scraper, status code: %d\n", resp.StatusCode)
}

func getVenuesInDb() ([]Venue, error) {
	apiKey := os.Getenv("publicApiKey")
	url := "https://kzxuobrnlppliqiwwgvu.supabase.co/rest/v1/venues"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var venues []Venue
	if err := json.NewDecoder(resp.Body).Decode(&venues); err != nil {
		return nil, err
	}

	return venues, nil
}

func findVenueIDByName(venues []Venue, scraped ScrapedData) (int, bool) {
	for _, venue := range venues {
		if venue.VenueName == scraped.VenueName {
			return venue.VenueID, true
		}
	}
	return 0, false
}

func addEventVenueToDB(eventId int, event ScrapedData) error {
	apiKey := os.Getenv("publicApiKey")
	venues, _ := getVenuesInDb()

	venueID, found := findVenueIDByName(venues, event)
	if !found {
		return fmt.Errorf("venue not found")
	}

	url := "https://kzxuobrnlppliqiwwgvu.supabase.co/rest/v1/event_venues"
	data := map[string]interface{}{
		"venue_id": venueID,
		"event_id": eventId,
	}

	payload, _ := json.Marshal(data)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to add connection between event and venue: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

func handleScrapeCallback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	var data []ScrapedData
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// add event to database
	apiKey := os.Getenv("publicApiKey")
	url := "https://kzxuobrnlppliqiwwgvu.supabase.co/rest/v1/events"
	for _, event := range data {
		eventData := map[string]interface{}{
			"event_name":        event.EventName,
			"event_description": event.EventDescription,
			"event_image_url":   event.EventImageUrl,
			"event_url":         event.EventUrl,
		}
		payload, _ := json.Marshal(eventData)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("apikey", apiKey)
		req.Header.Set("Authorization", "Bearer " + apiKey)
		req.Header.Set("Prefer", "return=representation")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, "Failed to insert event into database", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		
		if resp.StatusCode == http.StatusCreated {
			log.Printf("Event %s added to database successfully\n", event.EventName)

			body, _ := io.ReadAll(resp.Body)
			var insertedEvent []Event
			if err := json.Unmarshal(body, &insertedEvent); err != nil {
				log.Printf("unmarshal error: %v", err)
			}
			addEventVenueToDB(insertedEvent[0].EventID, event)
		} else {
			http.Error(w, "Failed to insert event into database", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Scraped data received and added into database successfully")
}

func CallbackServer() {
	http.HandleFunc("/ticketmaster-scrape-callback", handleScrapeCallback)
	log.Println("Server listening on port 3002")
	log.Fatal(http.ListenAndServe(":3002", nil))
}