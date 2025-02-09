package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Artist struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	StartYear  int      `json:"creationDate"`
	FirstAlbum string   `json:"firstAlbum"`
	Members    []string `json:"members"`
}

// Cache variables
var (
	cachedArtists         []Artist
	cachedLocations       map[int][]string
	lastFetchTime         time.Time
	lastLocationFetchTime time.Time
	cacheMutex            sync.Mutex
	cacheDuration         = 5 * time.Minute // Cache validity duration
)

// Fetch artists with caching
func fetchArtistsCached(url string) ([]Artist, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Check if the cache is still valid
	if time.Since(lastFetchTime) < cacheDuration {
		fmt.Println("Returning cached data")
		return cachedArtists, nil
	}

	// Fetch fresh data
	artists, err := fetchArtists(url)
	if err != nil {
		return nil, err
	}

	// Update cache
	cachedArtists = artists
	lastFetchTime = time.Now()

	return artists, nil
}

// Fetch artists from API
func fetchArtists(url string) ([]Artist, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artists: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var artists []Artist
	if err := json.Unmarshal(body, &artists); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", response.Status)
	}

	return artists, nil
}

// Fetch Concert Locations with Caching
func fetchLocationsCached(url string) (map[int][]string, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Check if cache is still valid
	if time.Since(lastLocationFetchTime) < cacheDuration {
		return cachedLocations, nil
	}

	// Fetch fresh location data
	locations, err := fetchLocations(url)
	if err != nil {
		return nil, err
	}

	// Update cache
	cachedLocations = locations
	lastLocationFetchTime = time.Now()

	return locations, nil
}

// Fetch Concert Locations from API
// func fetchLocations(url string) (map[int][]string, error) {
// 	response, err := http.Get(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to fetch locations: %v", err)
// 	}
// 	defer response.Body.Close()

// 	body, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read response body: %v", err)
// 	}

// 	// Define struct to match API response
// 	var locationsData map[string]struct {
// 		DatesLocations map[string][]string `json:"datesLocations"`
// 	}

// 	if err := json.Unmarshal(body, &locationsData); err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
// 	}

// 	// Convert to map[int][]string for filtering
// 	locations := make(map[int][]string)
// 	for artistID, locData := range locationsData {
// 		id, err := strconv.Atoi(artistID)
// 		if err != nil {
// 			continue
// 		}

// 		// Extract only the location names
// 		locations[id] = make([]string, 0)
// 		for loc := range locData.DatesLocations {
// 			locations[id] = append(locations[id], loc)
// 		}
// 	}

// 	return locations, nil
// }

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

	// Define struct to match API response
	type LocationData struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
	}

	type APIResponse struct {
		Index []LocationData `json:"index"`
	}

	// Unmarshal JSON into the struct
	var responseData APIResponse
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	// Convert to map[int][]string for filtering
	locations := make(map[int][]string)
	for _, entry := range responseData.Index {
		locations[entry.ID] = entry.Locations
	}

	// Debugging: Print the locations map
	//fmt.Println("Fetched locations:", locations)

	return locations, nil
}
