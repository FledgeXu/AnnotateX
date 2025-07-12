package api

import (
	"annotate-x/models"
	"annotate-x/service"
	"annotate-x/utils"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UserService  service.IAuthService
	CacheService service.ICacheService
}

func RegisterAuthRouters(rg *gin.RouterGroup, userService service.IAuthService, cacheService service.ICacheService) {
	handler := &AuthHandler{userService, cacheService}
	group := rg.Group("/auth")
	group.POST("/login", handler.login)
	group.POST("/logout", handler.logout)
}

func (h *AuthHandler) login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	user, token, err := h.UserService.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil || !user.IsActive {
		utils.Unauthorized(c, err.Error())
		return
	}
	userResp := models.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	utils.OK(c, gin.H{
		"token": token,
		"user":  userResp,
	})
}

func (h *AuthHandler) logout(c *gin.Context) {
	fmt.Println(c.GetHeader(models.XUserID))
	userId, err := strconv.ParseInt(c.GetHeader(models.XUserID), 10, 64)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	expiration, err := utils.UnixStringToTime(c.GetHeader(models.XExpiresAt))
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	err = h.UserService.Logout(c.Request.Context(), userId, expiration)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.OK(c, gin.H{})
}
