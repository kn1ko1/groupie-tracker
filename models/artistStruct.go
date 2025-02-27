package models

type Artist struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Image      string   `json:"image"`
	StartYear  int      `json:"creationDate"`
	FirstAlbum string   `json:"firstAlbum"`
	Members    []string `json:"members"`
	Locations  []string `json:"locnames"` // New field with fake name to ignore
}
