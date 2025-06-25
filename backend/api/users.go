package api

import (
	"annotate-x/model"
	"annotate-x/service"
	"annotate-x/utils"
	"strconv"
	"strings"

	casbin_auth "annotate-x/internal/auth"
	"annotate-x/internal/context"
	"annotate-x/internal/middleware"
	"annotate-x/internal/security"

	"github.com/gin-gonic/gin"

	"net/http"
)

func RegisterUsersRouters(rg *gin.RouterGroup) {
	group := rg.Group("/users")
	group.Use(middleware.AuthMiddleware(), middleware.RequirePermissionMiddleware(casbin_auth.Enforcer))

	group.GET("/list", list)
	group.GET("/me", middleware.AuthMiddleware(), middleware.UserInjectionMiddleware(), me)
	group.PUT("/me", middleware.AuthMiddleware(), middleware.UserInjectionMiddleware(), updateMe)
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
		utils.JSONError(c, http.StatusInternalServerError, "Failed to get users")
		return
	}

	utils.JSONSuccess(c, http.StatusOK, gin.H{
		"limit":   limit,
		"offset":  offset,
		"total":   total,
		"results": users,
	})
}

func me(c *gin.Context) {
	user := c.MustGet("currentUser").(*model.User)

	utils.JSONSuccess(c, http.StatusCreated, model.UserCreateResponse{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Role:        user.Role,
	})
}

func updateMe(c *gin.Context) {
	appCtx := c.MustGet("appCtx").(*context.AppContext)
	user := c.MustGet("currentUser").(*model.User)

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.JSONError(c, http.StatusBadRequest, err.Error())
		return
	}

	if req.Password != "" {
		hashedPassword, err := security.HashPassword(req.Password)
		if err != nil {
			utils.JSONError(c, http.StatusBadRequest, err.Error())
			return
		}
		user.Password = hashedPassword
	}

	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}

	if user.Email != "" {
		user.Email = req.Email
	}

	updatedUser, err := appCtx.UserRepo.UpdateUser(user)
	if err != nil {
		utils.JSONError(c, http.StatusInternalServerError, err.Error())
	}

	utils.JSONSuccess(c, http.StatusCreated, model.UserCreateResponse{
		Username:    updatedUser.Username,
		DisplayName: updatedUser.DisplayName,
		Email:       updatedUser.Email,
		Role:        updatedUser.Role,
	})
}
