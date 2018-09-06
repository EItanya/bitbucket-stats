package cache

import (
	"bitbucket/logger"
	"bitbucket/models"
	"errors"
	"fmt"
	"strings"

	"github.com/gomodule/redigo/redis"
)

var defaultRedisCacheConfig = &RedisCacheConfig{
	Port:     "6379",
	Protocol: "tcp",
}

// RedisCache structure of redis cache, implements cache
type RedisCache struct {
	Conn   redis.Conn
	Config *RedisCacheConfig
}

// RedisCacheConfig structure of config for redis cache
type RedisCacheConfig struct {
	Port     string
	Protocol string
}

func (r *RedisCache) write(key string, entity CacheEntity) error {
	var data interface{}
	err := entity.Unmarshal(&data)
	if err != nil {
		return err
	}
	switch typedData := data.(type) {
	case models.Files:
		if len(typedData) == 0 {
			return nil
		}
		_, err = r.Conn.Do("SADD", redis.Args{}.Add(key).AddFlat(typedData)...)
	case models.Project:
		_, err = r.Conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(typedData)...)
	case models.Repository:
		rrm := redisRepositoryModel{}
		rrm.initializeRedisRepoModel(typedData)
		_, err = r.Conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(rrm)...)
	default:
		return errors.New("Redis does not support saving of data type passed in")
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisCache) read(keys []string) ([]CacheEntity, error) {
	if len(keys) == 1 {
		key := keys[0]
		if key == AllProjectConst || key == AllRepositoryConst || key == AllFilesConst {
			keysPointer, err := r.getKeys(key)
			if err != nil {
				return nil, err
			}
			keys = keysPointer
		}
	}
	results := make([]CacheEntity, 0)
	for _, key := range keys {
		if splitKey := strings.Split(key, ":"); len(splitKey) > 0 {
			lookupType := splitKey[0]
			switch lookupType {
			case "project", "repository":
				values, err := r.getValues("HGETALL", key)
				if lookupType == "project" {
					var dat models.Project
					if err = redis.ScanStruct(values, &dat); err != nil {
						return nil, err
					}
					results = append(results, &dat)
				} else {
					var dat redisRepositoryModel
					if err = redis.ScanStruct(values, &dat); err != nil {
						return nil, err
					}
					pm, _ := r.getProject(dat.Project)
					rm := dat.revertToRepositoryModel(*pm)
					results = append(results, &rm)
				}
			case "files":
				resp, err := r.Conn.Do("SMEMBERS", redis.Args{}.Add(key)...)
				if err != nil {
					logger.Log.Info(err)
					return nil, err
				}
				values, err := redis.Values(resp, nil)
				if err != nil {
					logger.Log.Info(err)
					return nil, err
				}
				var dat models.Files
				if err = redis.ScanSlice(values, &dat); err != nil {
					return nil, err
				}
				extFiles := models.FilesID{
					Files:      dat,
					ProjectKey: splitKey[1],
					RepoSlug:   splitKey[2],
				}
				results = append(results, &extFiles)
			}

		} else {
			logger.Log.Infof("Key passed to redis cache (%s) not of the correct form.\n", key)

		}
	}
	return results, nil
}

func (r *RedisCache) check(keyGroup string) (bool, error) {
	if keyGroup == AllFilesConst || keyGroup == AllProjectConst || keyGroup == AllRepositoryConst {
		keys, err := r.getKeys(keyGroup)
		if err != nil {
			return false, err
		}
		if len(keys) > 0 {
			return true, nil
		}
	}
	return false, nil
}

func (r *RedisCache) clear() error {
	_, err := r.Conn.Do("FLUSHDB")
	return err
}

func (r *RedisCache) initialize() error {
	var conn redis.Conn
	var err error
	if r.Config != nil {
		if r.Config.Protocol == "" || r.Config.Port == "" {
			logger.Log.Fatal(errors.New("Redis config is missing information"))
		}
		conn, err = redis.Dial(r.Config.Protocol, fmt.Sprintf(":%s", r.Config.Port))
	} else {
		conn, err = redis.Dial("tcp", ":6379")
	}
	if err != nil {
		logger.Log.Info(err)
	}
	r.Conn = conn
	return nil
}

func (r *RedisCache) getKeys(key string) ([]string, error) {
	keys := make([]string, 0)
	keyParam := fmt.Sprintf("%s*", key[4:len(key)])
	values, err := redis.Values(r.Conn.Do("KEYS", redis.Args{}.Add(keyParam)...))
	if err != nil {
		return nil, err
	}

	err = redis.ScanSlice(values, &keys)
	if err != nil {
		return nil, err
	}
	return keys, err
}
