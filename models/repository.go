package models

// --------------------------------------
//
//  Bitbucket API repo types
//
// --------------------------------------

// RepoModel JSON form of single project
type Repository struct {
	Slug          string    `json:"slug"`
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	ScmID         string    `json:"scmId"`
	State         string    `json:"state"`
	StatusMessage string    `json:"statusMessage"`
	Forkable      bool      `json:"forkable"`
	Project       Project   `json:"project"`
	Public        bool      `json:"public"`
	Links         RepoLinks `json:"links"`
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
