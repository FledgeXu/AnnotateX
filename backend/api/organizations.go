package api

import (
	casbin_auth "annotate-x/internal/auth"
	"annotate-x/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterOrganizationsRouters(rg *gin.RouterGroup) {
	group := rg.Group("/organizations")
	group.Use(middleware.AuthMiddleware(), middleware.RequirePermissionMiddleware(casbin_auth.Enforcer))
}
