package services

import (
	"groupie-tracker/models"
)

// PaginatedArtists applies pagination to the artist list.
func PaginateArtists(artists []models.Artist, page, itemsPerPage int) ([]models.Artist, int, int) {
	totalPages := (len(artists) + itemsPerPage - 1) / itemsPerPage

	// Ensure the page number is within bounds
	if page < 1 {
		page = 1
	} else if page > totalPages {
		page = totalPages
	}

	startIndex := (page - 1) * itemsPerPage
	endIndex := startIndex + itemsPerPage
	if endIndex > len(artists) {
		endIndex = len(artists)
	}

	paginatedArtists := artists[startIndex:endIndex]

	// Precompute previous and next page numbers
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}
	nextPage := page + 1
	if nextPage > totalPages {
		nextPage = totalPages
	}

	return paginatedArtists, prevPage, nextPage
}
