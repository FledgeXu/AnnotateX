package service_test

import (
	"annotate-x/mocks" // 引入 mockery 生成的 mock
	"annotate-x/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheService_BlacklistToken(t *testing.T) {
	mockRepo := mocks.NewMockICacheRepository(t)
	cacheService := service.NewCacheService(mockRepo)

	token := "abc123"
	key := "blacklist:" + token
	ttl := 3600

	mockRepo.EXPECT().Set(key, "1", ttl).Return(nil)

	err := cacheService.BlacklistToken(token, ttl)

	assert.NoError(t, err)
}

func TestCacheService_IsTokenBlacklisted(t *testing.T) {
	mockRepo := mocks.NewMockICacheRepository(t)
	cacheService := service.NewCacheService(mockRepo)

	token := "expired-token"
	key := "blacklist:" + token

	mockRepo.EXPECT().Exists(key).Return(true, nil)

	result, err := cacheService.IsTokenBlacklisted(token)

	assert.NoError(t, err)
	assert.True(t, result)
}
