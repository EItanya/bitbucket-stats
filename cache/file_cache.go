package cache

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Type of FileCache
type FileCache struct {
	Config *FileCacheConfig
	Files  map[string]fileCacheTable
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

	cacheKey, err := c.translateCacheKey(key)
	if err != nil {
		return err
	}

	cacheTable, err := c.getCacheTable(cacheKey)
	if err != nil {
		return err
	}

	err = fileCacheSet(cacheTable, cacheKey, data)
	if err != nil {
		return err
	}
	// switch cacheKey.Location {
	// case projectConst:
	// 	err := fileCacheSet(projectsCacheTable, cacheKey, data)
	// 	if err != nil {
	// 		return err
	// 	}
	// case repositoryConst:
	// 	err := fileCacheSet(repositoriesCacheTable, cacheKey, data)
	// 	if err != nil {
	// 		return err
	// 	}
	// case filesConst:
	// 	err := fileCacheSet(filesCacheTable, cacheKey, data)
	// 	if err != nil {
	// 		return err
	// 	}
	// default:
	// 	return errors.New("Cache key did not include valid prefix")
	// }
	// if err != nil {
	// 	return err
	// }
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

	c.Files = make(map[string]fileCacheTable)
	c.Files[projectConst] = projectsCacheTable
	c.Files[repositoryConst] = repositoriesCacheTable
	c.Files[filesConst] = filesCacheTable

	return nil
}

func (c *FileCache) getCacheTable(key fileCacheKey) (fileCacheTable, error) {
	if val, ok := c.Files[key.Location]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("Unable to find cache table for given key location: (%s)", key.Location)
}

func (c FileCache) translateCacheKey(key string) (fileCacheKey, error) {
	var cacheKey fileCacheKey
	splitKey := strings.SplitN(key, ":", 2)
	if len(splitKey) != 2 {
		return cacheKey, errors.New("Key is not in the correct format to write to cache")
	}
	cacheKey = fileCacheKey{
		Key:      splitKey[1],
		Location: splitKey[0],
	}
	return cacheKey, nil
}
