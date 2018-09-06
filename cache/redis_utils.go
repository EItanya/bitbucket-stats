package cache

import (
	"bitbucket-stats/logger"
	"bitbucket-stats/models"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

type redisRepositoryModel struct {
	Slug          string
	ID            int
	Name          string
	ScmID         string
	State         string
	StatusMessage string
	Forkable      bool
	Public        bool
	Project       string
}

func (rp *redisRepositoryModel) initializeRedisRepoModel(r models.Repository) {
	rp.Slug = r.Slug
	rp.ID = r.ID
	rp.Name = r.Name
	rp.ScmID = r.ScmID
	rp.State = r.State
	rp.StatusMessage = r.StatusMessage
	rp.Forkable = r.Forkable
	rp.Public = r.Public
	rp.Project = fmt.Sprintf("project:%s", r.Project.Key)
}

func (rp *redisRepositoryModel) revertToRepositoryModel(p models.Project) models.Repository {
	r := models.Repository{
		Forkable:      rp.Forkable,
		ID:            rp.ID,
		Name:          rp.Name,
		Project:       p,
		Public:        rp.Public,
		ScmID:         rp.ScmID,
		Slug:          rp.Slug,
		State:         rp.State,
		StatusMessage: rp.StatusMessage,
	}
	return r
}

func (r *RedisCache) getValues(redisCommand, key string) ([]interface{}, error) {
	resp, err := r.Conn.Do(redisCommand, redis.Args{}.Add(key)...)
	if err != nil {
		logger.Log.Info(err)
		return nil, err
	}
	values, err := redis.Values(resp, nil)
	if err != nil {
		logger.Log.Info(err)
		return nil, err
	}
	return values, nil
}

func (r *RedisCache) getProject(key string) (*models.Project, error) {
	values, _ := r.getValues("HGETALL", key)
	var dat models.Project
	if err := redis.ScanStruct(values, &dat); err != nil {
		return nil, err
	}
	return &dat, nil
}
