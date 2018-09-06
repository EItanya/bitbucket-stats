package stats

import (
	"encoding/json"
	"sync"
)

type languageMap map[string]int

func (l *languageMap) ToJSON() (string, error) {
	byt, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return "", err
	}
	return string(byt), nil
}

type languageData struct {
	Data  languageMap
	Total int
	sync.WaitGroup
	sync.Mutex
}

func (counter *languageData) countFiles(files []string) {
	defer counter.Done()
	counter.Total += len(files)
	for _, val := range files {
		fileType, validFileType := getFiletype(val)
		counter.Lock()
		if languageKey, ok := extensionMap[fileType]; ok && validFileType {
			if val, ok := counter.Data[languageKey]; ok {
				counter.Data[languageKey] = val + 1
			} else {
				counter.Data[languageKey] = 1
			}
		} else {
			if val, ok := counter.Data[other]; ok {
				counter.Data[other] = val + 1
			} else {
				counter.Data[other] = 1
			}
		}
		counter.Unlock()
	}
}

type projectLanguageData struct {
	Stats      *languageData
	ProjectKey string
}

func (l *projectLanguageData) ToJSON() (string, error) {
	byt, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return "", err
	}
	return string(byt), nil
}

type repoLanguageData struct {
	Stats      *languageData
	ProjectKey string
	RepoSlug   string
}

func (l *repoLanguageData) ToJSON() (string, error) {
	byt, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		return "", err
	}
	return string(byt), nil
}

type organizedLanguageData struct {
	data []repoLanguageData
	sync.WaitGroup
	sync.Mutex
}

func (counter *organizedLanguageData) countFiles(files []string, repoSlug, projectKey string) {
	defer counter.Done()
	internalCounter := languageData{
		Data: make(languageMap),
	}
	internalCounter.Add(1)
	internalCounter.countFiles(files)
	internalCounter.Wait()
	counter.Lock()
	extendedRepo := repoLanguageData{
		Stats:      &internalCounter,
		RepoSlug:   repoSlug,
		ProjectKey: projectKey,
	}
	counter.data = append(counter.data, extendedRepo)
	counter.Unlock()
}

type dataByLanguage struct {
	Repos    []string
	Language string
	sync.WaitGroup
	sync.Mutex
}
