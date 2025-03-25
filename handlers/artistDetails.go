package handlers

import (
	"fmt"
	"groupie-tracker/models"
	"groupie-tracker/services"
	"net/http"
	"strings"
)

func GetArtistDetailsPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("🔍 Processing Request:", r.URL.Path)

	// Extract artist name from URL
	nameStr := strings.TrimPrefix(r.URL.Path, "/artist/")

	nameStr = strings.ReplaceAll(nameStr, "-", " ")
	nameStr = strings.TrimSpace(nameStr)

	fmt.Println("Extracted Artist Name from URL:", nameStr) // 🔍 Debugging

	// Fetch all artists
	apiURL := "https://groupietrackers.herokuapp.com/api/artists"
	artists, err := services.FetchArtistsCached(apiURL)
	if err != nil {
		http.Error(w, "Failed to fetch artist details", http.StatusInternalServerError)
		return
	}

	fmt.Println("Available artists:")
	for _, artist := range artists {
		fmt.Println(artist.Name) // Check if the artist exists in the list
	}

	var selectedArtist *models.Artist
	for _, artist := range artists {
		if strings.EqualFold(strings.TrimSpace(artist.Name), nameStr) {
			selectedArtist = &artist
			break
		}
	}

	// If artist not found, return error
	if selectedArtist == nil {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	// Render artist details template
	renderTemplate(w, "singleview.html", selectedArtist)
}
