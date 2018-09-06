package cache

import (
	"fmt"
	"log"
	"strings"
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
		log.Println("Error while attempting to write to cache")
		log.Println(err)
	}
	return err
}

func ReadFromCache(c Cache, keys []string) ([]CacheEntity, error) {
	entities, err := c.read(keys)
	if err != nil {
		log.Println("Error while attempting to read from cache")
		log.Println(err)
	}
	return entities, err
}

func ClearCache(c Cache) error {
	err := c.clear()
	if err != nil {
		log.Println("Error while attempting to clear cache")
		log.Println(err)
	}
	return err
}

func CheckCache(c Cache, keyGroup string) (bool, error) {
	ok, err := c.check(keyGroup)
	if err != nil {
		log.Println("Error while attempting to check cache")
		log.Println(err)
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
	err := ce.marshal(dat)
	if err != nil {
		errSlice := []string{
			"Error while attempting to Marshal Cache entity\n",
			fmt.Sprintf("Error: %s\n", err.Error()),
			fmt.Sprintf("Entity: %+v\n", ce),
			fmt.Sprintf("Data: %+v", dat),
		}
		log.Println(strings.Join(errSlice, ""))
	}
	return err
}

func UnmarshalEntity(ce CacheEntity, dat interface{}) error {
	if ce == nil {
		log.Println("Cache Entity is undefined, skipping translation")
		return nil
	}
	err := ce.unmarshal(dat)
	if err != nil {
		errSlice := []string{
			"Error while attempting to Marshal Cache entity\n",
			fmt.Sprintf("Error: %s\n", err.Error()),
			fmt.Sprintf("Entity: %+v\n", ce),
			fmt.Sprintf("Data: %+v", dat),
		}
		log.Println(strings.Join(errSlice, ""))
	}
	return err
}
