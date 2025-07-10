package repo_test

import (
	"annotate-x/repo"
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupTestRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   9,
	})
}

func clearTestRedisData(client *redis.Client) {
	client.FlushDB(context.Background())
}

func TestCacheRepository(t *testing.T) {
	client := setupTestRedisClient()
	defer client.Close()
	clearTestRedisData(client)

	cacheRepo := repo.NewCacheRepo(client)
	context := context.Background()

	t.Run("Set and Get", func(t *testing.T) {
		err := cacheRepo.Set(context, "test_key", "test_value", 10)
		assert.NoError(t, err)

		val, err := cacheRepo.Get(context, "test_key")
		assert.NoError(t, err)
		assert.Equal(t, "test_value", val)
	})

	t.Run("Get non-existing key", func(t *testing.T) {
		_, err := cacheRepo.Get(context, "non_existent_key")
		assert.Error(t, err)
		assert.Equal(t, redis.Nil, err)
	})

	t.Run("Delete", func(t *testing.T) {
		err := cacheRepo.Set(context, "delete_key", "to_delete", 10)
		assert.NoError(t, err)

		err = cacheRepo.Delete(context, "delete_key")
		assert.NoError(t, err)

		_, err = cacheRepo.Get(context, "delete_key")
		assert.Error(t, err)
		assert.Equal(t, redis.Nil, err)
	})

	t.Run("Exists", func(t *testing.T) {
		cacheRepo.Set(context, "exist_key", "yes", 10)

		exists, err := cacheRepo.Exists(context, "exist_key")
		assert.NoError(t, err)
		assert.True(t, exists)

		exists, err = cacheRepo.Exists(context, "not_exist_key")
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("TTL expiry", func(t *testing.T) {
		err := cacheRepo.Set(context, "short_ttl", "temp", 1)
		assert.NoError(t, err)

		time.Sleep(2 * time.Second)

		_, err = cacheRepo.Get(context, "short_ttl")
		assert.Error(t, err)
		assert.Equal(t, redis.Nil, err)
	})
}
