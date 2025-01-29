package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

var templates = template.Must(template.ParseGlob("./frontend/templates/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		fmt.Printf("Error rendering template '%s': %v\n", tmpl, err)
		http.Error(w, "Unable to load template", http.StatusInternalServerError)
	}
}

func getArtistsPage(w http.ResponseWriter, r *http.Request) {
	apiURL := "https://groupietrackers.herokuapp.com/api/artists"
	artists, err := fetchArtistsCached(apiURL)
	if err != nil {
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		return
	}

	// Get filter values from query parameters
	nameFilter := r.URL.Query().Get("name")
	yearFilter := r.URL.Query().Get("year")

	// Filter artists
	filteredArtists := make([]Artist, 0)
	if len(artist) == 0{
		fmt.Println("No artists to filter")
		return
	}

	for _, artist := range artists {
		// Check both filters
		nameMatch := nameFilter == "" || containsIgnoreCase(artist.Name, nameFilter)
		yearMatch := yearFilter == "" || strconv.Itoa(artist.StartYear) == yearFilter

		if nameMatch && yearMatch {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	data := struct {
		Title   string
		Artists []Artist
	}{
		Title:   "Artists - Band Info",
		Artists: artists,
	}
	//renderTemplate(w, "artists.html", data)
}

// Helper function for case-insensitive substring matching
func containsIgnoreCase(str, substr string) bool {
	//fmt.Printf("Comparing '%s' with '%s'\n", strings.ToLower(str), strings.ToLower(substr))
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

func getArtistsHandler(w http.ResponseWriter, r *http.Request) {
	apiURL := "https://groupietrackers.herokuapp.com/api/artists"
	artists, err := fetchArtistsCached(apiURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching artists: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
}

// Main function to set up the server
func main() {
	// Serve static files for styles
	http.Handle("/styles.css", http.FileServer(http.Dir("./frontend")))

	// Define routes
	http.HandleFunc("/", getArtistsPage)          // Default route renders the artists page
	http.HandleFunc("artists", getArtistsHandler) // JSON API endpoint (optional)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
