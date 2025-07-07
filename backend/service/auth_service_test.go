package service_test

import (
	"annotate-x/mocks"
	"annotate-x/models"
	"annotate-x/service"
	"annotate-x/utils/security"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Login_Success(t *testing.T) {
	userRepo := mocks.NewMockIUserRepository(t)
	cacheService := mocks.NewMockICacheService(t)

	password := "secret"
	hashed, _ := security.HashPassword(password)

	user := &models.User{
		ID:       1,
		Username: "testuser",
		Password: hashed,
	}

	userRepo.On("GetUserByUsername", "testuser").Return(user, nil)
	userRepo.On("UpdateUserPassword", user.ID, mock.Anything).Return(nil).Maybe()

	authService := service.NewAuthService(userRepo, cacheService)

	loggedInUser, token, err := authService.Login("testuser", password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, user.Username, loggedInUser.Username)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	userRepo := mocks.NewMockIUserRepository(t)
	cacheService := mocks.NewMockICacheService(t)

	hashed, _ := security.HashPassword("correct-password")
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Password: hashed,
	}

	userRepo.On("GetUserByUsername", "testuser").Return(user, nil)

	authService := service.NewAuthService(userRepo, cacheService)

	_, _, err := authService.Login("testuser", "wrong-password")

	assert.Error(t, err)
	assert.Equal(t, "Invalid username or password", err.Error())
	userRepo.AssertExpectations(t)
}

func TestAuthService_Register_Success(t *testing.T) {
	userRepo := mocks.NewMockIUserRepository(t)
	cacheService := mocks.NewMockICacheService(t)

	userRepo.On("UsernameExists", "newuser").Return(false, nil)
	userRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(int64(1), nil)

	authService := service.NewAuthService(userRepo, cacheService)

	req := &models.CreateUserRequest{
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
	cacheService := mocks.NewMockICacheService(t)

	userRepo.On("UsernameExists", "existinguser").Return(true, nil)

	authService := service.NewAuthService(userRepo, cacheService)

	req := &models.CreateUserRequest{
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

func TestAuthService_Logout_Success(t *testing.T) {
	userRepo := mocks.NewMockIUserRepository(t)
	cacheService := mocks.NewMockICacheService(t)

	authService := service.NewAuthService(userRepo, cacheService)

	tokenStr, _ := security.GenerateToken(1, "testuser")

	claims, _ := security.ParseToken(tokenStr)
	expiration := time.Until(claims.ExpiresAt.Time)

	cacheService.On("BlacklistToken", tokenStr, mock.MatchedBy(func(i int) bool {
		return float64(i) > expiration.Seconds()-2 && float64(i) < expiration.Seconds()+2
	})).Return(nil)

	err := authService.Logout(tokenStr)

	assert.NoError(t, err)
	cacheService.AssertExpectations(t)
}

func TestAuthService_Logout_InvalidToken(t *testing.T) {
	userRepo := mocks.NewMockIUserRepository(t)
	cacheService := mocks.NewMockICacheService(t)

	authService := service.NewAuthService(userRepo, cacheService)

	tokenStr := "invalid.token.string"

	err := authService.Logout(tokenStr)

	assert.Error(t, err)
	assert.Equal(t, "invalid token", err.Error())
}

func TestAuthService_Logout_BlacklistError(t *testing.T) {
	userRepo := mocks.NewMockIUserRepository(t)
	cacheService := mocks.NewMockICacheService(t)

	authService := service.NewAuthService(userRepo, cacheService)

	tokenStr, _ := security.GenerateToken(1, "testuser")

	cacheService.On("BlacklistToken", tokenStr, mock.AnythingOfType("int")).Return(errors.New("redis down"))

	err := authService.Logout(tokenStr)

	assert.Error(t, err)
	assert.Equal(t, "failed to logout", err.Error())
	cacheService.AssertExpectations(t)
}
