package main

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/models"
	"io"
	"net/http"
	"sync"
	"time"
)

// Cache for locations
var (
	cachedLocations       map[int][]string
	lastLocationFetchTime time.Time
	locationCacheMutex    sync.Mutex
)

// Fetch locations with caching
func fetchLocationsCached(url string) (map[int][]string, error) {
	locationCacheMutex.Lock()
	defer locationCacheMutex.Unlock()

	// Check cache validity
	if time.Since(lastLocationFetchTime) < cacheDuration {
		fmt.Println("Returning cached location data")
		return cachedLocations, nil
	}

	// Fetch fresh data
	locations, err := fetchLocations(url)
	if err != nil {
		return nil, err
	}

	// Log only when fresh data is fetched
	fmt.Println("Fetched fresh location data from API...")

	// Update cache
	cachedLocations = locations
	lastLocationFetchTime = time.Now()

	return locations, nil
}

// Fetch locations from API
func fetchLocations(url string) (map[int][]string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Define API response structure
	type APIResponse struct {
		Index []models.LocationData `json:"index"`
	}

	var responseData APIResponse
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Convert to map[int][]string
	locations := make(map[int][]string)
	for _, entry := range responseData.Index {
		locations[entry.ID] = entry.Locations
	}

	return locations, nil
}
