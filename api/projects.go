package api

import (
	"bitbucket/cache"
	"bitbucket/models"
	"errors"
	"fmt"
	"os"

	"github.com/gosuri/uiprogress"
)

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
	if present, err := cache.CheckCache(client.cache, cache.AllProjectConst); !present && err == nil {
		bar := uiprogress.AddBar(2)
		bar.AppendCompleted()
		bar.PrependFunc(prependFormatFunc(func(b *uiprogress.Bar) string {
			if b.Current() == b.Total {
				return "Project data retrieved"
			} else if b.Current() >= 1 {
				return "Writing to cache"
			}
			return "Dowloading project data"
		}))
		opts := urlOptions{
			limit: 250,
		}
		resp, err := client.api.Get(projectsURLPath, opts)

		if err != nil {
			return nil, err
		}
		bar.Incr()
		var jsonResp ProjectResponse
		err = jsonResp.UnmarshalHTTP(resp)
		if err != nil {
			return nil, err
		}

		for _, v := range jsonResp.Values {
			key := fmt.Sprintf("project:%s", v.Key)
			err = cache.SetCacheValue(client.cache, key, &v)
			if err != nil {
				return nil, err
			}
		}
		bar.Incr()
		results := models.FilterProjects(&jsonResp.Values, projects)
		return &results, nil
	} else if present && err == nil {
		entities, err := cache.GetCacheValue(client.cache, []string{cache.AllProjectConst})
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
		results := models.FilterProjects(&translatedEntities, projects)
		return &results, nil
	}
	return nil, errors.New("Reached end of GetProjects function with no data, check logic")
}

func removeLocalProjectsData() error {
	return os.Remove(projectsFilePath)
}
