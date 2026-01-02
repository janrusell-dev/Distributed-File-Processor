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
