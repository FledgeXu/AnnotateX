package middleware

import (
	"annotate-x/errors"
	"annotate-x/utils"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			switch e := err.(type) {
			case *errors.BadRequestError:
				utils.BadRequest(c, e.Error())
			default:
				utils.InternalServerError(c, e.Error())
			}
		}
	}
}
