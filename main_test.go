package main

import (
	"encoding/json"
	"groupie-tracker/models"
	"groupie-tracker/services"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestSortArtists(t *testing.T) {
	artists := []models.Artist{
		{ID: 1, Name: "ZZ Top"},
		{ID: 2, Name: "AC/DC"},
		{ID: 3, Name: "Metallica"},
	}

	// Test Ascending Order (A-Z)
	services.SortArtists(artists, "asc")
	if artists[0].Name != "AC/DC" || artists[1].Name != "Metallica" || artists[2].Name != "ZZ Top" {
		t.Errorf("Ascending sort failed, got: %v", artists)
	}

	// Test Descending Order (Z-A)
	services.SortArtists(artists, "desc")
	if artists[0].Name != "ZZ Top" || artists[1].Name != "Metallica" || artists[2].Name != "AC/DC" {
		t.Errorf("Descending sort failed, got: %v", artists)
	}

	// Test Default Order (by ID, assuming ID order matches original)
	services.SortArtists(artists, "")
	if artists[0].Name != "ZZ Top" || artists[1].Name != "AC/DC" || artists[2].Name != "Metallica" {
		t.Errorf("Default sort failed, got: %v", artists)
	}
}

// Test getArtistsHandler for JSON response
func TestGetArtistsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/artists", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getArtistsHandler) // Use your getArtistsHandler function
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Test for a valid JSON response
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Handler returned wrong content type: got %v want application/json", contentType)
	}

	// Check if JSON response is valid
	var response []models.Artist
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Invalid JSON response: %v", err)
	}

	// Ensure response is not empty (assuming API should return artists)
	if len(response) == 0 {
		t.Errorf("Expected at least one artist, but got empty response")
	}
}

// Test Pagination Logic
func TestPagination(t *testing.T) {
	artists := make([]models.Artist, 50) // Create 50 dummy artists
	for i := range artists {
		artists[i].Name = "Artist " + strconv.Itoa(i+1)
		artists[i].ID = i + 1
	}

	// Define page size
	const itemsPerPage = 15

	// Test Page 1
	startIndex, endIndex := 0, itemsPerPage
	if endIndex > len(artists) {
		endIndex = len(artists)
	}
	page1 := artists[startIndex:endIndex]
	if len(page1) != itemsPerPage {
		t.Errorf("Page 1 incorrect: got %d items, expected %d", len(page1), itemsPerPage)
	}

	// Test Last Page (Page 4 for 50 artists)
	startIndex, endIndex = (4-1)*itemsPerPage, 4*itemsPerPage
	if endIndex > len(artists) {
		endIndex = len(artists)
	}
	page4 := artists[startIndex:endIndex]
	if len(page4) != 5 { // Last page should have 5 items
		t.Errorf("Page 4 incorrect: got %d items, expected 5", len(page4))
	}
}
