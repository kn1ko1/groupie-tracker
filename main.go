package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
)

type Artists []struct {
	//This is our first struct created from API URL. Locations, ConcertDates and Relations all provide another URL and have to be dealt with later.
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Artist struct {
	//This struct becomes the new home of all of our info including the Locations, ConcertDates and Relations info from their respective URLs.
	ID           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Locations    []string            `json:"locations"`
	ConcertDates []string            `json:"concertDates"`
	Relations    map[string][]string `json:"datesLocations"`
}

type Locations struct {
	//This is our Locations struct.
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
	Dates     string   `json:"dates"`
}

type Dates struct {
	//This is our Dates struct.
	ID    int      `json:"id"`
	Dates []string `json:"dates"`
}

type Relations struct {
	//This is our Relations struct.
	ID                int                 `json:"id"`
	DatesandLocations map[string][]string `json:"datesLocations"`
}

func artistsConvertor(url string, object interface{}) error {
	//We have created this function to convert the information from the Artists JSON files. We provide a URL and an objest called responseArtists noted below in our main function.

	response, err := http.Get(url)
	//We name this variable response and use Get to get the information needed from the URL as well as dealing with an error.
	if err != nil {
		fmt.Println("Problem with artist conversion.")
		os.Exit(1)
		//If we have an error we print the above message and then exit. If no error then we continue on.
	}
	defer response.Body.Close()
	//This defer means that when this entire function completes its purpose we can come back to this position and close the function.

	responseData, err := io.ReadAll(response.Body)
	//We read the info from response using ReadAll and name it to the variable responseData.
	if err != nil {
		fmt.Println("Problem with artist conversion.")
		os.Exit(1)
		//If we have an error we print the above message and then exit. If no error then we continue on.
	}
	return json.Unmarshal(responseData, &object)
	//We return the unmarshalled data from the URL.
}

func locationConvertor(url string) []string {
	//We have created this function to convert the information from the Locations JSON files. We provide a URL and get a slice of string.

	var locations Locations
	//We name a variable locations.

	response, err := http.Get(url)
	//We name this variable response and use Get to get the information needed from the URL as well as dealing with an error.
	if err != nil {
		fmt.Println("Problem with location conversion.")
		os.Exit(1)
		//If we have an error we print the above message and then exit. If no error then we continue on.
	}
	defer response.Body.Close()
	//This defer means that when this entire function completes its purpose we can come back to this position and close the function.

	responseData, err := io.ReadAll(response.Body)
	//We read the info from response using ReadAll and name it to the variable responseData.
	if err != nil {
		fmt.Println("Problem with location conversion.")
		os.Exit(1)
		//If we have an error we print the above message and then exit. If no error then we continue on.
	}
	err = json.Unmarshal(responseData, &locations)
	//We are unmarshalling the locations information and putting it into the locations struct.
	if err != nil {
		fmt.Println("Problem with location conversion.")
		os.Exit(1)
		//If we have an error we print the above message and then exit. If no error then we continue on.
	}
	return locations.Locations
	//We are returning the locations information from our locations struct.
}

func datesConvertor(url string) []string {
	//We have created this function to convert the information from the Dates JSON files. We provide a URL and get a slice of string.

	var dates Dates
	//We name a variable dates.

	response, err := http.Get(url)
	//We name this variable response and use Get to get the information needed from the URL as well as dealing with an error.
	if err != nil {
		fmt.Println("Problem with dates conversion.")
		os.Exit(1)
		//If we have an error we print the above message and then exit. If no error then we continue on.
	}
	defer response.Body.Close()
	//This defer means that when this entire function completes its purpose we can come back to this position and close the function.

	responseData, err := io.ReadAll(response.Body)
	//We read the info from response using ReadAll and name it to the variable responseData.
	if err != nil {
		fmt.Println("Problem with dates conversion.")
		os.Exit(1)
		//If we have an error we print the above message and then exit. If no error then we continue on.
	}
	err = json.Unmarshal(responseData, &dates)
	//We are unmarshalling the dates information and putting it into the dates struct.
	if err != nil {
		fmt.Println("Problem with dates conversion.")
		os.Exit(1)
		//If we have an error we print the above message and then exit. If no error then we continue on.
	}
	return dates.Dates
	//We are returning the dates information from our dates struct.
}

func relationsConvertor(url string) map[string][]string {
	//We are returning the relations information from our relations struct.

	var relations Relations
	//We name a variable relations.

	response, err := http.Get(url)
	//We name this variable response and use Get to get the information needed from the URL as well as dealing with an error.
	if err != nil {
		fmt.Println("Problem with relations conversion.")
		os.Exit(1)
		//If we have an error we print the above message and then exit. If no error then we continue on.
	}
	defer response.Body.Close()
	//This defer means that when this entire function completes its purpose we can come back to this position and close the function.

	responseData, err := io.ReadAll(response.Body)
	//We read the info from response using ReadAll and name it to the variable responseData.
	if err != nil {
		fmt.Println("Problem with relations conversion.")
		os.Exit(1)
		//If we have an error we print the above message and then exit. If no error then we continue on.
	}
	err = json.Unmarshal(responseData, &relations)
	//We are unmarshalling the relations information and putting it into the relations struct.
	if err != nil {
		fmt.Println("Problem with relations conversion.")
		os.Exit(1)
		//If we have an error we print the above message and then exit. If no error then we continue on.
	}
	return relations.DatesandLocations
	//We are returning the relations information from our relations struct.
}

func fromJsonArtist(ja Artists) ([]Artist, error) {
	//We are providing this function with our original Artists struct however we are making a new copy and populating it with the information needed from locations, dates and relations.

	outputArtists := []Artist{}
	//We are naming our new struct to the variable name outputArtists. This is the structure of the struct.

	for _, eachArtist := range ja {
		//We are range loop through each artist populating this struct with the data we need.
		//We call up our locationConvertor, datesConvertor and relationsConvertor functions to get the necessary information.
		artist := Artist{
			//We are naming the information for our struct artist.
			ID:           eachArtist.ID,
			Image:        eachArtist.Image,
			Name:         eachArtist.Name,
			Members:      eachArtist.Members,
			CreationDate: eachArtist.CreationDate,
			FirstAlbum:   eachArtist.FirstAlbum,
			Locations:    locationConvertor(eachArtist.Locations),
			ConcertDates: datesConvertor(eachArtist.ConcertDates),
			Relations:    relationsConvertor(eachArtist.Relations),
		}

		outputArtists = append(outputArtists, artist)
		//We append the data from outputArtists, which is the structure and artist, which is the data. This is now the new value of outputArtists
	}
	return outputArtists, nil
	//We return outputArtists and the nil error message.
}

func showError(w http.ResponseWriter, message string, statusCode int) {
	//This funciton is created to show our 404 error message.
	t, err := template.ParseFiles("template/errors.html")
	//We parse the html error file and display.
	if err == nil {
		w.WriteHeader(statusCode)
		t.Execute(w, message)
		return
	}
}

func main() {

	fs := http.FileServer(http.Dir("./static"))
	//We are creating a variable named fs which is getting the informatin from the folder named static.
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	//We are serving the files in the statis folder. This folder caontains our CSS layout.

	var api string = "https://groupietrackers.herokuapp.com/api/artists"
	//The URL which contains all of our data has been named to the variable name api.

	var ResponseArtists Artists
	//We are copying the structure of the Artists struct to be named ResponseArtists and then populated in our main function below.

	artistsConvertor(api, &ResponseArtists)
	//We are taking the information from api and the location of ResponseArtists and putting them through the artistsConvertor.
	finalArtists, err := fromJsonArtist(ResponseArtists)
	//We now have the variable finalArtists which uses the function fromJsonArtists to populate our struct. This contains our final data set.
	if err != nil {
		os.Exit(1)
		//If the error does not equal nil then we exit.
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//This function writes to our webpage.
		if r.URL.Path != "/" {
			showError(w, "404 Page not found.", 404)
			//If our URL does not end with / then we call the showError function and show the 404 error.
			return
		}
		template, _ := template.ParseFiles("template/index.html")
		//We parse the html index file and display.
		template.Execute(w, finalArtists)
		//We execute the template with the information from finalArtists.
	})
	fmt.Printf("Starting server at port 8088\n")
	//Print this message to the terminal to signal that the webpage is ready.
	if err := http.ListenAndServe(":8088", nil); err != nil {
		fmt.Println("HTTP status 500 - Internal Server Errors")
		os.Exit(1)
		//If there are server errors then print the above message to the terminal and exit.
	}
}
