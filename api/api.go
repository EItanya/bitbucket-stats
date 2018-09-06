package api

import (
	"fmt"
	"log"
	"net/http"
)

// API basis for API class
type API struct {
	BaseURL  string
	username string
	password string
	http.Client
}

const batchNumber = 10

func (api *API) creds() (string, string) {
	return api.username, api.password
}

// Get wrapper for http GET
func (api *API) Get(path string, limit int) (*http.Response, error) {
	link := fmt.Sprintf("%s%s?limit=%d", api.BaseURL, path, limit)
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(api.creds())
	req.Close = true
	resp, err := api.Do(req)
	return resp, err
}
