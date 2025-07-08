//go:build wireinject
// +build wireinject

package wire

import (
	"annotate-x/cache"
	"annotate-x/db"
	"annotate-x/repo"
	"annotate-x/service"

	"github.com/google/wire"
)

var userRepo = wire.NewSet(
	repo.NewUserRepo,
	wire.Bind(new(repo.IUserRepo), new(*repo.UserRepo)),
)

var cacheRepo = wire.NewSet(
	repo.NewCacheRepo,
	wire.Bind(new(repo.ICacheRepo), new(*repo.CacheRepo)),
)

var cacheService = wire.NewSet(
	service.NewCacheService,
	wire.Bind(new(service.ICacheService), new(*service.CacheService)),
)

var authService = wire.NewSet(
	service.NewAuthService,
	wire.Bind(new(service.IAuthService), new(*service.AuthService)),
)

func InitICacheService(cacheConfig cache.RedisConfig) service.ICacheService {
	wire.Build(
		cache.InitRedis,
		cacheRepo,
		cacheService,
	)
	return nil
}

func InitIAuthService(dsn string, cacheConfig cache.RedisConfig) service.IAuthService {
	wire.Build(
		db.InitDB,
		cache.InitRedis,
		userRepo,
		cacheRepo,
		cacheService,
		authService,
	)
	return nil
}
