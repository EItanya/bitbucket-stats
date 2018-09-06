package api

import (
	"bitbucket/models"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// --------------------------------------
//
//  Bitbucket response API types
//
// --------------------------------------

// Response main JSON response from Bitbucket v1 API
type Response struct {
	Size          int    `json:"size"`
	Limit         int    `json:"limit"`
	IsLastPage    bool   `json:"isLastPage"`
	Start         int    `json:"start"`
	Filter        string `json:"filter"`
	NextPageStart int    `json:"nextPageStart"`
}

// ProjectResponse JSON form of single project
type ProjectResponse struct {
	Values []models.Project `json:"values"`
	Response
}

func (p *ProjectResponse) UnmarshalHTTP(resp *http.Response) error {
	byt, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(byt, p)
	return err
}

// RepoResponse JSON form of single project
type RepoResponse struct {
	Values []models.Repository `json:"values"`
	Response
}

func (r *RepoResponse) UnmarshalHTTP(resp *http.Response) error {
	byt, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(byt, r)
	return err
}

// FileResponse JSON form of single project
type FileResponse struct {
	Values models.Files `json:"values"`
	Response
}

func (f *FileResponse) UnmarshalHTTP(resp *http.Response) error {
	byt, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	err = json.Unmarshal(byt, f)
	return err
}

// ErrorResponse main JSON Error Response
type ErrorResponse struct {
	Errors []ErrorValue `json:"errors"`
}

// ErrorValue value of errors in Error response
type ErrorValue struct {
	Context       interface{} `json:"context"`
	Message       string      `json:"message"`
	ExceptionName interface{} `json:"exceptionName"`
}

// Entity interface for json http response
type Entity interface {
	UnmarshalHTTP(*http.Response) error
}

func GetEntity(r *http.Response, v Entity) error {
	return v.UnmarshalHTTP(r)
}

func readEntityFromResp(resp *http.Response, dat Entity) error {
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(byt, &dat)
}
