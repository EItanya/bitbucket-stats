package api

import (
	"fmt"
	"os"
)

// SavedRepos is the format by which repos are saved
type SavedRepos []RepoModel

// Filter is the function to filter repos
func (data *SavedRepos) Filter(repos map[string][]string) {

}
func (data SavedRepos) filterRepos() {

}

var reposURLPath = func(projKey string) string {
	return fmt.Sprintf("/projects/%s/repos", projKey)
}

const reposFilePath = "data/repos.json"

// GetRepos get all repos from ***REMOVED*** Bitbucket
func (client *Client) GetRepos(repos map[string][]string) (SavedRepos, error) {
	var reposJSON SavedRepos
	var repoChan = make(chan []RepoModel)
	if _, err := os.Stat(reposFilePath); os.IsNotExist(err) {
		projectJSON, err := getProjectsJSON()
		if err != nil {
			return nil, err
		}
		for _, v := range projectJSON.Values {
			go client.getReposInternal(v, repoChan)
		}
		for range projectJSON.Values {
			reposJSON = append(reposJSON, <-repoChan...)
		}
		err = writeJSONToFile(&reposJSON, reposFilePath)
		if err != nil {
			return nil, err
		}
	} else {
		err := readJSONFromFile(reposFilePath, &reposJSON)
		if err != nil {
			return nil, err
		}
	}
	return reposJSON, nil
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
