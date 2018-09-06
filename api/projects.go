package api

import (
	"bitbucket/cache"
	"bitbucket/models"
	"fmt"
	"os"
)

// SavedProjects is the format by which projects are saved
type SavedProjects ProjectResponse

// Filter method to filter saved projects
func (data *SavedProjects) Filter(projects []string) []models.Project {
	if len(projects) == 0 {
		return data.Values
	}
	filteredProjects := make([]models.Project, 0)
	ch := make(chan []models.Project)
	for _, val := range projects {
		go data.filterProjects(val, data.Values, ch)
	}
	for range projects {
		filteredProjects = append(filteredProjects, <-ch...)
	}
	return filteredProjects
}

func (data SavedProjects) filterProjects(val string, p []models.Project, ch chan []models.Project) {
	pm := make([]models.Project, 0)
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
func (client *Client) GetProjects(projects []string) (*[]models.Project, error) {

	redisCache, err := cache.NewRedisCache(cache.RedisConfig{
		Port:     "6379",
		Protocol: "tcp",
	})
	entities, err := cache.ReadFromCache(redisCache, []string{"all_project"})
	if err != nil {
		return nil, err
	}
	translatedEntities := make([]models.Project, 0)
	for _, ce := range entities {
		var dat models.Project
		err = cache.UnmarshalEntity(ce, &dat)
		if err != nil {
			return nil, err
		}
		translatedEntities = append(translatedEntities, dat)
	}
	result := &translatedEntities
	return result, nil
}

func removeLocalProjectsData() error {
	return os.Remove(projectsFilePath)
}
