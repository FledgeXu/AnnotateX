package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository struct {
	rdb *redis.Client
}

func NewCacheRepository(rdb *redis.Client) *CacheRepository {
	return &CacheRepository{rdb: rdb}
}

func (r *CacheRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	exists, err := r.rdb.Exists(ctx, "blacklist:"+token).Result()
	return exists == 1, err
}

func (r *CacheRepository) BlacklistToken(ctx context.Context, token string, ttl time.Duration) error {
	return r.rdb.Set(ctx, "blacklist:"+token, "1", ttl).Err()
}
