package api

import (
	"encoding/json"
	"fmt"
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
	Values []ProjectModel `json:"values"`
	Response
}

func (p *ProjectResponse) UnmarshalHTTP(resp *http.Response) error {
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byt, p)
	return err
}

// ProjectModel JSON form of single project
type ProjectModel struct {
	Key         string `json:"key"`
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	Type        string `json:"type"`
}

// RepoResponse JSON form of single project
type RepoResponse struct {
	Values []RepoModel `json:"values"`
	Response
}

func (r *RepoResponse) UnmarshalHTTP(resp *http.Response) error {
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byt, r)
	fmt.Println(r)
	return err
}

// RepoModel JSON form of single project
type RepoModel struct {
	Slug          string       `json:"slug"`
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	ScmID         string       `json:"scmId"`
	State         string       `json:"state"`
	StatusMessage string       `json:"statusMessage"`
	Forkable      bool         `json:"forkable"`
	Project       ProjectModel `json:"project"`
	Public        bool         `json:"public"`
	Links         RepoLinks    `json:"links"`
}

// RepoLinks JSON form of single project
type RepoLinks struct {
	Clone []Link `json:"clone"`
	Self  []Link `json:"self"`
}

// Link JSON form of single project
type Link struct {
	Href string `json:"href"`
	Name string `json:"name"`
}

// FileResponse JSON form of single project
type FileResponse struct {
	Values []string `json:"values"`
	Response
}

func (f *FileResponse) UnmarshalHTTP(resp *http.Response) error {
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byt, f)
	fmt.Println(f)
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
