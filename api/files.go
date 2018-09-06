package api

import (
	"bitbucket/cache"
	"bitbucket/logger"
	"bitbucket/models"
	"errors"
	"fmt"

	"github.com/gosuri/uiprogress"
)

var filesURLPath = func(projKey, repoSlug string) string {
	return fmt.Sprintf("/projects/%s/repos/%s/files", projKey, repoSlug)
}

// GetFiles get all repos from Bitbucket
func (client *Client) GetFiles(repos map[string][]string) (*[]models.FilesID, error) {

	if present, err := cache.CheckCache(client.cache, cache.AllFilesConst); !present && err == nil {
		var fileChan = make(chan *models.FilesID)
		reposJSONPtr, err := client.GetRepos(make([]string, 0))
		reposJSON := *reposJSONPtr
		if err != nil {
			return nil, err
		}

		fileListResults := make([]models.FilesID, len(reposJSON))
		numRepos := len(reposJSON)
		bar := uiprogress.AddBar(numRepos + 1)
		bar.AppendCompleted()
		bar.PrependFunc(prependFormatFunc(func(b *uiprogress.Bar) string {

			if b.Current() == b.Total {
				return "File data retrieved"
			} else if b.Current() >= numRepos {
				return "Saving file data"
			}
			return "Dowloading file data"
		}))
		for i := 0; i < numRepos; i += batchNumber {
			var r []models.Repository
			if i+batchNumber > numRepos {
				r = reposJSON[i:numRepos]
			} else {
				r = reposJSON[i : i+batchNumber]
			}
			go client.getFilesInternal(r, fileChan)
		}
		for range reposJSON {
			fileList := <-fileChan
			bar.Incr()
			if fileList == nil || len(fileList.Files) == 0 {
				continue
			}
			key := fmt.Sprintf("files:%s:%s", fileList.ProjectKey, fileList.RepoSlug)
			err = cache.SetCacheValue(client.cache, key, &fileList.Files)
			if err != nil {
				return nil, err
			}
			fileListResults = append(fileListResults, *fileList)
		}
		bar.Incr()
		return &fileListResults, nil
	} else if present && err == nil {
		entities, err := cache.GetCacheValue(client.cache, []string{cache.AllFilesConst})
		if err != nil {
			return nil, err
		}
		translatedEntities := make([]models.FilesID, 0)
		for _, ce := range entities {
			var dat models.FilesID
			if ce == nil {
				continue
			}
			err = cache.UnmarshalEntity(ce, &dat)
			if err != nil {
				return nil, err
			}
			translatedEntities = append(translatedEntities, dat)
		}
		// results := models.FilterFiles(&translatedEntities, repos)
		return &translatedEntities, nil
	}
	logger.Log.Info(errors.New("Reached end of GetFiles function with no data, check logic"))
	return nil, nil
}

func (client *Client) getFilesInternal(r []models.Repository, c chan *models.FilesID) {
	for _, v3 := range r {
		var repoFiles models.FilesID
		// Make sure both keys are set
		if v3.Project.Key != "" && v3.Slug != "" {
			var collector FileResponse
			nextPageStart := 0
			for {
				var filesJSON FileResponse
				opts := urlOptions{
					limit: 1000,
				}
				if nextPageStart > 0 {
					opts.start = nextPageStart
				}
				resp, err := client.api.Get(filesURLPath(v3.Project.Key, v3.Slug), opts)
				if err != nil {
					logger.Log.Fatal(err)
				}
				err = filesJSON.UnmarshalHTTP(resp)
				if err != nil {
					logger.Log.Fatal(err)
				}
				collector.Values = append(collector.Values, filesJSON.Values...)
				if filesJSON.IsLastPage {
					break
				}
				nextPageStart = filesJSON.NextPageStart
			}

			repoFiles = models.FilesID{
				ProjectKey: v3.Project.Key,
				RepoSlug:   v3.Slug,
				Files:      collector.Values,
			}
		}

		c <- &repoFiles
	}
}
