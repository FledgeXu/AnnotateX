package service

import (
	"net/http"
	"time"

	"annotate-x/internal/middleware"

	"annotate-x/internal/context"
	"annotate-x/repository"

	"annotate-x/internal/security"

	"annotate-x/model"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context, req model.UserCreateRequest) (*model.User, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Username: req.Username,
		Password: hashedPassword,
		Role:     req.Role,
	}

	err = model.InsertUser(ctx, &user)
	return &user, err
}
