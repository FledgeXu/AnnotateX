package service

import (
	"annotate-x/repo"
	"context"
)

type ICacheService interface {
	BlacklistToken(token string, ttlSeconds int) error
	IsTokenBlacklisted(token string) (bool, error)
}

type CacheService struct {
	Repo repo.ICacheRepository
	Ctx  context.Context
}

func NewCacheService(repo repo.ICacheRepository) *CacheService {
	return &CacheService{
		Repo: repo,
		Ctx:  context.Background(),
	}
}

func (s *CacheService) BlacklistToken(token string, ttlSeconds int) error {
	key := "blacklist:" + token
	return s.Repo.Set(key, "1", ttlSeconds)
}

func (s *CacheService) IsTokenBlacklisted(token string) (bool, error) {
	key := "blacklist:" + token
	return s.Repo.Exists(key)
}
