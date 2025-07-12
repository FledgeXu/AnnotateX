package router

import (
	"annotate-x/api"
	"annotate-x/bootstrap"
	"annotate-x/cache"
	"annotate-x/config"
	"annotate-x/middleware"
	"annotate-x/wire"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	bootstrap.CreateSuperAdmin()

	r := gin.Default()
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	cacheConfig := cache.RedisConfig{
		Addr:     config.GetConfig().REDIS_ADDRESS,
		Password: config.GetConfig().REDIS_PASSWORD,
		DB:       config.GetConfig().REDIS_DB,
	}
	userService := wire.InitIUserService(config.GetConfig().DATABASE_URL)
	authService := wire.InitIAuthService(config.GetConfig().DATABASE_URL, cacheConfig)
	projectService := wire.InitIProjectService(config.GetConfig().DATABASE_URL)
	cacheService := wire.InitICacheService(cacheConfig)

	r.Use(bootstrap.SetupCors())

	v1 := r.Group("/v1")
	v1.Use(middleware.InjectUserHeaderMiddleware())
	api.RegisterAuthRouters(v1, authService, cacheService)
	api.RegisterUserRouters(v1, userService)
	api.RegisterProjectRouters(v1, projectService)

	return r
}
