package router

import (
	"annotate-x/bootstrap"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	bootstrap.CreateSuperAdmin()

	r := gin.Default()

	r.Use(bootstrap.SetupCors())

	// v1 := r.Group("/v1")

	// api.RegisterAuthRouters(v1, injectObject.UserRepo, injectObject.CacheRepo, injectObject.UserService)

	return r
}
