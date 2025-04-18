package maps

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/labstack/echo/v4"
)

type AutoCompleteResponse struct {
	Predictions []struct {
		Description string `json:"description"`
	} `json:"predictions"`
	Status string `json:"status"`
}

func AutocompleteHandler(c echo.Context) error {
	input := c.QueryParam("input")
	if input == "" {
		return c.String(http.StatusBadRequest, "input parameter missing")
	}

	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return c.String(http.StatusInternalServerError, "API key not set")
	}

	baseURL := "https://maps.googleapis.com/maps/api/place/autocomplete/json"
	params := url.Values{}
	params.Add("input", input)
	params.Add("key", apiKey)

	urlStr := baseURL + "?" + params.Encode()
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Println("Error performing GET request:", err)
		return c.String(http.StatusInternalServerError, "Failed to fetch suggestions")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Non-200 response from Google API:", resp.Status)
		return c.String(resp.StatusCode, "Error from Google API")
	}

	var acResp AutoCompleteResponse
	if err := json.NewDecoder(resp.Body).Decode(&acResp); err != nil {
		log.Println("Error decoding response:", err)
		return c.String(http.StatusInternalServerError, "Error decoding Google response")
	}

	descriptions := make([]string, 0, len(acResp.Predictions))
	for _, prediction := range acResp.Predictions {
		descriptions = append(descriptions, prediction.Description)
	}

	return c.JSON(http.StatusOK, descriptions)
}

func MapsJSHandler(c echo.Context) error {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return c.String(http.StatusInternalServerError, "API key not set")
	}

	query := c.QueryParams()
	query.Set("key", apiKey)

	targetURL := "https://maps.googleapis.com/maps/api/js?" + query.Encode()
	resp, err := http.Get(targetURL)
	if err != nil {
		return c.String(http.StatusBadGateway, "Error reaching Google Maps API")
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		c.Response().Header().Set(k, v[0])
	}
	c.Response().WriteHeader(resp.StatusCode)
	io.Copy(c.Response().Writer, resp.Body)

	return nil
}
