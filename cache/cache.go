package cache

import (
	"fmt"
	"groupie-tracker/models"
	"sync"
	"time"
)

// Cache variables
var (
	cachedArtists         []models.Artist
	cachedLocations       map[int][]string
	lastFetchTime         time.Time
	lastLocationFetchTime time.Time
	cacheMutex            sync.Mutex
	cacheDuration         = 5 * time.Minute
)

// Fetch artists with caching
func FetchArtistsCached(fetchFunc func() ([]models.Artist, error)) ([]models.Artist, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if time.Since(lastFetchTime) < cacheDuration {
		fmt.Println("✅ Returning cached artist data")
		return cachedArtists, nil
	}

	artists, err := fetchFunc()
	if err != nil {
		return nil, err
	}

	cachedArtists = artists
	lastFetchTime = time.Now()
	return artists, nil
}

// Fetch locations with caching
func FetchLocationsCached(fetchFunc func() (map[int][]string, error)) (map[int][]string, error) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if time.Since(lastLocationFetchTime) < cacheDuration {
		fmt.Println("✅ Returning cached location data")
		return cachedLocations, nil
	}

	locations, err := fetchFunc()
	if err != nil {
		return nil, err
	}

	cachedLocations = locations
	lastLocationFetchTime = time.Now()
	return locations, nil
}
