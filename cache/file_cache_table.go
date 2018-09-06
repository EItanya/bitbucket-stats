package cache

import (
	"bitbucket/logger"
	"sync"
	"time"

	"go.uber.org/zap"
)

type fileCacheTable interface {
	get(key fileCacheKey) (CacheEntity, error)
	set(key fileCacheKey, dat CacheEntity) error
	keys() []string
	read() error
	write() error
	clear() error
}

func fileCacheGet(t fileCacheTable, key fileCacheKey) (CacheEntity, error) {
	dat, err := t.get(key)
	if err != nil {
		logger.Log.Errorw("Error while attempting to get item from file cache",
			zap.Error(err),
			zap.Reflect("Key", key),
			zap.Reflect("Data Pointer", dat),
		)
	}
	return dat, err
}

func fileCacheSet(t fileCacheTable, key fileCacheKey, dat CacheEntity) error {
	err := t.set(key, dat)
	if err != nil {
		logger.Log.Errorw("Error while attempting to set item in file cache map",
			zap.Error(err),
			zap.Reflect("Key", key),
			zap.Reflect("Data Pointer", dat),
		)
	}
	return err
}

func fileCacheRead(t fileCacheTable) error {
	err := t.read()
	if err != nil {
		logger.Log.Errorw("Error while attempting to read cache table from file",
			zap.Error(err),
		)
	}
	return err
}

func fileCacheWrite(t fileCacheTable) error {
	err := t.write()
	if err != nil {
		logger.Log.Errorw("Error while attempting to write cache to file",
			zap.Error(err),
		)
	}
	return err
}

func fileCacheClear(t fileCacheTable) error {
	err := t.clear()
	if err != nil {
		logger.Log.Errorw("Error while attempting to clear cache table",
			zap.Error(err),
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
	lastSync time.Time
}
