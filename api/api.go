package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

// API basis for API class
type API struct {
	BaseURL string
	user    UserInfo
	http.Client
}

const batchNumber = 10

var errBadRequest = errors.New("Bad Request attempted")
var errAuthorization = errors.New("Authorization error")
var errNotFound = errors.New("Resource Not Found")

func (api *API) creds() (string, string) {
	return api.user.Username, api.user.Password
}

func (api *API) doExt(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(api.creds())
	req.Close = true
	resp, err := api.Do(req)
	if err != nil {
		return nil, err
	}
	switch {
	case resp.StatusCode == 200:
		return resp, nil
	case resp.StatusCode >= 400 && resp.StatusCode <= 499:
		return resp, authErrors(resp)
	case resp.StatusCode >= 500:
		return resp, errors.New("500 error, I don't even know how you managed that\nGive it a couple minutes")
	default:
		return resp, err
	}
}

// Get wrapper for http GET
func (api *API) Get(path string, limit int) (*http.Response, error) {
	link := fmt.Sprintf("%s%s?limit=%d", api.BaseURL, path, limit)
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Fatal(err)
	}
	return api.doExt(req)
}
