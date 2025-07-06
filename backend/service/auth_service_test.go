package service_test

import (
	"annotate-x/mocks"
	"annotate-x/model"
	"annotate-x/service"
	"annotate-x/util/security"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Login_Success(t *testing.T) {
	userRepo := mocks.NewMockIUserRepository(t)
	cacheRepo := mocks.NewMockICacheRepository(t)

	password := "secret"
	hashed, _ := security.HashPassword(password)

	user := &model.User{
		ID:       1,
		Username: "testuser",
		Password: hashed,
	}

	userRepo.On("GetUserByUsername", "testuser").Return(user, nil)
	userRepo.On("UpdateUserPassword", user.ID, mock.Anything).Return(nil).Maybe()

	authService := service.NewAuthService(userRepo, cacheRepo)

	loggedInUser, token, err := authService.Login("testuser", password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, user.Username, loggedInUser.Username)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	userRepo := mocks.NewMockIUserRepository(t)
	cacheRepo := mocks.NewMockICacheRepository(t)

	hashed, _ := security.HashPassword("correct-password")
	user := &model.User{
		ID:       1,
		Username: "testuser",
		Password: hashed,
	}

	userRepo.On("GetUserByUsername", "testuser").Return(user, nil)

	authService := service.NewAuthService(userRepo, cacheRepo)

	_, _, err := authService.Login("testuser", "wrong-password")

	assert.Error(t, err)
	assert.Equal(t, "Invalid username or password", err.Error())
	userRepo.AssertExpectations(t)
}

func TestAuthService_Register_Success(t *testing.T) {
	userRepo := mocks.NewMockIUserRepository(t)
	cacheRepo := mocks.NewMockICacheRepository(t)

	userRepo.On("UsernameExists", "newuser").Return(false, nil)
	userRepo.On("CreateUser", mock.AnythingOfType("*model.User")).Return(int64(1), nil)

	authService := service.NewAuthService(userRepo, cacheRepo)

	req := &model.CreateUserRequest{
		Username:    "newuser",
		Password:    "pass123",
		DisplayName: "New User",
		Email:       "new@example.com",
	}

	err := authService.Register(req)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Register_UsernameExists(t *testing.T) {
	userRepo := mocks.NewMockIUserRepository(t)
	cacheRepo := mocks.NewMockICacheRepository(t)

	userRepo.On("UsernameExists", "existinguser").Return(true, nil)

	authService := service.NewAuthService(userRepo, cacheRepo)

	req := &model.CreateUserRequest{
		Username:    "existinguser",
		Password:    "pass123",
		DisplayName: "Existing User",
		Email:       "existing@example.com",
	}

	err := authService.Register(req)

	assert.Error(t, err)
	assert.Equal(t, "Username already exists.", err.Error())
	userRepo.AssertExpectations(t)
}
