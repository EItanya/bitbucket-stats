package api

import (
	"fmt"
	"log"
	"os"

	"github.com/gosuri/uiprogress"
)

// SavedFiles is the format by which files are saved
type SavedFiles map[string]map[string]FileResponse

var filesURLPath = func(projKey, repoSlug string) string {
	return fmt.Sprintf("/projects/%s/repos/%s/files", projKey, repoSlug)
}

const filesFilePath = "data/files.json"

// GetFiles get all repos from Bitbucket
func (client *Client) GetFiles(repos map[string][]string) (*SavedFiles, error) {
	var allFilesJSON = make(SavedFiles)

	if _, err := os.Stat(filesFilePath); os.IsNotExist(err) {
		var fileChan = make(chan SavedFiles)
		reposJSON, err := getReposJSON()
		if err != nil {
			return nil, err
		}
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
			var r []RepoModel
			if i+batchNumber > numRepos {
				r = reposJSON[i:numRepos]
			} else {
				r = reposJSON[i : i+batchNumber]
			}
			go client.getFilesInternal(r, fileChan)
		}
		for range reposJSON {
			keyedFileList := <-fileChan
			bar.Incr()
			for projectKey, value := range keyedFileList {
				project, repoExists := allFilesJSON[projectKey]
				if repoExists {
					for repoSlug, value1 := range value {
						project[repoSlug] = value1
						allFilesJSON[projectKey] = project
					}
				} else {
					allFilesJSON[projectKey] = value
				}
			}
		}
		err = writeJSONToFile(&allFilesJSON, filesFilePath)
		if err != nil {
			return nil, err
		}
		bar.Incr()
	} else {
		err := readJSONFromFile(filesFilePath, &allFilesJSON)
		if err != nil {
			return nil, err
		}
	}

	return &allFilesJSON, nil
}

func (client *Client) getFilesInternal(r []RepoModel, c chan SavedFiles) {
	for _, v3 := range r {
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
				log.Fatal(err)
			}
			err = readJSONFromResp(resp, &filesJSON)
			if err != nil {
				log.Fatal(err)
			}
			collector.Values = append(collector.Values, filesJSON.Values...)
			if filesJSON.IsLastPage {
				break
			}
			nextPageStart = filesJSON.NextPageStart
		}

		repoFiles := SavedFiles{
			v3.Project.Key: {
				v3.Slug: collector,
			},
		}
		c <- repoFiles
	}
}

func removeLocalFilesData() error {
	return os.Remove(filesFilePath)
}
