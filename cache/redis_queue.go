package cache

import (
	"bitbucket-stats/logger"
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type RedisAction struct {
	command, key string
}

type RedisValueKey struct {
	Values []interface{}
	Key    string
}

func (r *RedisCache) redisActionQueue(ctx context.Context, actionChan chan RedisAction) chan *RedisValueKey {
	valueChan := make(chan *RedisValueKey)
	go func(ctx context.Context, actionChan chan RedisAction, valueChan chan *RedisValueKey) {
		wg := &sync.WaitGroup{}
		defer func() {
			wg.Wait()
			close(valueChan)
		}()

		for {
			select {
			case val := <-actionChan:
				wg.Add(1)
				fmt.Println(val)
				go r.queueRedisAction(val, valueChan, wg)
			case <-ctx.Done():
				return
			}
		}
	}(ctx, actionChan, valueChan)
	return valueChan
}

func (r *RedisCache) queueRedisAction(action RedisAction, valueChan chan *RedisValueKey, wg *sync.WaitGroup) {
	var result *RedisValueKey
	values, err := r.getValues(action.command, action.key)
	if err != nil {
		logger.Log.Warnw("Redis action failed to complete",
			zap.Reflect("Action:", action),
			zap.Error(err),
		)
	} else {
		result = new(RedisValueKey)
		result.Values = values
		result.Key = action.key
	}
	wg.Done()
	valueChan <- result
}
