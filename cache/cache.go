package cache

import (
	"bitbucket/logger"

	"go.uber.org/zap"
)

// Cache interface for Caches
type Cache interface {
	write(key string, entity CacheEntity) error
	read(keys []string) ([]CacheEntity, error)
	check(keyGroup string) (bool, error)
	clear() error
	initialize() error
}

func SaveToCache(c Cache, key string, entity CacheEntity) error {
	err := c.write(key, entity)
	if err != nil {
		logger.Log.Info("Error while attempting to write to cache")
		logger.Log.Info(err)
	}
	return err
}

func ReadFromCache(c Cache, keys []string) ([]CacheEntity, error) {
	entities, err := c.read(keys)
	if err != nil {
		logger.Log.Info("Error while attempting to read from cache")
		logger.Log.Info(err)
	}
	return entities, err
}

func ClearCache(c Cache) error {
	err := c.clear()
	if err != nil {
		logger.Log.Info("Error while attempting to clear cache")
		logger.Log.Info(err)
	}
	return err
}

func CheckCache(c Cache, keyGroup string) (bool, error) {
	ok, err := c.check(keyGroup)
	if err != nil {
		logger.Log.Info("Error while attempting to check cache")
		logger.Log.Info(err)
	}
	return ok, err
}

func NewRedisCache(config RedisConfig) (*RedisCache, error) {
	cache := &RedisCache{
		Config: &config,
	}
	err := cache.initialize()
	return cache, err
}

type CacheEntity interface {
	marshal(dat interface{}) error
	unmarshal(dat interface{}) error
}

func MarshalEntity(ce CacheEntity, dat interface{}) error {
	if ce == nil {
		logger.Log.Info("Cache Entity is undefined, skipping marshal to cache entity")
		return nil
	}
	err := ce.marshal(dat)
	if err != nil {
		logger.Log.Infow("Error while attempting to Marshal Cache entity",
			zap.Error(err),
			zap.Reflect("Entity", ce),
			zap.Reflect("Data", dat),
		)
	}
	return err
}

func UnmarshalEntity(ce CacheEntity, dat interface{}) error {
	if ce == nil {
		logger.Log.Info("Cache Entity is undefined, skipping Unmarshal from cache entity")
		return nil
	}
	err := ce.unmarshal(dat)
	if err != nil {
		logger.Log.Infow("Error while attempting to Unmarshal Cache entity",
			zap.Error(err),
			zap.Reflect("Entity", ce),
			zap.Reflect("Data", dat),
		)
	}
	return err
}
