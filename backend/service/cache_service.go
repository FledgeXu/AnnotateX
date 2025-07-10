package service

import (
	"annotate-x/repo"
	"context"
)

type ICacheService interface {
	BlacklistToken(ctx context.Context, token string, ttlSeconds int) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}

type CacheService struct {
	Repo repo.ICacheRepo
	Ctx  context.Context
}

func NewCacheService(repo repo.ICacheRepo) *CacheService {
	return &CacheService{
		Repo: repo,
		Ctx:  context.Background(),
	}
}

func (s *CacheService) BlacklistToken(ctx context.Context, token string, ttlSeconds int) error {
	key := "blacklist:" + token
	return s.Repo.Set(ctx, key, "1", ttlSeconds)
}

func (s *CacheService) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := "blacklist:" + token
	return s.Repo.Exists(ctx, key)
}
