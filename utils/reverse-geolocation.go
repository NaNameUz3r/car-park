package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var APIKEY = "r6YsPZAcZw218_dYQta0r5zo-Y2H8ItP4TSRN6iwuv4"

type output struct {
	Items []struct {
		Title string `json:"title"`
	} `json:"items"`
}

// 42.3824936,-83.0747712
func GeocodeToAddress(latitude, longitude float64) (address string) {
	url := "https://revgeocode.search.hereapi.com/v1/revgeocode?apiKey=" + APIKEY + "&at=" + fmt.Sprint(latitude) + "," + fmt.Sprint(longitude)

	res, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var data output
	json.Unmarshal(body, &data)
	// fmt.Println(data)
	for _, add := range data.Items {
		address = add.Title
	}
	return address
}
