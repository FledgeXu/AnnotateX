package api

import (
	"annotate-x/model"
	"annotate-x/service"
	"strconv"
	"strings"

	casbin_auth "annotate-x/internal/auth"
	"annotate-x/internal/context"
	"annotate-x/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterUsersRouters(rg *gin.RouterGroup) {
	auth := rg.Group("/users")
	auth.Use(middleware.AuthMiddleware(), middleware.RequirePermissionMiddleware(casbin_auth.Enforcer))

	auth.GET("/list", list)
}

func list(c *gin.Context) {
	appCtx := c.MustGet("appCtx").(*context.AppContext)
	// Parse pagination
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Parse sorting
	sortBy := c.DefaultQuery("sort_by", "created_at")
	order := strings.ToLower(c.DefaultQuery("order", "desc"))
	if order != "asc" && order != "desc" {
		order = "desc"
	}

	// Optional search and filter
	keyword := c.Query("keyword")
	isActive := c.Query("is_active") // "true", "false", or ""

	filter := model.UserFilter{
		Keyword:  keyword,
		IsActive: isActive,
		SortBy:   sortBy,
		Order:    order,
		Limit:    limit,
		Offset:   offset,
	}

	users, total, err := service.NewUserService(appCtx.UserRepo).GetFilteredUserList(filter)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to get users"})
		return
	}

	c.JSON(200, gin.H{
		"limit":   limit,
		"offset":  offset,
		"total":   total,
		"results": users,
	})
}
