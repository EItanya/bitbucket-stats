package cache

import (
	"bitbucket/models"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gomodule/redigo/redis"
)

const AllProjectConst = "all_project"
const AllRepositoryConst = "all_repository"
const AllFilesConst = "all_files"

// RedisCache structure of redis cache, implements cache
type RedisCache struct {
	Conn   redis.Conn
	Config *RedisConfig
}

// RedisConfig structure of config for redis cache
type RedisConfig struct {
	Port     string
	Protocol string
}

func (r *RedisCache) write(key string, entity CacheEntity) error {
	var data interface{}
	err := entity.unmarshal(&data)
	if err != nil {
		return err
	}
	switch typedData := data.(type) {
	case []interface{}:
		_, err = r.Conn.Do("SADD", redis.Args{}.Add(key).AddFlat(typedData)...)
	case interface{}:
		_, err = r.Conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(typedData)...)
	default:
		return errors.New("Redis does not support saving of data type passed in")
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisCache) read(keys []string) ([]CacheEntity, error) {
	results := make([]CacheEntity, len(keys))
	if len(keys) == 1 {
		key := keys[0]
		if key == AllProjectConst || key == AllRepositoryConst || key == AllFilesConst {
			keyParam := fmt.Sprintf("%s*", key[4:len(key)])
			values, err := redis.Values(r.Conn.Do("KEYS", redis.Args{}.Add(keyParam)...))
			if err != nil {
				return results, err
			}
			keysPointer := make([]string, len(values))
			err = redis.ScanSlice(values, &keysPointer)
			if err != nil {
				return nil, err
			}
			keys = keysPointer
		}

	}
	for _, key := range keys {
		if splitKey := strings.Split(key, ":"); len(splitKey) > 0 {
			lookupType := splitKey[0]
			re := &RedisEntity{}
			resp, err := r.Conn.Do("HGETALL", redis.Args{}.Add(key)...)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			values, err := redis.Values(resp, nil)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			switch lookupType {
			case "project":
				var dat models.Project
				if err = redis.ScanStruct(values, &dat); err != nil {
					return nil, err
				}
				re.marshal(dat)
			case "repository":
				var dat models.Repository
				if err = redis.ScanStruct(values, &dat); err != nil {
					return nil, err
				}
				re.marshal(dat)
			case "files":
				var dat models.Files
				if err = redis.ScanStruct(values, &dat); err != nil {
					return nil, err
				}
				re.marshal(dat)
			}
			results = append(results, re)

		} else {
			log.Printf("Key passed to redis cache (%s) not of the correct form.\n", key)

		}
	}
	return results, nil
}

func (r *RedisCache) clear() error {
	return r.Conn.Flush()
}

func (r *RedisCache) initialize() error {
	var conn redis.Conn
	var err error
	if r.Config != nil {
		if r.Config.Protocol == "" || r.Config.Port == "" {
			log.Fatalln(errors.New("Redis config is missing information"))
		}
		conn, err = redis.Dial(r.Config.Protocol, fmt.Sprintf(":%s", r.Config.Port))
	} else {
		conn, err = redis.Dial("tcp", ":6379")
	}
	if err != nil {
		log.Println(err)
	}
	r.Conn = conn
	return nil
}
