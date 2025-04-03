package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type AutoCompleteResponse struct {
	Predictions []struct {
		Description string `json:"description"`
	} `json:"predictions"`
	Status string `json:"status"`
}

func autocompleteHandler(w http.ResponseWriter, r *http.Request) {
	// Allow cross-origin requests if needed
	w.Header().Set("Access-Control-Allow-Origin", "*")

	input := r.URL.Query().Get("input")
	if input == "" {
		http.Error(w, "input parameter missing", http.StatusBadRequest)
		return
	}

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		http.Error(w, "API key not set", http.StatusInternalServerError)
		return
	}

	baseURL := "https://maps.googleapis.com/maps/api/place/autocomplete/json"
	params := url.Values{}
	params.Add("input", input)
	params.Add("key", apiKey)

	urlStr := baseURL + "?" + params.Encode()
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Println("Error performing GET request:", err)
		http.Error(w, "Failed to fetch suggestions", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Println("GET", urlStr)
	log.Println("Response status:", resp.Status)
	log.Println("Response headers:", resp.Header)

	if resp.StatusCode != http.StatusOK {
		log.Println("Non-200 response from Google API:", resp.Status)
		http.Error(w, "Error from Google API", resp.StatusCode)
		return
	}

	var acResp AutoCompleteResponse
	if err := json.NewDecoder(resp.Body).Decode(&acResp); err != nil {
		log.Println("Error decoding response:", err)
		http.Error(w, "Error decoding Google response", http.StatusInternalServerError)
		return
	}

	// Extract the descriptions from predictions and return as a list of strings.
	descriptions := make([]string, 0, len(acResp.Predictions))
	for _, prediction := range acResp.Predictions {
		descriptions = append(descriptions, prediction.Description)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(descriptions); err != nil {
		log.Println("Error encoding response:", err)
	}
}

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: .env file not found, relying on environment variables")
	}

	http.HandleFunc("/api/autocomplete", autocompleteHandler)
	log.Println("Server starting on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
