package repo

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type ICacheRepo interface {
	Set(ctx context.Context, key string, value string, ttlSeconds int) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

type CacheRepo struct {
	Client *redis.Client
}

func NewCacheRepo(client *redis.Client) *CacheRepo {
	return &CacheRepo{
		Client: client,
	}
}

func (r *CacheRepo) Set(ctx context.Context, key string, value string, ttlSeconds int) error {
	return r.Client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *CacheRepo) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *CacheRepo) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

func (r *CacheRepo) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.Client.Exists(ctx, key).Result()
	return n > 0, err
}
