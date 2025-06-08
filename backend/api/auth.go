package api

import (
	"net/http"
	"time"

	"annotate-x/internal/middleware"
	"annotate-x/model"
	"annotate-x/service"

	"annotate-x/internal/context"

	"annotate-x/internal/security"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func RegisterAuthRouters(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", login)
		auth.POST("/register", register)
		auth.GET("/me", middleware.AuthMiddleware(), middleware.UserInjectionMiddleware(), me)
		auth.POST("/logout", middleware.AuthMiddleware(), logout)
	}
}

func login(c *gin.Context) {
	appCtx := c.MustGet("appCtx").(*context.AppContext)
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := appCtx.UserRepo.GetUserByUsername(req.Username)
	if err != nil || !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	match, needsRehash, err := security.VerifyPassword(req.Password, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !match {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func register(c *gin.Context) {
	appCtx := c.MustGet("appCtx").(*context.AppContext)
	var req model.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// for this endpoint role must be unassigned.
	req.Role = string(model.RoleUnassigned)

	user, err := service.NewUserService(appCtx.UserRepo).CreateUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model.UserCreateResponse{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		Role:        user.Role,
	})
}

func me(c *gin.Context) {
	user := c.MustGet("currentUser").(*model.User)

	c.JSON(http.StatusCreated, model.UserCreateResponse{
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing token"})
		return
	}
	tokenStr := tokenRaw.(string)

	claimsRaw, exists := c.Get("jwtClaims")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing claims"})
		return
	}
	claims := claimsRaw.(*security.Claims)

	// Add the token to the Redis blacklist with an expiration time matching the original token.
	expiration := time.Until(claims.ExpiresAt.Time)
	err := appCtx.CacheRepo.BlacklistToken(c.Request.Context(), tokenStr, expiration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
