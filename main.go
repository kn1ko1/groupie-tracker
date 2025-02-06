package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

var templates = template.Must(template.ParseGlob("./frontend/templates/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		// Prevent multiple response writes
		http.Error(w, fmt.Sprintf("Unable to load template: %v", err), http.StatusInternalServerError)
		return
	}
}

func getArtistsPage(w http.ResponseWriter, r *http.Request) {
	apiURL := "https://groupietrackers.herokuapp.com/api/artists"
	artists, err := fetchArtistsCached(apiURL)
	if err != nil {
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		return
	}

	// Fetch locations
	locationsAPI := "https://groupietrackers.herokuapp.com/api/locations"
	locations, err := fetchLocationsCached(locationsAPI)
	if err != nil {
		http.Error(w, "Failed to fetch locations", http.StatusInternalServerError)
		return
	}

	// Extract and sort artist names
	artistNames := make([]string, 0, len(artists))
	startYears := make([]int, 0, len(artists))
	locationSet := make(map[string]bool)
	uniqueLocations := make([]string, 0)

	for _, artist := range artists {
		artistNames = append(artistNames, artist.Name)
		startYears = append(startYears, artist.StartYear)

		// Store unique locations
		for _, loc := range locations[artist.ID] {
			if !locationSet[loc] {
				locationSet[loc] = true
				uniqueLocations = append(uniqueLocations, loc)
			}
		}
	}

	sort.Strings(artistNames) // Sort names alphabetically
	// sort.Sort(sort.Reverse(sort.IntSlice(startYears))) // Sort years in descending order

	// Sort years in ascending order first
	sort.Ints(startYears)

	// Reverse the sorted years to make them descending
	uniqueYears := make([]int, 0)
	yearSet := make(map[int]bool)
	for i := len(startYears) - 1; i >= 0; i-- { // Reverse iteration for descending order
		year := startYears[i]
		if !yearSet[year] {
			yearSet[year] = true
			uniqueYears = append(uniqueYears, year)
		}
	}

	//SortedLocations: uniqueLocations, // ✅ Ensure locations are passed to the template

	sort.Strings(uniqueLocations) // Sort locations alphabetically

	// Get filter values from query parameters
	nameFilter := strings.TrimSpace(r.URL.Query().Get("name"))
	yearFilter := strings.TrimSpace(r.URL.Query().Get("year"))
	locationFilter := strings.TrimSpace(r.URL.Query().Get("location")) // Future concert filtering

	// Filter artists basedon user selection
	filteredArtists := make([]Artist, 0)
	for _, artist := range artists {
		nameMatch := nameFilter == "" || artist.Name == nameFilter
		yearMatch := yearFilter == "" || strconv.Itoa(artist.StartYear) == yearFilter
		locationMatch := locationFilter == ""

		// Check if artist performed at the selected location
		if !locationMatch {
			for _, loc := range locations[artist.ID] {
				if loc == locationFilter {
					locationMatch = true
					break
				}
			}
		}

		if nameMatch && yearMatch {
			filteredArtists = append(filteredArtists, artist)
		}
	}

	// Pagination logic (5 columns x 3 rows = 15 artists per page)
	const itemsPerPage = 15
	page := 1
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}

	startIndex := (page - 1) * itemsPerPage
	endIndex := startIndex + itemsPerPage
	if endIndex > len(filteredArtists) {
		endIndex = len(filteredArtists)
	}

	paginatedArtists := filteredArtists[startIndex:endIndex]

	// Calculate total pages
	totalPages := (len(filteredArtists) + itemsPerPage - 1) / itemsPerPage

	// Precompute previous and next page numbers
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}
	nextPage := page + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	// Determine if full view (/artists) or minimal view (/)
	isFullView := r.URL.Path == "/artists"

	// Define template data with dynamic content for home ('/home') vs Artists ('/artists')
	data := struct {
		Title       string
		Artists     []Artist
		ShowDetail  bool // Now properly defined
		SortedNames []string
		SortedYears []int
		Page        int
		TotalPages  int
		PrevPage    int
		NextPage    int
	}{
		Title:       "Artists - Band Info",
		Artists:     paginatedArtists,
		ShowDetail:  isFullView, // Show full details only on `/artists`
		SortedNames: artistNames,
		SortedYears: uniqueYears, // Now sorted in descending order
		Page:        page,
		TotalPages:  totalPages,
		PrevPage:    prevPage,
		NextPage:    nextPage,
	}
	renderTemplate(w, "artists.html", data)
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
	http.Handle("/styles.css", http.FileServer(http.Dir("./")))

	// Define routes
	http.HandleFunc("/", getArtistsPage)        // Default route renders the artists page
	http.HandleFunc("/artists", getArtistsPage) // JSON API endpoint (optional)

	fmt.Println("Server running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
