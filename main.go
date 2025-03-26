package main

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/handlers"
	"groupie-tracker/services"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func getArtistsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	artistID := vars["id"]
	apiURL := "https://groupietrackers.herokuapp.com/api/artists"
	// Find the artist with the given ID
	artists, err := services.FetchArtistsCached(apiURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching artists: %v", err), http.StatusInternalServerError)
		return
	}
	// Render the singleview template with the artist data
	tmpl, err := template.ParseFiles("singleview.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, artistID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artists)
}

// Main function to set up the server
func main() {
	// Serve static files for styles
	http.Handle("/styles.css", http.FileServer(http.Dir("./")))

	// Define routes
	http.HandleFunc("/", handlers.GetArtistsPage)        // Default route renders the artists page
	http.HandleFunc("/artists", handlers.GetArtistsPage) // JSON API endpoint (optional)
	http.HandleFunc("/artist/", handlers.GetArtistDetailsPage)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
