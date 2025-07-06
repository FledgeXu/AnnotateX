package middleware

import (
	"net/http"

	"annotate-x/internal/security"
	"annotate-x/repository"
	"annotate-x/utils"

	"github.com/gin-gonic/gin"
)

func UserInjectionMiddleware(userRepo *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		claimsRaw, exists := c.Get("jwtClaims")
		if !exists {
			utils.AbortJSON(c, http.StatusUnauthorized, "JWT claims missing")
			return
		}
		claims := claimsRaw.(*security.Claims)

		user, err := userRepo.GetUserByID(claims.UserID)
		if err != nil {
			utils.AbortJSON(c, http.StatusNotFound, "User not found")
			return
		}

		c.Set("currentUser", user)
		c.Next()
	}
}
