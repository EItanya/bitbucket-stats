package stats

import (
	"bitbucket-stats/models"
	"sync"
)

// Projects base type for projects/project stats
type Projects struct {
	data []models.Project
	sync.Mutex
}

// Repositories base type for repo/repo stats
type Repositories struct {
	data []models.Repository
	sync.Mutex
}

// Files base type for file/file stats
type Files struct {
	data []string
	sync.Mutex
}
