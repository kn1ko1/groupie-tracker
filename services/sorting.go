package services

import (
	"groupie-tracker/models"
	"sort"
)

// SortArtists sorts the list of artists based on the sorting order
func SortArtists(artists []models.Artist, order string) {
	if order == "desc" {
		sort.SliceStable(artists, func(i, j int) bool {
			return artists[i].Name > artists[j].Name
		})
	} else if order == "asc" {
		sort.SliceStable(artists, func(i, j int) bool {
			return artists[i].Name < artists[j].Name
		})
	} else {
		sort.SliceStable(artists, func(i, j int) bool {
			return artists[i].ID < artists[j].ID // Reset to original order
		})
	}
}
