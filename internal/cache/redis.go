package cache

import (
	"context"

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
