package service

import (
	"annotate-x/model"
	"annotate-x/repo"
	"annotate-x/util/security"
	"errors"
)

type IAuthService interface {
	Login(username, password string) (*model.User, string, error)
	Register(createRequest *model.CreateUserRequest) error
	Logout(userID int64) error
}

type AuthService struct {
	UserRepo  repo.IUserRepository
	CacheRepo repo.ICacheRepository
}

func NewAuthService(userRepo repo.IUserRepository, cacheRepo repo.ICacheRepository) *AuthService {
	return &AuthService{
		UserRepo:  userRepo,
		CacheRepo: cacheRepo,
	}
}

func (s *AuthService) Login(username, password string) (*model.User, string, error) {
	user, err := s.UserRepo.GetUserByUsername(username)
	if err != nil {
		return nil, "", errors.New("Invalid username or password")
	}

	match, needsRehash, err := security.VerifyPassword(password, user.Password)
	if err != nil {
		return nil, "", err
	}

	// Auto-upgrade hash if parameters are outdated
	if needsRehash {
		if newHash, ok, err := security.RehashIfNeeded(password, user.Password); err == nil && ok {
			user.Password = newHash
			s.UserRepo.UpdateUserPassword(user.ID, newHash)
		}
	}

	if !match {
		return nil, "", errors.New("Invalid username or password")
	}

	token, err := security.GenerateToken(user.ID, user.Username)

	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Register(createRequest *model.CreateUserRequest) error {
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

	newUser := &model.User{
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
