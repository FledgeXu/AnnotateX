package utils

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func JSONResponse(c *gin.Context, httpCode int, message string, data any) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(httpCode, Response{
		Code:    httpCode,
		Message: message,
		Data:    data,
	})
}

func JSONSuccess(c *gin.Context, httpCode int, data any) {
	JSONResponse(c, httpCode, "Success", data)
}

func JSONError(c *gin.Context, httpCode int, message string) {
	JSONResponse(c, httpCode, message, nil)
}

func JSONAbortWithError(c *gin.Context, httpCode int, message string) {
	c.AbortWithStatusJSON(httpCode, Response{
		Code:    httpCode,
		Message: message,
		Data:    gin.H{},
	})
}
