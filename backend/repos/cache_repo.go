package repos

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type ICacheRepository interface {
	Set(key string, value string, ttlSeconds int) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)
}

type CacheRepository struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewCacheRepository(client *redis.Client) *CacheRepository {
	return &CacheRepository{
		Client: client,
		Ctx:    context.Background(),
	}
}

func (r *CacheRepository) Set(key string, value string, ttlSeconds int) error {
	return r.Client.Set(r.Ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *CacheRepository) Get(key string) (string, error) {
	return r.Client.Get(r.Ctx, key).Result()
}

func (r *CacheRepository) Delete(key string) error {
	return r.Client.Del(r.Ctx, key).Err()
}

func (r *CacheRepository) Exists(key string) (bool, error) {
	n, err := r.Client.Exists(r.Ctx, key).Result()
	return n > 0, err
}
