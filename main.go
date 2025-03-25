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

//var templates = template.Must(template.ParseGlob("./frontend/templates/*.html"))

// var funcMap = template.FuncMap{
// 	"replaceSpaces": func(s string) string {
// 		return strings.ReplaceAll(s, " ", "-") // Convert spaces to hyphens
// 	},
// }

// var templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("./frontend/templates/*.html"))

// func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
// 	err := templates.ExecuteTemplate(w, tmpl, data)
// 	if err != nil {
// 		// Prevent multiple response writes
// 		http.Error(w, fmt.Sprintf("Unable to load template: %v", err), http.StatusInternalServerError)
// 		return
// 	}
// }

// func getArtistDetailsPage(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("🔍 Processing Request:", r.URL.Path)

// 	// Extract artist name from URL
// 	nameStr := strings.TrimPrefix(r.URL.Path, "/artist/")

// 	nameStr = strings.ReplaceAll(nameStr, "-", " ")
// 	nameStr = strings.TrimSpace(nameStr)

// 	fmt.Println("Extracted Artist Name from URL:", nameStr) // 🔍 Debugging

// 	// Fetch all artists
// 	apiURL := "https://groupietrackers.herokuapp.com/api/artists"
// 	artists, err := services.FetchArtistsCached(apiURL)
// 	if err != nil {
// 		http.Error(w, "Failed to fetch artist details", http.StatusInternalServerError)
// 		return
// 	}

// 	fmt.Println("Available artists:")
// 	for _, artist := range artists {
// 		fmt.Println(artist.Name) // Check if the artist exists in the list
// 	}

// 	var selectedArtist *models.Artist
// 	for _, artist := range artists {
// 		if strings.EqualFold(strings.TrimSpace(artist.Name), nameStr) {
// 			selectedArtist = &artist
// 			break
// 		}
// 	}

// 	// If artist not found, return error
// 	if selectedArtist == nil {
// 		http.Error(w, "Artist not found", http.StatusNotFound)
// 		return
// 	}

// 	// Render artist details template
// 	renderTemplate(w, "singleview.html", selectedArtist)
// }

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
