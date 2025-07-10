package service_test

import (
	"annotate-x/mocks"
	"annotate-x/models"
	"annotate-x/service"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Create_Success(t *testing.T) {
	userRepo := mocks.NewMockIUserRepo(t)
	context := context.Background()

	userRepo.On("UsernameExists", mock.Anything, "newuser").Return(false, nil)
	userRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(int64(1), nil)

	userService := service.NewUserService(userRepo)

	req := &models.CreateUserRequest{
		Username:    "newuser",
		Password:    "pass123",
		DisplayName: "New User",
		Email:       "new@example.com",
	}

	err := userService.Create(context, req)

	assert.NoError(t, err)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Create_UsernameExists(t *testing.T) {
	userRepo := mocks.NewMockIUserRepo(t)
	context := context.Background()

	userRepo.On("UsernameExists", mock.Anything, "existinguser").Return(true, nil)

	userService := service.NewUserService(userRepo)

	req := &models.CreateUserRequest{
		Username:    "existinguser",
		Password:    "pass123",
		DisplayName: "Existing User",
		Email:       "existing@example.com",
	}

	err := userService.Create(context, req)

	assert.Error(t, err)
	assert.Equal(t, "Username already exists.", err.Error())
	userRepo.AssertExpectations(t)
}
