package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

func JSON[T any](c *gin.Context, status int, message string, data T) {
	c.JSON(status, Response[T]{
		Code:    status,
		Message: message,
		Data:    data,
	})
}

func AbortJSON[T any](c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, Response[T]{
		Code:    status,
		Message: message,
		Data:    *new(T),
	})
}

func OK[T any](c *gin.Context, data T) {
	JSON(c, http.StatusOK, "Success", data)
}

func Created[T any](c *gin.Context, data T) {
	JSON(c, http.StatusCreated, "Created", data)
}

func Error(c *gin.Context, status int, message string) {
	JSON(c, status, message, gin.H{})
}

func AbortError(c *gin.Context, status int, message string) {
	AbortJSON[gin.H](c, status, message)
}

func BadRequest(c *gin.Context, msg string)   { Error(c, http.StatusBadRequest, msg) }
func Unauthorized(c *gin.Context, msg string) { Error(c, http.StatusUnauthorized, msg) }
func Forbidden(c *gin.Context, msg string)    { Error(c, http.StatusForbidden, msg) }
func NotFound(c *gin.Context, msg string)     { Error(c, http.StatusNotFound, msg) }
func InternalServerError(c *gin.Context, msg string) {
	Error(c, http.StatusInternalServerError, msg)
}
