package cache

import (
	"bitbucket/logger"

	"go.uber.org/zap"
)

const defaultDir = "data"

const projectConst = "project"
const repositoryConst = "repository"
const filesConst = "files"

const AllProjectConst = "all_" + projectConst
const AllRepositoryConst = "all_" + repositoryConst
const AllFilesConst = "all_" + filesConst

// Cache interface for Caches
type Cache interface {
	set(key string, entity CacheEntity) error
	get(keys []string) ([]CacheEntity, error)
	check(keyGroup string) (bool, error)
	clear() error
	initialize() error
}

func SetCacheValue(c Cache, key string, entity CacheEntity) error {
	err := c.set(key, entity)
	if err != nil {
		logger.Log.Errorw(
			"Error while attempting to write to cache",
			zap.Error(err),
		)
	}
	return err
}

func GetCacheValue(c Cache, keys []string) ([]CacheEntity, error) {
	entities, err := c.get(keys)
	if err != nil {
		logger.Log.Errorw(
			"Error while attempting to read from cache",
			zap.Error(err),
		)
	}
	return entities, err
}

func ClearCache(c Cache) error {
	err := c.clear()
	if err != nil {
		logger.Log.Errorw(
			"Error while attempting to clear cache",
			zap.Error(err),
		)
	}
	return err
}

func CheckCache(c Cache, keyGroup string) (bool, error) {
	ok, err := c.check(keyGroup)
	if err != nil {
		logger.Log.Errorw(
			"Error while attempting to check cache",
			zap.Error(err),
		)
	}
	return ok, err
}

func NewRedisCache(config *RedisCacheConfig) (*RedisCache, error) {
	if config == nil {
		config = defaultRedisCacheConfig
	}
	cache := &RedisCache{
		Config: config,
	}
	err := cache.initialize()
	return cache, err
}

func NewFileCache(config *FileCacheConfig) (*FileCache, error) {
	if config == nil {
		config = defaultFileCacheConfig
	}
	cache := &FileCache{
		Config: config,
	}
	err := cache.initialize()
	return cache, err
}

type CacheEntity interface {
	Marshal(dat interface{}) error
	Unmarshal(dat interface{}) error
}

func MarshalEntity(ce CacheEntity, dat interface{}) error {
	if ce == nil {
		logger.Log.Info("Cache Entity is undefined, skipping marshal to cache entity")
		return nil
	}
	err := ce.Marshal(dat)
	if err != nil {
		logger.Log.Errorw("Error while attempting to Marshal Cache entity",
			zap.Error(err),
			zap.Reflect("Entity", ce),
			zap.Reflect("Data", dat),
		)
	}
	return err
}

func UnmarshalEntity(ce CacheEntity, dat interface{}) error {
	if ce == nil {
		logger.Log.Infof("Cache Entity (%+v) is undefined, skipping Unmarshal from cache entity", ce)
		return nil
	}
	err := ce.Unmarshal(dat)
	if err != nil {
		logger.Log.Errorw("Error while attempting to Unmarshal Cache entity",
			zap.Error(err),
			zap.Reflect("Entity", ce),
			zap.Reflect("Data", dat),
		)
	}
	return err
}
