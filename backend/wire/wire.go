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

var authRepoProvider = wire.NewSet(
	db.InitDB,
	userRepo,
)

var cacheRepoProvider = wire.NewSet(
	cache.InitRedis,
	cacheRepo,
)

var cacheServiceProvider = wire.NewSet(
	cacheRepoProvider,
	cacheService,
)

func InitICacheService(cacheConfig cache.RedisConfig) service.ICacheService {
	wire.Build(
		cacheService,
		cacheRepoProvider,
	)
	return nil
}

func InitIAuthService(dsn string, cacheConfig cache.RedisConfig) service.IAuthService {
	wire.Build(
		authService,
		authRepoProvider,
		cacheServiceProvider,
	)
	return nil
}
