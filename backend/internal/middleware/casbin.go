package middleware

import (
	"annotate-x/internal/security"

	"github.com/gin-gonic/gin"

	"fmt"
	"github.com/casbin/casbin/v2"
)

func CasbinMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsRaw, exists := c.Get("jwtClaims")
		if !exists {
			c.AbortWithStatusJSON(401, gin.H{"error": "unauthorized"})
			return
		}
		role := claimsRaw.(*security.Claims).Role // Get role from jwt claims.

		obj := c.FullPath() // API path，like /labels
		act := c.Request.Method
		fmt.Printf("[CASBIN] obj=%s role=%s act=%s\n", obj, role, act)

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
