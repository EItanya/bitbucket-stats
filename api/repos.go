package api

import (
	"bitbucket-stats/cache"
	"bitbucket-stats/logger"
	"bitbucket-stats/models"
	"errors"
	"fmt"
	"os"

	"github.com/gosuri/uiprogress"
)

var reposURLPath = func(projKey string) string {
	return fmt.Sprintf("/projects/%s/repos", projKey)
}

const reposFilePath = "data/repos.json"

// GetRepos get all repos from Bitbucket
func (client *Client) GetRepos(repos []string) (*[]models.Repository, error) {
	var reposJSON []models.Repository
	var repoChan = make(chan []models.Repository)

	if present, err := cache.CheckCache(client.cache, cache.AllRepositoryConst); !present && err == nil {
		projectJSON, err := client.GetProjects(make([]string, 0))
		if err != nil {
			return nil, err
		}
		numProjects := len(*projectJSON)
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
		for _, v := range *projectJSON {
			go client.getReposInternal(v, repoChan)
		}
		for range *projectJSON {
			reposJSON = append(reposJSON, <-repoChan...)
			bar.Incr()
		}

		for _, v := range reposJSON {
			key := fmt.Sprintf("repository:%s", v.Slug)
			err = cache.SetCacheValue(client.cache, key, &v)
			if err != nil {
				return nil, err
			}
		}
		bar.Incr()
		result := models.FilterRepos(&reposJSON, repos)
		return &result, nil
	} else if present && err == nil {
		entities, err := cache.GetCacheValue(client.cache, []string{cache.AllRepositoryConst})
		if err != nil {
			return nil, err
		}
		translatedEntities := make([]models.Repository, len(entities))
		for index, ce := range entities {
			var dat models.Repository
			err = cache.UnmarshalEntity(ce, &dat)
			if err != nil {
				return nil, err
			}
			if dat.Slug != "" && dat.Project.Key != "" {
				translatedEntities[index] = dat
			}
		}
		results := models.FilterRepos(&translatedEntities, repos)
		return &results, nil
	} else if err != nil {
		logger.Log.Info(err)
		return nil, err
	}
	return nil, errors.New("Reached end of GetRepos function with no data, check logic")
}

func (client *Client) getReposInternal(v models.Project, c chan []models.Repository) {
	var repoJSON RepoResponse
	opts := urlOptions{
		limit: 500,
	}
	resp, _ := client.api.Get(reposURLPath(v.Key), opts)
	_ = repoJSON.UnmarshalHTTP(resp)
	c <- repoJSON.Values
}

func removeLocalReposData() error {
	return os.Remove(reposFilePath)
}
