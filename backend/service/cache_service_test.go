package service_test

import (
	"annotate-x/mocks"
	"annotate-x/service"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheService_BlacklistToken(t *testing.T) {
	mockRepo := mocks.NewMockICacheRepo(t)
	context := context.Background()
	cacheService := service.NewCacheService(mockRepo)

	token := "abc123"
	key := "blacklist:" + token
	ttl := 3600

	mockRepo.EXPECT().Set(context, key, "1", ttl).Return(nil)

	err := cacheService.BlacklistToken(context, token, ttl)

	assert.NoError(t, err)
}

func TestCacheService_IsTokenBlacklisted(t *testing.T) {
	mockRepo := mocks.NewMockICacheRepo(t)
	context := context.Background()
	cacheService := service.NewCacheService(mockRepo)

	token := "expired-token"
	key := "blacklist:" + token

	mockRepo.EXPECT().Exists(context, key).Return(true, nil)

	result, err := cacheService.IsTokenBlacklisted(context, token)

	assert.NoError(t, err)
	assert.True(t, result)
}
