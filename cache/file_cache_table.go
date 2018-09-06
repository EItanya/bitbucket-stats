package cache

import (
	"bitbucket/logger"
	"sync"

	"go.uber.org/zap"
)

type fileCacheTable interface {
	get(key fileCacheKey) (interface{}, error)
	set(key fileCacheKey, dat interface{}) error
}

func fileCacheGet(t fileCacheTable, key fileCacheKey) (interface{}, error) {
	dat, err := t.get(key)
	if err != nil {
		logger.Log.Infow("Error while attempting to get item from file cache",
			zap.Error(err),
			zap.Reflect("Key", key),
			zap.Reflect("Data Pointer", dat),
		)
	}
	return dat, err
}

func fileCacheSet(t fileCacheTable, key fileCacheKey, dat interface{}) error {
	err := t.set(key, dat)
	if err != nil {
		logger.Log.Infow("Error while attempting to set item in file cache map",
			zap.Error(err),
			zap.Reflect("Key", key),
			zap.Reflect("Data Pointer", dat),
		)
	}
	return err
}

// Struct used to hold file cache keys for map lookup
type fileCacheKey struct {
	Location, Key string
}

type fileCacheTableBasic struct {
	dataMtx  *sync.RWMutex
	fileMtx  *sync.RWMutex
	title    string
	filename string
}
