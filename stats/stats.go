package stats

import (
	"bitbucket/api"
	"sync"
)

// Projects base type for projects/project stats
type Projects struct {
	data api.SavedProjects
	sync.Mutex
}

// Repositories base type for repo/repo stats
type Repositories struct {
	data api.SavedRepos
	sync.Mutex
}

// Files base type for file/file stats
type Files struct {
	data api.SavedRepos
	sync.Mutex
}
