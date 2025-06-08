package api

import (
	"net/http"
	"time"

	"annotate-x/internal/middleware"

	"annotate-x/internal/context"
	"annotate-x/repository"

	"annotate-x/internal/security"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required,min=8,max=64"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email" binding:"omitempty,email"`
}

func RegisterAuthRouters(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", postLogin)
		auth.POST("/register", postRegister)
		auth.GET("/me", middleware.AuthMiddleware(), getMe)
		auth.POST("/logout", middleware.AuthMiddleware(), postLogout)
	}
}

func postLogin(c *gin.Context) {
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

	token, err := security.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func postRegister(c *gin.Context) {
	appCtx := c.MustGet("appCtx").(*context.AppContext)
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	exists, err := appCtx.UserRepo.UsernameExists(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username is already taken"})
		return
	}

	hash, err := security.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user := &repository.User{
		Username:    req.Username,
		Password:    hash,
		DisplayName: req.DisplayName,
		Email:       req.Email,
		IsActive:    true,
	}

	if err := appCtx.UserRepo.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": user.ID,
	})
}

func getMe(c *gin.Context) {
	appCtx := c.MustGet("appCtx").(*context.AppContext)
	claimsRaw, exists := c.Get("jwtClaims")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}
	claims := claimsRaw.(*security.Claims)

	userID := claims.UserID

	user, err := appCtx.UserRepo.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	roles, err := appCtx.UserRepo.GetUserRoles(userID)
	if err != nil {
		roles = []string{} // graceful fallback
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"display_name": user.DisplayName,
		"email":        user.Email,
		"avatar_url":   user.AvatarURL,
		"roles":        roles,
	})
}

func postLogout(c *gin.Context) {
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
