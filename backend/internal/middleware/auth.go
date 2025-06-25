package middleware

import (
	ctxpkg "annotate-x/internal/context"
	"annotate-x/internal/security"
	"annotate-x/utils"
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
			utils.JSONAbortWithError(c, http.StatusUnauthorized, "Authorization header missing or malformed")
			return
		}

		tokenStr := authHeader[7:]

		ctx := context.Background()
		exist, err := appCtx.CacheRepo.IsBlacklisted(ctx, tokenStr)
		if err != nil && err != redis.Nil {
			utils.JSONAbortWithError(c, http.StatusInternalServerError, "Failed to check token blacklist")
			return
		}
		if exist {
			utils.JSONAbortWithError(c, http.StatusUnauthorized, "Token is blacklisted")
			return
		}

		claims, err := security.ParseToken(tokenStr)
		if err != nil {
			utils.JSONAbortWithError(c, http.StatusUnauthorized, "Invalid or expired token")
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
			utils.JSONAbortWithError(c, http.StatusUnauthorized, "Unauthorized")
			return
		}
		role := claimsRaw.(*security.Claims).Role // Get role from jwt claims.

		obj := c.FullPath() // API path，like /labels
		act := c.Request.Method

		allowed, err := enforcer.Enforce(obj, role, act)
		if err != nil {
			utils.JSONAbortWithError(c, http.StatusInternalServerError, err.Error())
			return
		}
		if !allowed {
			utils.JSONAbortWithError(c, http.StatusForbidden, "Forbidden")
			return
		}
		c.Next()
	}
}
