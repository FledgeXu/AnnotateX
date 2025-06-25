package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func JSONResponse(c *gin.Context, status int, message string, data any) {
	if data == nil {
		data = gin.H{}
	}
	c.JSON(status, Response{
		Code:    status,
		Message: message,
		Data:    data,
	})
}

func OK(c *gin.Context, data any) {
	JSONResponse(c, http.StatusOK, "Success", data)
}

func Created(c *gin.Context, data any) {
	JSONResponse(c, http.StatusCreated, "Created", data)
}

func BadRequest(c *gin.Context, message string) {
	JSONResponse(c, http.StatusBadRequest, message, nil)
}

func Unauthorized(c *gin.Context, message string) {
	JSONResponse(c, http.StatusUnauthorized, message, nil)
}

func Forbidden(c *gin.Context, message string) {
	JSONResponse(c, http.StatusForbidden, message, nil)
}

func NotFound(c *gin.Context, message string) {
	JSONResponse(c, http.StatusNotFound, message, nil)
}

func InternalServerError(c *gin.Context, message string) {
	JSONResponse(c, http.StatusInternalServerError, message, nil)
}

func AbortJSON(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, Response{
		Code:    status,
		Message: message,
		Data:    gin.H{},
	})
}
