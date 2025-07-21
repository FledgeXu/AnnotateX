package service

import (
	"annotate-x/models"
	"annotate-x/repo"
	"annotate-x/utils/security"
	"context"
	"errors"
)

type IUserService interface {
	Create(ctx context.Context, createRequest *models.CreateUserRequest) error
	GetUserById(ctx context.Context, userId int64) (*models.UserResponse, error)
}

type UserService struct {
	UserRepo repo.IUserRepo
}

func NewUserService(userRepo repo.IUserRepo) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) Create(ctx context.Context, createRequest *models.CreateUserRequest) error {
	isExist, err := s.UserRepo.UsernameExists(ctx, createRequest.Username)
	if err != nil {
		return err
	}
	if isExist {
		return errors.New("username already exists")
	}

	hashedPassword, err := security.HashPassword(createRequest.Password)
	if err != nil {
		return err
	}

	newUser := &models.User{
		Username:    createRequest.Username,
		Password:    hashedPassword,
		DisplayName: createRequest.DisplayName,
		Email:       createRequest.Email,
		IsActive:    true,
	}

	err = s.UserRepo.CreateUser(ctx, newUser)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetUserById(ctx context.Context, userId int64) (*models.UserResponse, error) {
	user, err := s.UserRepo.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}
	userResp := &models.UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	return userResp, nil
}
