package handlers

import (
	"fmt"
	"groupie-tracker/models"
	"groupie-tracker/services"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// var templates = template.Must(template.ParseGlob("./frontend/templates/*.html"))

var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"add": func(a, b int) int { return a + b },
	"until": func(count int) []int {
		var i []int
		for j := 1; j <= count; j++ {
			i = append(i, j)
		}
		return i
	},
	"replaceSpaces": func(s string) string {
		return strings.ReplaceAll(s, " ", "-") // Convert spaces to hyphens
	},
}).ParseGlob("./frontend/templates/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to load template: %v", err), http.StatusInternalServerError)
		return
	}
}

func GetArtistsPage(w http.ResponseWriter, r *http.Request) {
	apiURL := "https://groupietrackers.herokuapp.com/api/artists"
	locationsAPI := "https://groupietrackers.herokuapp.com/api/locations"

	artists, err := services.FetchArtistsCached(apiURL)
	if err != nil {
		http.Error(w, "Failed to fetch artists", http.StatusInternalServerError)
		return
	}

	locations, err := services.FetchLocationsCached(locationsAPI)
	if err != nil {
		http.Error(w, "Failed to fetch locations", http.StatusInternalServerError)
		return
	}

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

	artistNames := make([]string, 0, len(artists))
	startYears := make([]int, 0, len(artists))

	for i, artist := range artists {
		artistNames = append(artistNames, artist.Name)
		startYears = append(startYears, artist.StartYear)

		if locs, exists := locations[artist.ID]; exists {
			artists[i].Locations = locs
		}
	}

	sort.Strings(artistNames)
	sort.Ints(startYears)
	sort.Strings(uniqueLocations)

	locationFilter := strings.TrimSpace(r.URL.Query().Get("location"))
	nameFilter := strings.TrimSpace(r.URL.Query().Get("name"))
	yearFilter := strings.TrimSpace(r.URL.Query().Get("year"))

	uniqueYears := make([]int, 0)
	yearSet := make(map[int]bool)
	for i := len(startYears) - 1; i >= 0; i-- {
		year := startYears[i]
		if !yearSet[year] {
			yearSet[year] = true
			uniqueYears = append(uniqueYears, year)
		}
	}

	var filteredArtists []models.Artist
	sortOrder := strings.TrimSpace(r.URL.Query().Get("sort"))
	searchQuery := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("search")))

	//var filteredArtists []models.Artist
	if len(searchQuery) >= 2 { // Only search if the query is at least 2 characters
		for _, artist := range artists {
			// Match on artist name
			if strings.Contains(strings.ToLower(artist.Name), searchQuery) {
				filteredArtists = append(filteredArtists, artist)
				continue
			}

			// Match on member names
			for _, member := range artist.Members {
				if strings.Contains(strings.ToLower(member), searchQuery) {
					filteredArtists = append(filteredArtists, artist)
					break
				}
			}

			// Match on concert locations
			for _, loc := range artist.Locations {
				if strings.Contains(strings.ToLower(loc), searchQuery) {
					filteredArtists = append(filteredArtists, artist)
					break
				}
			}

			// Match on first album (e.g. "2005-04-01")
			if strings.Contains(strings.ToLower(artist.FirstAlbum), searchQuery) ||
				strings.Contains(strconv.Itoa(artist.StartYear), searchQuery) {
				filteredArtists = append(filteredArtists, artist)
				continue
			}
		}
	} else {
		// Match on creation year
		filteredArtists = services.FilterArtists(artists, locations, nameFilter, yearFilter, locationFilter)
	}

	services.SortArtists(filteredArtists, sortOrder)

	const itemsPerPage = 15
	page := 1
	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}

	paginatedArtists, prevPage, nextPage := services.PaginateArtists(filteredArtists, page, itemsPerPage)
	totalPages := (len(filteredArtists) + itemsPerPage - 1) / itemsPerPage

	isFullView := r.URL.Path == "/artists"

	data := struct {
		Title            string
		Artists          []models.Artist
		ShowDetail       bool
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
		ShowDetail:       isFullView,
		SortedNames:      artistNames,
		SortedYears:      uniqueYears,
		SortedLocations:  uniqueLocations,
		Page:             page,
		TotalPages:       totalPages,
		PrevPage:         prevPage,
		NextPage:         nextPage,
		SelectedName:     nameFilter,
		SelectedYear:     yearFilter,
		SelectedLocation: locationFilter,
		SelectedSort:     sortOrder,
	}

	renderTemplate(w, "artists.html", data)
}
