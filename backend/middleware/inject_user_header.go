package middleware

import (
	"annotate-x/models"
	"annotate-x/utils/security"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func InjectUserHeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return
		}

		tokenStr := strings.TrimSpace(authHeader[7:])
		if tokenStr == "" {
			return
		}

		claims, err := security.ParseToken(tokenStr)
		if err != nil {
			return
		}

		if claims != nil {
			c.Request.Header.Set(models.XUserID, strconv.FormatInt(claims.UserID, 10))
			c.Request.Header.Set(models.XUserName, claims.Username)
			c.Request.Header.Set(models.XExpiresAt, strconv.FormatInt(claims.ExpiresAt.Time.Unix(), 10))
		}

		c.Next()
	}
}
