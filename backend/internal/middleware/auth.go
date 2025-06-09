package middleware

import (
	ctxpkg "annotate-x/internal/context"
	"annotate-x/internal/security"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/casbin/casbin/v2"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appCtx := c.MustGet("appCtx").(*ctxpkg.AppContext)
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or malformed"})
			return
		}

		tokenStr := authHeader[7:]

		ctx := context.Background()
		exist, err := appCtx.CacheRepo.IsBlacklisted(ctx, tokenStr)
		if err != nil && err != redis.Nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to check token blacklist"})
			return
		}
		if exist {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted"})
			return
		}

		claims, err := security.ParseToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("jwtClaims", claims)
		c.Set("rawToken", tokenStr)

		c.Next()
	}
}

func RequirePermissionMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsRaw, exists := c.Get("jwtClaims")
		if !exists {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		role := claimsRaw.(*security.Claims).Role // Get role from jwt claims.

		obj := c.FullPath() // API path，like /labels
		act := c.Request.Method

		allowed, err := enforcer.Enforce(obj, role, act)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
