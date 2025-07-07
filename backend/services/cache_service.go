package services

import (
	"annotate-x/repos"
	"context"
)

type ICacheService interface {
	BlacklistToken(token string, ttlSeconds int) error
	IsTokenBlacklisted(token string) (bool, error)
}

type CacheService struct {
	Repo repos.ICacheRepository
	Ctx  context.Context
}

func NewCacheService(repo repos.ICacheRepository) *CacheService {
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
