package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func removeAllLocalData() error {
	err := removeLocalFilesData()
	if err != nil {
		return err
	}
	err = removeLocalReposData()
	if err != nil {
		return err
	}
	err = removeLocalProjectsData()
	if err != nil {
		return err
	}
	return nil
}

func writeJSONToFile(data interface{}, filename string) error {
	if !strings.Contains(filename, ".json") {
		return errors.New("write file must be a valid json filename")
	}
	byt, err := json.MarshalIndent(&data, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, byt, 0600)
	if err != nil {
		return err
	}
	return nil
}

func readJSONFromFile(filename string, dataType interface{}) error {
	if !strings.Contains(filename, ".json") {
		return errors.New("Read file must be a valid json filename")
	}
	byt, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byt, dataType)
	if err != nil {
		return err
	}
	return nil
}

func readJSONFromResp(resp *http.Response, dat interface{}) error {
	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(byt, &dat)
}

func getProjectsJSON() (SavedProjects, error) {
	var projectJSON SavedProjects
	err := readJSONFromFile(projectsFilePath, &projectJSON)
	return projectJSON, err
}

func getReposJSON() (SavedRepos, error) {
	var repoJSON SavedRepos
	err := readJSONFromFile(reposFilePath, &repoJSON)
	return repoJSON, err
}
