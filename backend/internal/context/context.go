package context

import (
	"annotate-x/repository"

	"github.com/gin-gonic/gin"
)

type AppContext struct {
	UserRepo  *repository.UserRepository
	CacheRepo *repository.CacheRepository
}

func InjectAppContext(app *AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("appCtx", app)
		c.Next()
	}
}
