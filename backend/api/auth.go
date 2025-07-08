package api

import (
	"annotate-x/models"
	"annotate-x/service"
	"annotate-x/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	UserService service.IAuthService
}

func RegisterAuthRouters(rg *gin.RouterGroup, userService service.IAuthService) {
	handler := &AuthHandler{userService}
	group := rg.Group("/auth")
	group.POST("/login", handler.login)
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
