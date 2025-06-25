package api

import (
	"time"

	"annotate-x/internal/middleware"
	"annotate-x/model"
	"annotate-x/service"

	"annotate-x/internal/context"

	"annotate-x/internal/security"
	"annotate-x/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterAuthRouters(rg *gin.RouterGroup) {
	group := rg.Group("/auth")
	{
		group.POST("/login", login)
		group.POST("/register", register)
		group.POST("/logout", middleware.AuthMiddleware(), logout)
	}
}

func login(c *gin.Context) {
	appCtx := c.MustGet("appCtx").(*context.AppContext)
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := appCtx.UserRepo.GetUserByUsername(req.Username)
	if err != nil || !user.IsActive {
		utils.Unauthorized(c, "Invalid username or password")
		return
	}

	match, needsRehash, err := security.VerifyPassword(req.Password, user.Password)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}
	if !match {
		utils.Unauthorized(c, "Invalid username or password")
		return
	}

	// Auto-upgrade hash if parameters are outdated
	if needsRehash {
		if newHash, ok, err := security.RehashIfNeeded(req.Password, user.Password); err == nil && ok {
			user.Password = newHash
			_ = appCtx.UserRepo.UpdateUserPassword(user.ID, newHash) // optional error handling
		}
	}

	token, err := security.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.InternalServerError(c, "Failed to generate token")
		return
	}
	utils.OK(c, gin.H{
		"token": token,
	})
}

func register(c *gin.Context) {
	appCtx := c.MustGet("appCtx").(*context.AppContext)
	var req model.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	// for this endpoint role must be unassigned.
	req.Role = string(model.RoleUnassigned)

	user, err := service.NewUserService(appCtx.UserRepo).CreateUser(req)
	if err != nil {
		utils.InternalServerError(c, err.Error())
		return
	}

	utils.Created(c, model.UserCreateResponse{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Role:        user.Role,
	})
}

func logout(c *gin.Context) {
	appCtx := c.MustGet("appCtx").(*context.AppContext)
	tokenRaw, exists := c.Get("rawToken")
	if !exists {
		utils.BadRequest(c, "Missing token")
		return
	}
	tokenStr := tokenRaw.(string)

	claimsRaw, exists := c.Get("jwtClaims")
	if !exists {
		utils.BadRequest(c, "Missing claims")
		return
	}
	claims := claimsRaw.(*security.Claims)

	// Add the token to the Redis blacklist with an expiration time matching the original token.
	expiration := time.Until(claims.ExpiresAt.Time)
	err := appCtx.CacheRepo.BlacklistToken(c.Request.Context(), tokenStr, expiration)
	if err != nil {
		utils.InternalServerError(c, "Failed to logout")
		return
	}
	utils.OK(c, gin.H{})
}
