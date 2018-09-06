package api

import (
	"fmt"
	"log"
	"os"
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
		for i := 0; i < len(reposJSON); i += batchNumber {
			var r []RepoModel
			if i+batchNumber > len(reposJSON) {
				r = reposJSON[i:len(reposJSON)]
			} else {
				r = reposJSON[i : i+batchNumber]
			}
			go client.getFilesInternal(r, fileChan)
		}
		for range reposJSON {
			keyedFileList := <-fileChan
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
		var filesJSON FileResponse
		resp, err := client.api.Get(filesURLPath(v3.Project.Key, v3.Slug), 1000)
		if err != nil {
			log.Fatal(err)
		}
		err = readJSONFromResp(resp, &filesJSON)
		if err != nil {
			log.Fatal(err)
		}
		repoFiles := SavedFiles{
			v3.Project.Key: {
				v3.Slug: filesJSON,
			},
		}
		c <- repoFiles
	}
}

func removeLocalFilesData() error {
	return os.Remove(filesFilePath)
}
