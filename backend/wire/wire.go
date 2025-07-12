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

// Repos
var projectRepo = wire.NewSet(
	repo.NewProjectRepo,
	wire.Bind(new(repo.IProjectRepo), new(*repo.ProjectRepo)),
)
var userRepo = wire.NewSet(
	repo.NewUserRepo,
	wire.Bind(new(repo.IUserRepo), new(*repo.UserRepo)),
)

var cacheRepo = wire.NewSet(
	repo.NewCacheRepo,
	wire.Bind(new(repo.ICacheRepo), new(*repo.CacheRepo)),
)

// Service
var cacheService = wire.NewSet(
	service.NewCacheService,
	wire.Bind(new(service.ICacheService), new(*service.CacheService)),
)

var authService = wire.NewSet(
	service.NewAuthService,
	wire.Bind(new(service.IAuthService), new(*service.AuthService)),
)

var userService = wire.NewSet(
	service.NewUserService,
	wire.Bind(new(service.IUserService), new(*service.UserService)),
)

var projectService = wire.NewSet(
	service.NewProjectService,
	wire.Bind(new(service.IProjectService), new(*service.ProjectService)),
)

// Repo Providers
var projectRepoProvider = wire.NewSet(
	db.InitDB,
	projectRepo,
)

var userRepoProvider = wire.NewSet(
	db.InitDB,
	userRepo,
)

var cacheRepoProvider = wire.NewSet(
	cache.InitRedis,
	cacheRepo,
)

// Service Providers
var cacheServiceProvider = wire.NewSet(
	cacheRepoProvider,
	cacheService,
)

func InitIAuthService(dsn string, cacheConfig cache.RedisConfig) service.IAuthService {
	wire.Build(
		authService,
		userRepoProvider,
		cacheServiceProvider,
	)
	return nil
}

func InitIUserService(dsn string) service.IUserService {
	wire.Build(
		userService,
		userRepoProvider,
	)
	return nil
}

func InitIProjectService(dsn string) service.IProjectService {
	wire.Build(
		projectService,
		projectRepoProvider,
	)
	return nil
}
