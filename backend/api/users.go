package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterUsersRouters(rg *gin.RouterGroup) {
	auth := rg.Group("/users")
	{
		auth.GET("/list", list)
	}
}

func list(c *gin.Context) {
}
