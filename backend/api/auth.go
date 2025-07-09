package api

import (
	"annotate-x/models"
	"annotate-x/service"
	"annotate-x/utils"
	"annotate-x/utils/security"
	"time"

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
	user, token, err := h.UserService.Login(req.Username, req.Password)
	if err != nil || !user.IsActive {
		utils.Unauthorized(c, err.Error())
		return
	}
	userResp := models.UserCreateResponse{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
	}

	utils.OK(c, gin.H{
		"token": token,
		"user":  userResp,
	})
}

func (h *AuthHandler) logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	tokenStr := authHeader[7:]
	claims, err := security.ParseToken(tokenStr)
	if err != nil {
		utils.BadRequest(c, "Invalid or expired token")
		return
	}

	expiration := time.Until(claims.ExpiresAt.Time)
	err = h.CacheService.BlacklistToken(tokenStr, int(expiration))
	if err != nil {
		utils.InternalServerError(c, "Failed to logout")
		return
	}
	utils.OK(c, gin.H{})
}
