package middleware

import (
	"net/http"

	"annotate-x/internal/context"
	"annotate-x/internal/security"
	"annotate-x/utils"

	"github.com/gin-gonic/gin"
)

func UserInjectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		appCtx := c.MustGet("appCtx").(*context.AppContext)

		claimsRaw, exists := c.Get("jwtClaims")
		if !exists {
			utils.AbortJSON(c, http.StatusUnauthorized, "JWT claims missing")
			return
		}
		claims := claimsRaw.(*security.Claims)

		user, err := appCtx.UserRepo.GetUserByID(claims.UserID)
		if err != nil {
			utils.AbortJSON(c, http.StatusNotFound, "User not found")
			return
		}

		c.Set("currentUser", user)
		c.Next()
	}
}
