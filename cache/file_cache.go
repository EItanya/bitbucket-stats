package cache

import (
	"errors"
	"os"
	"strings"
)

// Type of FileCache
type FileCache struct {
	Config *FileCacheConfig
	Files  map[string]*fileCacheTable
}

// Type of FileCacheConfig
type FileCacheConfig struct {
	Dir string
}

func (c *FileCache) write(key string, entity CacheEntity) error {
	var data interface{}
	err := entity.unmarshal(&data)
	if err != nil {
		return err
	}
	splitKey := strings.SplitN(key, ":", 2)
	if len(splitKey) != 2 {
		return errors.New("Key is not in the correct format to write to cache")
	}
	cacheKey := fileCacheKey{
		Key:      splitKey[0],
		Location: splitKey[1],
	}
	switch splitKey[0] {
	case projectConst:
		err := fileCacheSet(projectsCacheTable, cacheKey, data)
		if err != nil {
			return err
		}
	case repositoryConst:
	case filesConst:
	default:
		return errors.New("Cache key did not include valid prefix")
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *FileCache) read(keys []string) ([]CacheEntity, error) {
	result := make([]CacheEntity, len(keys))
	return result, nil
}

func (c *FileCache) check(keyGroup string) (bool, error) {
	return false, nil
}

func (c *FileCache) clear() error {
	return nil
}

func (c *FileCache) initialize() error {
	if c.Config == nil {

	}
	// Create data directory if none exists
	if _, err1 := os.Stat(c.Config.Dir); os.IsNotExist(err1) {
		if err2 := os.Mkdir(c.Config.Dir, os.ModeDir); err2 != nil {
			return err2
		}
	}

	return nil
}
