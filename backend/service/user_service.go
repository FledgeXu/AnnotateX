package service

import (
	"annotate-x/models"
	"annotate-x/repo"
	"annotate-x/utils/security"
	"errors"
)

type IUserService interface {
	Create(createRequest *models.CreateUserRequest) error
	GetUserById(userId int64) (*models.UserResponse, error)
}

type UserService struct {
	UserRepo repo.IUserRepo
}

func NewUserService(userRepo repo.IUserRepo) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) Create(createRequest *models.CreateUserRequest) error {
	isExist, err := s.UserRepo.UsernameExists(createRequest.Username)
	if err != nil {
		return err
	}
	if isExist {
		return errors.New("Username already exists.")
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

	_, err = s.UserRepo.CreateUser(newUser)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetUserById(userId int64) (*models.UserResponse, error) {
	user, err := s.UserRepo.GetUserByID(userId)
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
