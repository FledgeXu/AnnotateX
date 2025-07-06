package api

import (
	"annotate-x/model"
	"annotate-x/repository"
	"annotate-x/service"
	"annotate-x/utils"
	"strconv"
	"strings"

	casbin_auth "annotate-x/internal/auth"
	"annotate-x/internal/middleware"
	"annotate-x/internal/security"

	"github.com/gin-gonic/gin"
)

type UsersHandler struct {
	UserRepo    *repository.UserRepository
	UserService *service.UserService
}

func RegisterUsersRouters(rg *gin.RouterGroup,
	userRepo *repository.UserRepository,
	cacheRepo *repository.CacheRepository,
	userService *service.UserService,
) {
	handler := &UsersHandler{userRepo, userService}
	group := rg.Group("/users")
	group.Use(middleware.AuthMiddleware(cacheRepo), middleware.RequirePermissionMiddleware(casbin_auth.Enforcer))

	group.GET("/list", handler.userList)
	group.GET("/me", middleware.AuthMiddleware(cacheRepo), middleware.UserInjectionMiddleware(userRepo), handler.me)
	group.PUT("/me", middleware.AuthMiddleware(cacheRepo), middleware.UserInjectionMiddleware(userRepo), handler.updateMe)
}

func (h *UsersHandler) userList(c *gin.Context) {
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

	users, total, err := h.UserService.GetFilteredUserList(filter)
	if err != nil {
		utils.InternalServerError(c, "Failed to get users")
		return
	}

	utils.OK(c, gin.H{
		"limit":   limit,
		"offset":  offset,
		"total":   total,
		"results": users,
	})
}

func (h *UsersHandler) me(c *gin.Context) {
	user := c.MustGet("currentUser").(*model.User)

	utils.OK(c, model.UserCreateResponse{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Role:        user.Role,
	})
}

func (h *UsersHandler) updateMe(c *gin.Context) {
	user := c.MustGet("currentUser").(*model.User)

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if req.Password != "" {
		hashedPassword, err := security.HashPassword(req.Password)
		if err != nil {
			utils.BadRequest(c, err.Error())
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

	updatedUser, err := h.UserRepo.UpdateUser(user)
	if err != nil {
		utils.InternalServerError(c, err.Error())
	}

	utils.OK(c, model.UserCreateResponse{
		Username:    updatedUser.Username,
		DisplayName: updatedUser.DisplayName,
		Email:       updatedUser.Email,
		Role:        updatedUser.Role,
	})
}
