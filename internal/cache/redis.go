package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	rdb *redis.Client
}

func NewRedisClient(addr string) *RedisClient {
	return &RedisClient{
		rdb: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (r *RedisClient) PushTask(ctx context.Context, fileID string) error {
	return r.rdb.LPush(ctx, "file_tasks", fileID).Err()
}

func (r *RedisClient) PopTask(ctx context.Context) (string, error) {
	results, err := r.rdb.BLPop(ctx, 0, "file_tasks").Result()
	if err != nil {
		return "", err
	}

	if len(results) > 1 {
		return results[1], nil
	}

	return "", nil
}

func (r *RedisClient) SetResult(ctx context.Context, fileID string, data string) error {
	return r.rdb.Set(ctx, "result:"+fileID, data, 24*time.Hour).Err()
}

func (r *RedisClient) GetResult(ctx context.Context, fileID string) (string, error) {
	if r == nil || r.rdb == nil {
		return "", fmt.Errorf("redis not initialized")
	}
	val, err := r.rdb.Get(ctx, "result:"+fileID).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}
