package stats

import (
	"bitbucket/api"
	"fmt"
	"sync"
)

// Context is the object which holds the counter data for the Context
type Context struct {
	RawFileData       languageMap
	FileDataByRepo    []repoLanguageData
	FileDataByProject []projectLanguageData
	files             *api.SavedFiles
	repos             *api.SavedRepos
	projects          *[]api.ProjectModel
	TotalFileCount    int
}

//Initialize initialize stats Context Object
func (c *Context) Initialize(client *api.Client) error {
	files, err := client.GetFiles(make(map[string][]string))
	if err != nil {
		return err
	}
	c.files = &files
	repos, err := client.GetRepos(make(map[string][]string))
	if err != nil {
		return err
	}
	c.repos = &repos
	projects, err := client.GetProjects(make([]string, 0))
	if err != nil {
		return err
	}
	c.projects = &projects
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
		data: make(languageMap),
	}
	for _, repo := range *c.files {
		for _, fileJSON := range repo {
			counter.Add(1)
			go counter.countFiles(fileJSON.Values)
		}
	}
	counter.Wait()
	for _, value := range counter.data {
		c.TotalFileCount += value
	}
	c.RawFileData = counter.data
}

// CountFilesByRepo by repo
func (c *Context) CountFilesByRepo() {
	counter := &organizedLanguageData{
		data: make([]repoLanguageData, 0),
	}
	for _, repo := range *c.repos {
		if projectRepoList, ok := (*c.files)[repo.Project.Key]; ok {
			if repoFileResponses, ok := projectRepoList[repo.Slug]; ok {
				counter.Add(1)
				go counter.countFiles(repoFileResponses.Values, repo.Slug, repo.Project.Key)
			}
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
	for _, project := range *c.projects {
		if projectRepoList, ok := (*c.files)[project.Key]; ok {
			for repoSlug, fileList := range projectRepoList {
				counter.Add(1)
				go counter.countFiles(fileList.Values, repoSlug, project.Key)
			}
		}
	}
	counter.Wait()
	c.FileDataByProject = combineReposIntoProjects(c.projects, counter.data)
}

// CountFilesByLanguage by project
func (c *Context) CountFilesByLanguage(langs []string) {
	c.RawFileData = make(languageMap)
}

// ReposWithNodeModules by project
func (c *Context) ReposWithNodeModules() []string {
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}
	result := make([]string, 0)
	for projectKey, repos := range *c.files {
		for repoSlug, repoFiles := range repos {
			wg.Add(1)
			go func(list []string, projectKey, repoSlug string) {
				defer wg.Done()
				if findItem("node_modules", list) {
					lock.Lock()
					result = append(result, fmt.Sprintf("%s:%s", projectKey, repoSlug))
					lock.Unlock()
				}
			}(repoFiles.Values, projectKey, repoSlug)
		}
	}
	wg.Wait()
	return result
}
