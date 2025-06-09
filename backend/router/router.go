package router

import (
	"annotate-x/api"

	"github.com/gin-gonic/gin"

	"annotate-x/config"
	"annotate-x/internal/context"

	"annotate-x/cache"
	"annotate-x/db"
	"annotate-x/repository"
)

func setupAppContext() *context.AppContext {
	appConfig := config.AppConfig

	db := db.InitDB(appConfig.DATABASE_URL)
	redis := cache.InitRedis(appConfig.REDIS_ADDRESS, appConfig.REDIS_PASSWORD, appConfig.REDIS_DB)
	userRepository := repository.NewUserRepository(db)
	cacheRepository := repository.NewCacheRepository(redis)

	appContext := context.AppContext{
		UserRepo:  userRepository,
		CacheRepo: cacheRepository,
	}

	return &appContext
}

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(context.InjectAppContext(setupAppContext()))

	v1 := r.Group("/v1")

	api.RegisterAuthRouters(v1)
	api.RegisterUsersRouters(v1)

	return r
}
