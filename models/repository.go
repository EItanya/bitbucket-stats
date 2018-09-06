package models

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
