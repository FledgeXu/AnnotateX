package router

import (
	"annotate-x/api"
	"annotate-x/bootstrap"
	"annotate-x/cache"
	"annotate-x/config"
	"annotate-x/wire"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	bootstrap.CreateSuperAdmin()

	r := gin.Default()

	cacheConfig := cache.RedisConfig{
		Addr:     config.GetConfig().REDIS_ADDRESS,
		Password: config.GetConfig().REDIS_PASSWORD,
		DB:       config.GetConfig().REDIS_DB,
	}
	authService := wire.InitIAuthService(config.GetConfig().DATABASE_URL, cacheConfig)
	cacheService := wire.InitICacheService(cacheConfig)

	r.Use(bootstrap.SetupCors())

	v1 := r.Group("/v1")

	api.RegisterAuthRouters(v1, authService, cacheService)

	return r
}
