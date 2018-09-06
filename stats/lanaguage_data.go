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
	data  languageMap
	total int
	sync.WaitGroup
	sync.Mutex
}

func (counter *languageData) countFiles(files []string) {
	defer counter.Done()
	counter.total += len(files)
	for _, val := range files {
		fileType, validFileType := getFiletype(val)
		counter.Lock()
		if languageKey, ok := extensionMap[fileType]; ok && validFileType {
			if val, ok := counter.data[languageKey]; ok {
				counter.data[languageKey] = val + 1
			} else {
				counter.data[languageKey] = 1
			}
		} else {
			if val, ok := counter.data[other]; ok {
				counter.data[other] = val + 1
			} else {
				counter.data[other] = 1
			}
		}
		counter.Unlock()
	}
}

type projectLanguageData struct {
	Stats      languageMap
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
	Stats      languageMap
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
		data: make(languageMap),
	}
	internalCounter.Add(1)
	internalCounter.countFiles(files)
	internalCounter.Wait()
	counter.Lock()
	extendedRepo := repoLanguageData{
		Stats:      internalCounter.data,
		RepoSlug:   repoSlug,
		ProjectKey: projectKey,
	}
	counter.data = append(counter.data, extendedRepo)
	counter.Unlock()
}
