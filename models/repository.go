package models

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// --------------------------------------
//
//  Bitbucket API repo types
//
// --------------------------------------

// RepoModel JSON form of single project
type Repository struct {
	Slug          string  `json:"slug"`
	ID            int     `json:"id"`
	Name          string  `json:"name"`
	ScmID         string  `json:"scmId"`
	State         string  `json:"state"`
	StatusMessage string  `json:"statusMessage"`
	Forkable      bool    `json:"forkable"`
	Project       Project `json:"project"`
	Public        bool    `json:"public"`
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

func (r *Repository) Unmarshal(dat interface{}) error {
	if cast, ok := dat.(*Repository); ok {
		*cast = *r
		return nil
	}
	return fmt.Errorf("Improper type (%s) was passed into unmarshal method for Project model", dat)
}

func (r *Repository) Marshal(dat interface{}) error {
	switch typedData := dat.(type) {
	case Repository:
		*r = typedData
	case map[string]interface{}:
		err := mapstructure.Decode(typedData, r)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("Improper type (%s) was passed into marshal method for Repository model", dat)
	}
	return nil
}

// FilterRepos function to filter repos
func FilterRepos(data *[]Repository, repos []string) []Repository {
	if len(repos) == 0 {
		return *data
	}
	filteredRepos := make([]Repository, 0)
	ch := make(chan []Repository)
	for _, val := range repos {
		filterRepos(val, data, ch)
	}
	for range repos {
		filteredRepos = append(filteredRepos, <-ch...)
	}
	return filteredRepos

}
func filterRepos(val string, r *[]Repository, ch chan []Repository) {
	rm := make([]Repository, 0)
	for _, v := range *r {
		if v.Slug == val {
			rm = append(rm, v)
		}
	}
	ch <- rm
}
