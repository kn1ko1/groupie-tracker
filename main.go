package main

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/models"
	"groupie-tracker/services"
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
	// ✅ Debugging: Print every request
	fmt.Println("🔍 Processing Request:", r.URL.Path)

	// Ignore favicon requests (prevents unnecessary calls)
	if r.URL.Path == "/favicon.ico" {
		return
	}

	apiURL := "https://groupietrackers.herokuapp.com/api/artists"
	locationsAPI := "https://groupietrackers.herokuapp.com/api/locations"

	// Only fetch once (no duplicate calls)
	artists, err := fetchArtistsCached(apiURL)
	if err != nil {
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		return
	}

	// ✅ Only fetch once (no duplicate call)
	locations, err := fetchLocationsCached(locationsAPI)
	if err != nil {
		http.Error(w, "Failed to fetch locations", http.StatusInternalServerError)
		return
	}

	// Extract unique locations
	locationSet := make(map[string]bool)
	uniqueLocations := make([]string, 0)
	for _, locs := range locations {
		for _, loc := range locs {
			if !locationSet[loc] {
				locationSet[loc] = true
				uniqueLocations = append(uniqueLocations, loc)
			}
		}
	}

	// Extract and sort artist names
	artistNames := make([]string, 0, len(artists))
	startYears := make([]int, 0, len(artists))
	for _, artist := range artists {
		artistNames = append(artistNames, artist.Name)
		startYears = append(startYears, artist.StartYear)
	}

	sort.Strings(artistNames) // Sort names alphabetically
	sort.Ints(startYears)
	sort.Strings(uniqueLocations) // Sort locations alphabetically

	// Get filter values
	locationFilter := strings.TrimSpace(r.URL.Query().Get("location"))

	//Reverse the sorted years to make them descending
	uniqueYears := make([]int, 0)
	yearSet := make(map[int]bool)
	for i := len(startYears) - 1; i >= 0; i-- { // Reverse iteration for descending order
		year := startYears[i]
		if !yearSet[year] {
			yearSet[year] = true
			uniqueYears = append(uniqueYears, year)
		}
	}

	// Get filter values from query parameters
	nameFilter := strings.TrimSpace(r.URL.Query().Get("name"))
	yearFilter := strings.TrimSpace(r.URL.Query().Get("year"))
	//locationFilter := strings.TrimSpace(r.URL.Query().Get("location")) // Future concert filtering
	//fmt.Println(locationFilter)

	filteredArtists := services.FilterArtists(artists, locations, nameFilter, yearFilter, locationFilter)

	// Get sorting order from query parameter
	sortOrder := strings.TrimSpace(r.URL.Query().Get("sort"))
	//fmt.Println(1)

	// Sort artist names list for dropdown
	// if sortOrder == "desc" {
	// 	sort.Sort(sort.Reverse(sort.StringSlice(artistNames))) // Z-A order
	// 	fmt.Println(2)
	// } else {
	// 	sort.Strings(artistNames) // A-Z order (default)
	// 	//fmt.Println(3)
	// }

	// Ensure sorting applies to the main artist list too
	if sortOrder == "desc" {
		fmt.Println(4)
		// Ensure artists are sorted in descending order
		sort.Slice(filteredArtists, func(i, j int) bool {
			return filteredArtists[i].Name > filteredArtists[j].Name // Sort A-Z
		})

	} else if sortOrder == "asc" {
		fmt.Println(4)
		// Ensure artists are sorted in descending order
		sort.Slice(filteredArtists, func(i, j int) bool {
			return filteredArtists[i].Name < filteredArtists[j].Name // Sort A-Z
		})
	} else {
		// Reset to default order using ID (this fixes "Default" order)
		sort.SliceStable(filteredArtists, func(i, j int) bool {
			return filteredArtists[i].ID < filteredArtists[j].ID // Sort by ID in ascending order // artists[i].ID < ists[j].ID
		})
		// sort.Slice(artists, func(i, j int) bool {
		// 	return artists[i].Name < artists[j].Name //Maintain orginal order
		// } else {
		// 	// Ensure Default order resets properly
		// 	sort.SliceStable(artists, func(i, j int) bool {
		// 		return artists[i].Name < artists[j].Name //Maintain orginal order
		// 	})
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
		Title            string
		Artists          []models.Artist
		ShowDetail       bool // Now properly defined
		SortedNames      []string
		SortedYears      []int
		SortedLocations  []string
		Page             int
		TotalPages       int
		PrevPage         int
		NextPage         int
		SelectedName     string
		SelectedYear     string
		SelectedLocation string
		SelectedSort     string
	}{
		Title:            "Artists - Band Info",
		Artists:          paginatedArtists,
		ShowDetail:       isFullView, // Show full details only on `/artists`
		SortedNames:      artistNames,
		SortedYears:      uniqueYears, // Now sorted in descending order
		SortedLocations:  uniqueLocations,
		Page:             page,
		TotalPages:       totalPages,
		PrevPage:         prevPage,
		NextPage:         nextPage,
		SelectedName:     nameFilter,
		SelectedYear:     yearFilter,
		SelectedLocation: locationFilter,
		SelectedSort:     sortOrder, // ✅ Ensure sorting selection is retained
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
