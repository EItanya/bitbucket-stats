package stats

import (
	"strings"
	"sync"
)

const sectionSize = 500

type Categories struct {
	Possibilities []string
	MostLikely    string
	lock          *sync.Mutex
	wg            *sync.WaitGroup
}

func (c *Categories) categorizeRepo(fileSlice []string) {
	defer c.wg.Done()
	for _, val := range fileSlice {
		if strings.Contains(val, "package.json") {
			c.Possibilities = append(c.Possibilities, jsProject)
		} else if strings.Contains(val, "pom.xml") {

		}
	}
}

// CategorizeRepo attempt to categorize a repo
func (c *Categories) CategorizeRepo(files []string, lmap languageMap) {
	numberOfFiles := len(files)
	for i := 0; i < numberOfFiles; i += sectionSize {
		var fileSlice []string
		if i+sectionSize >= numberOfFiles {
			fileSlice = files[i:numberOfFiles]
		} else {
			fileSlice = files[i : i+sectionSize]
		}
		c.wg.Add(1)
		go c.categorizeRepo(fileSlice)
	}
	c.wg.Wait()
}
