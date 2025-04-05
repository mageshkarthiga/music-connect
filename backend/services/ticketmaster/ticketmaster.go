package services

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "github.com/joho/godotenv"
    // Removed unused "bytes" import
)

// Event represents a ticketmaster event
type Event struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    URL  string `json:"url"`
}

// GetTicketmasterEvents fetches events from the Ticketmaster API
func GetTicketmasterEvents() ([]Event, error) {
    // Load environment variables from .env
    if err := godotenv.Load(); err != nil {
        return nil, fmt.Errorf("Error loading .env file")
    }

    apiKey := os.Getenv("TICKETMASTER_API_KEY")
    if apiKey == "" {
        return nil, fmt.Errorf("API key is missing")
    }

    // Define the endpoint for the Ticketmaster API
    endpoint := "https://app.ticketmaster.com/discovery/v2/events.json"

    // Set up query parameters (example: search for "music" events in San Francisco)
    params := "?keyword=music&city=San+Francisco&apikey=" + apiKey

    // Make the HTTP request
    resp, err := http.Get(endpoint + params)
    if err != nil {
        return nil, fmt.Errorf("Error making the request: %v", err)
    }
    defer resp.Body.Close()

    // Read the response body
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("Error reading the response body: %v", err)
    }

    // Check if we received a valid response (status code 200)
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("Error: Received non-OK response: %v", resp.StatusCode)
    }

    // Parse the response into a map (JSON format)
    var result map[string]interface{}
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, fmt.Errorf("Error unmarshalling response: %v", err)
    }

    // Check if the response contains "embedded" events
    var events []Event
    if _embedded, ok := result["_embedded"]; ok {
        eventsData, _ := _embedded.(map[string]interface{})["events"]
        if eventsData != nil {
            eventsArray := eventsData.([]interface{})
            for _, event := range eventsArray {
                eventData := event.(map[string]interface{})
                // Extract relevant details
                eventName := eventData["name"].(string)
                eventURL := eventData["url"].(string)

                events = append(events, Event{
                    Name: eventName,
                    URL:  eventURL,
                })
            }
        }
    }

    return events, nil
}
