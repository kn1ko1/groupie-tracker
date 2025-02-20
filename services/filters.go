package services

import (
	"groupie-tracker/models"
	"strconv"
	"strings"
)

// FilterArtists applies filtering based on name, year, and location
func FilterArtists(artists []models.Artist, locations map[int][]string, name, year, location string) []models.Artist {
	filtered := []models.Artist{}

	for _, artist := range artists {
		nameMatch := name == "" || strings.EqualFold(artist.Name, name)
		yearMatch := year == "" || strconv.Itoa(artist.StartYear) == year
		locationMatch := location == ""

		if !locationMatch {
			for _, loc := range locations[artist.ID] {
				if loc == location {
					locationMatch = true
					break
				}
			}
		}

		if nameMatch && yearMatch && locationMatch {
			filtered = append(filtered, artist)
		}
	}

	return filtered
}
