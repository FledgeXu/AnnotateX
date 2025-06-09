package middleware

import (
	"net/http"

	"annotate-x/internal/context"
	"annotate-x/internal/security"

	"github.com/gin-gonic/gin"
)

func UserInjectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appCtx := c.MustGet("appCtx").(*context.AppContext)

		claimsRaw, exists := c.Get("jwtClaims")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "jwtClaims missing"})
			return
		}
		claims := claimsRaw.(*security.Claims)

		user, err := appCtx.UserRepo.GetUserByID(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.Set("currentUser", user)
		c.Next()
	}
}
