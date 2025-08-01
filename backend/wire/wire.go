//go:build wireinject
// +build wireinject

package wire

import (
	"annotate-x/cache"
	"annotate-x/config"
	"annotate-x/db"
	"annotate-x/models"
	"annotate-x/mq"
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

var datasetRepo = wire.NewSet(
	repo.NewDatasetRepo,
	wire.Bind(new(repo.IDatasetRepo), new(*repo.DatasetRepo)),
)

var cacheRepo = wire.NewSet(
	repo.NewCacheRepo,
	wire.Bind(new(repo.ICacheRepo), new(*repo.CacheRepo)),
)

var s3Repo = wire.NewSet(
	repo.NewS3Repo,
	wire.Bind(new(repo.IS3Repo), new(*repo.S3Repo)),
)

var mqRepo = wire.NewSet(
	repo.NewMqRepo,
	wire.Bind(new(repo.IMqRepo), new(*repo.MqRepo)),
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

var datasetService = wire.NewSet(
	service.NewDatasetService,
	wire.Bind(new(service.IDatasetService), new(*service.DatasetService)),
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

var datasetRepoProvider = wire.NewSet(
	db.InitDB,
	datasetRepo,
)

var cacheRepoProvider = wire.NewSet(
	cache.InitRedis,
	cacheRepo,
)

var mqProvider = wire.NewSet(
	mq.InitMQ,
	mqRepo,
)

// Service Providers
var cacheServiceProvider = wire.NewSet(
	cacheRepoProvider,
	cacheService,
)

func InitIAuthService(dsn models.DataSourceName, cacheConfig cache.RedisConfig) service.IAuthService {
	wire.Build(
		authService,
		userRepoProvider,
		cacheServiceProvider,
	)
	return nil
}

func InitIUserService(dsn models.DataSourceName) service.IUserService {
	wire.Build(
		userService,
		userRepoProvider,
	)
	return nil
}

func InitIProjectService(dsn models.DataSourceName) service.IProjectService {
	wire.Build(
		projectService,
		projectRepoProvider,
	)
	return nil
}

func InitIDatasetService(dsn models.DataSourceName, s3Config config.S3Config, bucketName models.BucketName, mqUrl models.MQUrl) service.IDatasetService {
	wire.Build(
		datasetService,
		datasetRepoProvider,
		s3Repo,
		mqProvider,
	)
	return nil
}
