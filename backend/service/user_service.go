package service

import (
	"annotate-x/repository"

	"annotate-x/internal/security"

	"annotate-x/model"
	"errors"
)

var (
	ErrUsernameExists = errors.New("username exists.")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(req model.UserCreateRequest) (*repository.User, error) {
	exists, err := s.repo.UsernameExists(req.Username)
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

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, err
}
