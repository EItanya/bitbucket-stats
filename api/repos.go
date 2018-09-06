package api

import (
	"bitbucket/models"
	"fmt"
	"os"

	"github.com/gosuri/uiprogress"
)

// SavedRepos is the format by which repos are saved
type SavedRepos []models.Repository

// Filter is the function to filter repos
func (data SavedRepos) Filter(repos []string) []models.Repository {
	if len(repos) == 0 {
		return data
	}
	filteredRepos := make([]models.Repository, 0)
	ch := make(chan []models.Repository)
	for _, val := range repos {
		go data.filterRepos(val, data, ch)
	}
	for range repos {
		filteredRepos = append(filteredRepos, <-ch...)
	}
	return filteredRepos

}
func (data SavedRepos) filterRepos(val string, r []models.Repository, ch chan []models.Repository) {
	rm := make([]models.Repository, 0)
	for _, v := range r {
		if v.Slug == val {
			rm = append(rm, v)
		}
	}
	ch <- rm
}

var reposURLPath = func(projKey string) string {
	return fmt.Sprintf("/projects/%s/repos", projKey)
}

const reposFilePath = "data/repos.json"

// GetRepos get all repos from Bitbucket
func (client *Client) GetRepos(repos []string) (*[]models.Repository, error) {
	var reposJSON SavedRepos
	var repoChan = make(chan []models.Repository)
	if _, err := os.Stat(reposFilePath); os.IsNotExist(err) {
		projectJSON, err := getProjectsJSON()
		if err != nil {
			return nil, err
		}
		numProjects := len(projectJSON.Values)
		bar := uiprogress.AddBar(numProjects + 1)
		bar.AppendCompleted()
		bar.PrependFunc(prependFormatFunc(func(b *uiprogress.Bar) string {
			if b.Current() == b.Total {
				return "Repo data retrieved"
			} else if b.Current() >= numProjects {
				return "Saving repo data"
			}
			return "Dowloading repo data"
		}))
		for _, v := range projectJSON.Values {
			go client.getReposInternal(v, repoChan)
		}
		for range projectJSON.Values {
			reposJSON = append(reposJSON, <-repoChan...)
			bar.Incr()
		}
		err = writeJSONToFile(&reposJSON, reposFilePath)
		if err != nil {
			return nil, err
		}
		bar.Incr()
	} else {
		err := readJSONFromFile(reposFilePath, &reposJSON)
		if err != nil {
			return nil, err
		}
	}
	result := reposJSON.Filter(repos)
	return &result, nil
}

func (client *Client) getReposInternal(v models.Project, c chan []models.Repository) {
	var repoJSON RepoResponse
	opts := urlOptions{
		limit: 500,
	}
	resp, _ := client.api.Get(reposURLPath(v.Key), opts)
	_ = readJSONFromResp(resp, &repoJSON)
	c <- repoJSON.Values
}

func removeLocalReposData() error {
	return os.Remove(reposFilePath)
}
