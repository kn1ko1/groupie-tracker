package services

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/models"
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
	//Locations  []string `json:"locations"`
}

// Cache variables
var (
	cachedArtists []models.Artist
	lastFetchTime time.Time
	cacheMutex    sync.Mutex
	cacheDuration = 5 * time.Minute
)

// Fetch artists with caching
func FetchArtistsCached(url string) ([]models.Artist, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Check if the cache is still valid
	if time.Since(lastFetchTime) < cacheDuration {
		//fmt.Println("Returning cached data")
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
func fetchArtists(url string) ([]models.Artist, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch artists: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var artists []models.Artist
	if err := json.Unmarshal(body, &artists); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", response.Status)
	}

	return artists, nil
}
