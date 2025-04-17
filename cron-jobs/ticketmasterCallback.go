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

func getVenuesInDb() ([]Venue, error) {
	apiKey := os.Getenv("publicApiKey")
	url := SUPABASE_URL + "/rest/v1/venues"

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

	url := SUPABASE_URL + "/rest/v1/event_venues"
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

func getAllEventUrls(apiKey string) ([]string, error) {
	url := SUPABASE_URL + "/rest/v1/events?select=event_url"

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

	var eventUrls []struct {
		EventUrl string `json:"event_url"`
	}
	var urls []string
	if err := json.NewDecoder(resp.Body).Decode(&eventUrls); err != nil {
		return nil, err
	}
	for _, row := range eventUrls {
		urls = append(urls, row.EventUrl)
	}

	return urls, nil
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

	supabaseApiKey := os.Getenv("publicApiKey")
	url := SUPABASE_URL + "/rest/v1/events"

	// get all event_urls from database
	eventUrls, err := getAllEventUrls(supabaseApiKey)
	if err != nil {
		log.Fatalf("Error fetching events from database: %v", err)
	}
	eventUrlsSet := toSet(eventUrls)

	// add event to database if url does not exist
	for _, event := range data {
		if !eventUrlsSet[event.EventUrl] {
			eventData := map[string]interface{}{
				"event_name":        event.EventName,
				"event_description": event.EventDescription,
				"event_image_url":   event.EventImageUrl,
				"event_url":         event.EventUrl,
			}
			payload, _ := json.Marshal(eventData)
			req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("apikey", supabaseApiKey)
			req.Header.Set("Authorization", "Bearer "+ supabaseApiKey)
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
	}

	w.WriteHeader(http.StatusOK)
	log.Println("Scraped data received and added into database successfully")
}

func CallbackServer() {
	http.HandleFunc("/ticketmaster-scrape-callback", handleScrapeCallback)
	log.Println("Server listening on port 3002")
	log.Fatal(http.ListenAndServe(":3002", nil))
}