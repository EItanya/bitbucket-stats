package stats

import (
	"bitbucket-stats/models"
	"strings"
)

func combineLanguageMaps(langDataMapsExt []repoLanguageData) *languageData {
	langDataMaps := make([]*languageData, 0)
	for _, v := range langDataMapsExt {
		langDataMaps = append(langDataMaps, v.Stats)
	}
	result := languageData{
		Data:  make(languageMap),
		Total: 0,
	}
	for _, langData := range langDataMaps {
		for lang, counter := range langData.Data {
			result.Total += counter
			if _, ok := result.Data[lang]; ok {
				result.Data[lang] += counter
			} else {
				result.Data[lang] = counter
			}
		}
	}
	return &result
}

func combineProject(val splitRepos, c chan projectLanguageData) {
	data := combineLanguageMaps(val.data)
	c <- projectLanguageData{
		ProjectKey: val.projectKey,
		Stats:      data,
	}
}

type splitRepos struct {
	data       []repoLanguageData
	projectKey string
}

func splitReposIntoProjects(proj models.Project, langData []repoLanguageData, langChan chan splitRepos) {
	jointRepos := make([]repoLanguageData, 0)
	for _, lang := range langData {
		if lang.ProjectKey == proj.Key {
			jointRepos = append(jointRepos, lang)
		}
	}
	result := splitRepos{
		data:       jointRepos,
		projectKey: proj.Key,
	}
	langChan <- result
}

func combineReposIntoProjects(projects *[]models.Project, langData []repoLanguageData) []projectLanguageData {
	langChan := make(chan splitRepos)
	splitByProject := make([]splitRepos, 0)
	result := make([]projectLanguageData, 0)
	for _, proj := range *projects {
		go splitReposIntoProjects(proj, langData, langChan)

	}
	for range *projects {
		splitByProject = append(splitByProject, <-langChan)
	}
	projectDataChan := make(chan projectLanguageData)
	for _, val := range splitByProject {
		go combineProject(val, projectDataChan)
	}
	for range splitByProject {
		result = append(result, <-projectDataChan)
	}
	return result
}

func getFiletype(filepath string) (string, bool) {
	if splitVal := strings.Split(filepath, "."); len(splitVal) > 1 {
		return splitVal[len(splitVal)-1], true
	}

	if splitBySlash := strings.Split(filepath, "/"); len(splitBySlash) > 0 {
		return splitBySlash[len(splitBySlash)-1], true
	}

	return "", false
}

func findItem(item string, list []string) bool {
	for _, val := range list {
		if strings.Contains(val, item) {
			return true
		}
	}
	return false
}

func getFilesByRepoSlug(files *[]models.FilesID, repoSlug string, ch chan *models.FilesID) {
	var result *models.FilesID
	if repoSlug != "" {
		for _, val := range *files {
			if val.RepoSlug == repoSlug {
				result = &val
				break
			}
		}
	}
	ch <- result
}
