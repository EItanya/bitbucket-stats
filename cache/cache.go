package cache

import "log"

// Cache interface for Caches
type Cache interface {
	write(key string, entity CacheEntity) error
	read(keys []string) ([]CacheEntity, error)
	clear() error
	initialize() error
}

func SaveToCache(c Cache, key string, entity CacheEntity) error {
	return c.write(key, entity)
}

func ReadFromCache(c Cache, keys []string) ([]CacheEntity, error) {

	return c.read(keys)
}

func ClearCache(c Cache) error {
	return c.clear()
}

func NewRedisCache(config RedisConfig) (*RedisCache, error) {
	cache := &RedisCache{
		Config: &config,
	}
	err := cache.initialize()
	return cache, err
}

// func NewFileCache(config FileConfig) (FileCache, error) {

// }

type CacheEntity interface {
	marshal(dat interface{}) error
	unmarshal(dat interface{}) error
}

func MarshalEntity(ce CacheEntity, dat interface{}) error {
	return ce.marshal(dat)
}

func UnmarshalEntity(ce CacheEntity, dat interface{}) error {
	if ce == nil {
		log.Println("Cache Entity is undefined, skipping translation")
		return nil
	}
	return ce.unmarshal(dat)
}
