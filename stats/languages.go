package stats

import (
	"bitbucket/api"
	"bitbucket/models"
	"fmt"
	"log"
	"sync"
)

// Context is the object which holds the counter data for the Context
type Context struct {
	RawFileData       languageMap
	FileDataByRepo    []repoLanguageData
	FileDataByProject []projectLanguageData
	files             *[]models.FilesID
	repos             *[]models.Repository
	projects          *[]models.Project
	TotalFileCount    int
}

//Initialize initialize stats Context Object
func (c *Context) Initialize(client *api.Client) error {
	files, err := client.GetFiles(make(map[string][]string))
	if err != nil || files == nil {
		log.Fatalf("Error while trying to retrieve files\n%s\nFiles: %+v", err.Error(), files)
	}
	c.files = files
	repos, err := client.GetRepos(make([]string, 0))
	if err != nil || repos == nil {
		log.Fatalf("Error while trying to retrieve repos\n%s\nRepos: %+v", err.Error(), repos)
	}
	c.repos = repos
	projects, err := client.GetProjects(make([]string, 0))
	if err != nil || projects == nil {
		log.Fatalf("Error while trying to retrieve projects\n%s\nProjects: %+v", err.Error(), projects)
	}
	c.projects = projects
	c.TotalFileCount = 0
	return nil
}

// ToJSON Based on key passed in
func (c *Context) ToJSON(key string) (string, bool) {
	var result string
	var err error
	switch key {
	case "RawFileData":
		result, err = c.RawFileData.ToJSON()
	}
	if err != nil {
		return "", false
	}
	return result, true
}

// CountAllFiles counts all
func (c *Context) CountAllFiles() {
	counter := &languageData{
		Data: make(languageMap),
	}
	for _, files := range *c.files {
		counter.Add(1)
		go counter.countFiles(files.Files)
	}
	counter.Wait()
	for _, value := range counter.Data {
		c.TotalFileCount += value
	}
	c.RawFileData = counter.Data
}

// CountFilesByRepo by repo
func (c *Context) CountFilesByRepo() {
	counter := &organizedLanguageData{
		data: make([]repoLanguageData, 0),
	}
	filesChan := make(chan *models.FilesID)
	for _, repo := range *c.repos {
		go getFilesByRepoSlug(c.files, repo.Slug, filesChan)
	}

	for range *c.repos {
		files := <-filesChan
		if files != nil {
			counter.Add(1)
			go counter.countFiles(files.Files, files.RepoSlug, files.ProjectKey)
		}
	}
	counter.Wait()
	c.FileDataByRepo = counter.data
}

// CountFilesByProject by project
func (c *Context) CountFilesByProject() {
	counter := &organizedLanguageData{
		data: make([]repoLanguageData, 0),
	}
	filesChan := make(chan *[]models.FilesID)
	for _, project := range *c.projects {
		go func(files *[]models.FilesID, projectKey string, ch chan *[]models.FilesID) {
			var result []models.FilesID
			for _, val := range *files {
				if val.ProjectKey == projectKey {
					result = append(result, val)
				}
			}
			ch <- &result
		}(c.files, project.Key, filesChan)
	}

	for range *c.projects {
		files := <-filesChan
		for _, v := range *files {
			counter.Add(1)
			go counter.countFiles(v.Files, v.RepoSlug, v.ProjectKey)
		}
	}
	counter.Wait()
	c.FileDataByProject = combineReposIntoProjects(c.projects, counter.data)
}

// CountFilesByLanguage by project
// func (c *Context) CountFilesByLanguage(langs []string) {
// 	c.RawFileData = make(languageMap)
// }

// ReposWithNodeModules by project
func (c *Context) ReposWithNodeModules() []string {
	filesChan := make(chan *models.FilesID)
	for _, repo := range *c.repos {
		go getFilesByRepoSlug(c.files, repo.Slug, filesChan)
	}

	wg := sync.WaitGroup{}
	lock := sync.Mutex{}
	result := make([]string, 0)
	for range *c.repos {
		fileList := <-filesChan
		wg.Add(1)
		go func(list []string, projectKey, repoSlug string) {
			defer wg.Done()
			if findItem("node_modules", list) {
				lock.Lock()
				result = append(result, fmt.Sprintf("%s:%s", projectKey, repoSlug))
				lock.Unlock()
			}
		}(fileList.Files, fileList.ProjectKey, fileList.RepoSlug)
	}
	wg.Wait()
	return result
}

// GetDataByLanguage gets organized data by language
func (c *Context) GetDataByLanguage(langs []string) []string {
	lock := sync.Mutex{}
	wg := sync.WaitGroup{}
	result := make([]string, 0)
	if c.FileDataByRepo == nil {
		c.CountFilesByRepo()
	}

	wg.Add(len(c.FileDataByRepo))
	for _, val := range c.FileDataByRepo {
		go func(r repoLanguageData, langs []string) {
			defer wg.Done()
			for _, inputLang := range langs {
				for lang := range r.Stats.Data {
					if inputLang == lang {
						lock.Lock()
						result = append(result, fmt.Sprintf("%s:%s", r.ProjectKey, r.RepoSlug))
						lock.Unlock()
					}
				}
			}
		}(val, langs)
	}
	wg.Wait()

	return result
}
