package router

import (
	"annotate-x/api"

	"github.com/gin-gonic/gin"

	"annotate-x/config"
	"annotate-x/internal/context"

	"annotate-x/cache"
	"annotate-x/db"
	"annotate-x/internal/security"
	"annotate-x/model"
	"annotate-x/repository"
)

func createSuperAdmin() {
	appConfig := config.AppConfig
	db := db.InitDB(appConfig.DATABASE_URL)
	userRepository := repository.NewUserRepository(db)
	if exists, err := userRepository.UsernameExists(appConfig.SUPER_ADMIN_USERNAME); err != nil {
		panic(err.Error())
	} else if exists {
		return
	}
	hashedPassword, err := security.HashPassword(appConfig.SUPER_ADMIN_PASSWORD)
	if err != nil {
		panic(err.Error())
	}
	user := &model.User{
		Username:    appConfig.SUPER_ADMIN_USERNAME,
		Password:    hashedPassword,
		DisplayName: "superadmin",
		Email:       "",
		IsActive:    true,
		Role:        string(model.RoleSuperAdmin),
	}
	if err := userRepository.CreateUser(user); err != nil {
		panic(err.Error())
	}
}

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
	createSuperAdmin()
	r := gin.Default()

	r.Use(context.InjectAppContext(setupAppContext()))

	v1 := r.Group("/v1")

	api.RegisterAuthRouters(v1)
	api.RegisterUsersRouters(v1)

	return r
}
