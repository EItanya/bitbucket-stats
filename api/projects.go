package api

import (
	"fmt"
	"os"

	"github.com/gosuri/uiprogress"
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

// GetProjects get all projects from Bitbucket
func (client *Client) GetProjects(projects []string) (*[]ProjectModel, error) {
	var projectJSON SavedProjects
	// projectJSON := &ProjectResponse{}
	if _, err := os.Stat(projectsFilePath); os.IsNotExist(err) {
		bar := uiprogress.AddBar(2)
		bar.AppendCompleted()
		bar.PrependFunc(prependFormatFunc(func(b *uiprogress.Bar) string {
			if b.Current() == b.Total {
				return "Project data retrieved"
			} else if b.Current() >= 1 {
				return "Saving project data"
			}
			return "Dowloading project data"
		}))
		opts := urlOptions{
			limit: 250,
		}
		resp, err := client.api.Get(projectsURLPath, opts)

		projectJSON := &ProjectResponse{}
		// GetEntity(resp, projectJSON)

		if err != nil {
			return nil, err
		}
		bar.Incr()
		err = readJSONFromResp(resp, &projectJSON)
		if err != nil {
			return nil, err
		}
		err = writeJSONToFile(&projectJSON, projectsFilePath)
		if err != nil {
			return nil, err
		}
		bar.Incr()
	} else {
		err := readJSONFromFile(projectsFilePath, &projectJSON)
		if err != nil {
			return nil, err
		}
	}
	result := projectJSON.Filter(projects)
	return &result, nil
}

func removeLocalProjectsData() error {
	return os.Remove(projectsFilePath)
}
