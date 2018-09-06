package api

import (
	"fmt"
	"os"
)

// SavedProjects is the format by which projects are saved
type SavedProjects ProjectResponse

// Filter method to filter saved projects
func (data *SavedProjects) Filter(projects []string) []ProjectModel {
	if len(projects) == 0 {
		return data.Values
	}
	filteredProjects := make([]ProjectModel, 0)
	ch := make(chan []ProjectModel)
	for _, val := range projects {
		go data.filterProjects(val, data.Values, ch)
	}
	for range projects {
		filteredProjects = append(filteredProjects, <-ch...)
	}
	return filteredProjects
}

func (data SavedProjects) filterProjects(val string, p []ProjectModel, ch chan []ProjectModel) {
	pm := make([]ProjectModel, 0)
	for _, v := range p {
		if v.Key == val || string(v.ID) == val {
			pm = append(pm, v)
		}
	}
	ch <- pm
}

const projectsURLPath = "/projects"

var projectsURLPaths = func(projects []string) []string {
	if len(projects) == 0 {
		return []string{projectsURLPath}
	}
	urls := make([]string, 0)
	for _, val := range projects {
		urls = append(urls, fmt.Sprintf("%s/%s", projectsURLPath, val))
	}
	return urls
}

const projectsFilePath = "data/projects.json"

// GetProjects get all projects from ***REMOVED*** Bitbucket
func (client *Client) GetProjects(projects []string) ([]ProjectModel, error) {
	var projectJSON SavedProjects
	if _, err := os.Stat(projectsFilePath); os.IsNotExist(err) {
		resp, err := client.api.Get(projectsURLPath, 250)
		if err != nil {
			return nil, err
		}
		err = readJSONFromResp(resp, &projectJSON)
		if err != nil {
			return nil, err
		}
		err = writeJSONToFile(&projectJSON, projectsFilePath)
		if err != nil {
			return nil, err
		}
	} else {
		err := readJSONFromFile(projectsFilePath, &projectJSON)
		if err != nil {
			return nil, err
		}
	}
	return projectJSON.Filter(projects), nil
}

func removeLocalProjectsData() error {
	return os.Remove(projectsFilePath)
}
