package api

import (
	"annotate-x/models"
	"annotate-x/service"
	"annotate-x/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService service.IAuthService
}

func RegisterAuthRouters(rg *gin.RouterGroup, authService service.IAuthService) {
	handler := &AuthHandler{authService}
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
	user, token, err := h.AuthService.Login(c.Request.Context(), req.Username, req.Password)
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
	userId, err := strconv.ParseInt(c.GetHeader(models.XUserID), 10, 64)
	if err != nil {
		utils.Unauthorized(c, "Invalid Bearer Token")
		return
	}

	expiration, err := utils.UnixStringToTime(c.GetHeader(models.XExpiresAt))
	if err != nil {
		utils.Unauthorized(c, "Invalid Bearer Token")
		return
	}

	err = h.AuthService.Logout(c.Request.Context(), userId, expiration)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.OK(c, gin.H{})
}
