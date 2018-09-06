package cache

import (
	"errors"
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

// Redis structure of redis cache, implements cache
type Redis struct {
	Conn   redis.Conn
	Config *RedisConfig
}

// RedisConfig structure of config for redis cache
type RedisConfig struct {
	Port     string
	Protocol string
}

func (r *Redis) write(key string) error {
	return nil
}

func (r *Redis) read(key string) (string, error) {
	return "", nil
}

func (r *Redis) clear() error {
	return r.Conn.Flush()
}

func (r *Redis) initialize() error {
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
