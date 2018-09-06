package api

import (
	"fmt"
	"os"

	"github.com/gosuri/uiprogress"
)

// SavedRepos is the format by which repos are saved
type SavedRepos []RepoModel

// Filter is the function to filter repos
func (data SavedRepos) Filter(repos []string) []RepoModel {
	if len(repos) == 0 {
		return data
	}
	filteredRepos := make([]RepoModel, 0)
	ch := make(chan []RepoModel)
	for _, val := range repos {
		go data.filterRepos(val, data, ch)
	}
	for range repos {
		filteredRepos = append(filteredRepos, <-ch...)
	}
	return filteredRepos

}
func (data SavedRepos) filterRepos(val string, r []RepoModel, ch chan []RepoModel) {
	rm := make([]RepoModel, 0)
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
func (client *Client) GetRepos(repos []string) (*[]RepoModel, error) {
	var reposJSON SavedRepos
	var repoChan = make(chan []RepoModel)
	if _, err := os.Stat(reposFilePath); os.IsNotExist(err) {
		projectJSON, err := getProjectsJSON()
		if err != nil {
			return nil, err
		}
		numProjects := len(projectJSON.Values)
		bar := uiprogress.AddBar(numProjects + 1)
		bar.AppendCompleted()
		bar.PrependFunc(func(b *uiprogress.Bar) string {
			if b.Current() > numProjects {
				return fmt.Sprintf("Saving repo data:  %s", b.TimeElapsedString())
			}
			return fmt.Sprintf("Dowloading repo data:  %s", b.TimeElapsedString())
		})
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

func (client *Client) getReposInternal(v ProjectModel, c chan []RepoModel) {
	var repoJSON RepoResponse
	resp, _ := client.api.Get(reposURLPath(v.Key), 250)
	_ = readJSONFromResp(resp, &repoJSON)
	c <- repoJSON.Values
}

func removeLocalReposData() error {
	return os.Remove(reposFilePath)
}
