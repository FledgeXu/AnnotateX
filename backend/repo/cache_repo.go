package repo

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type ICacheRepo interface {
	Set(key string, value string, ttlSeconds int) error
	Get(key string) (string, error)
	Delete(key string) error
	Exists(key string) (bool, error)
}

type CacheRepo struct {
	Client *redis.Client
	Ctx    context.Context
}

func NewCacheRepo(client *redis.Client) *CacheRepo {
	return &CacheRepo{
		Client: client,
		Ctx:    context.Background(),
	}
}

func (r *CacheRepo) Set(key string, value string, ttlSeconds int) error {
	return r.Client.Set(r.Ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *CacheRepo) Get(key string) (string, error) {
	return r.Client.Get(r.Ctx, key).Result()
}

func (r *CacheRepo) Delete(key string) error {
	return r.Client.Del(r.Ctx, key).Err()
}

func (r *CacheRepo) Exists(key string) (bool, error) {
	n, err := r.Client.Exists(r.Ctx, key).Result()
	return n > 0, err
}
