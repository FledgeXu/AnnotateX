package service

import (
	"annotate-x/internal/context"
	"annotate-x/repository"

	"annotate-x/internal/security"

	"annotate-x/model"
	"errors"
)

var (
	ErrUsernameExists = errors.New("username exists.")
)

func CreateUser(appCtx *context.AppContext, req model.UserCreateRequest) (*repository.User, error) {
	exists, err := appCtx.UserRepo.UsernameExists(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrUsernameExists
	}

	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &repository.User{
		Username:    req.Username,
		Password:    hashedPassword,
		DisplayName: req.DisplayName,
		Email:       req.Email,
		IsActive:    true,
		Role:        req.Role,
	}

	if err := appCtx.UserRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, err
}
