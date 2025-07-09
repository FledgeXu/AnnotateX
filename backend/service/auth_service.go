package service

import (
	"annotate-x/models"
	"annotate-x/repo"
	"annotate-x/utils/security"
	"errors"
	"time"
)

type IAuthService interface {
	Login(username, password string) (*models.User, string, error)
	Logout(tokenStr string) error
}

type AuthService struct {
	UserRepo     repo.IUserRepo
	CacheService ICacheService
}

func NewAuthService(userRepo repo.IUserRepo, CacheService ICacheService) *AuthService {
	return &AuthService{
		UserRepo:     userRepo,
		CacheService: CacheService,
	}
}

func (s *AuthService) Login(username, password string) (*models.User, string, error) {
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

func (s *AuthService) Logout(tokenStr string) error {
	claims, err := security.ParseToken(tokenStr)
	if err != nil {
		return errors.New("invalid token")
	}

	expiration := time.Until(claims.ExpiresAt.Time)
	if expiration <= 0 {
		return nil
	}

	err = s.CacheService.BlacklistToken(tokenStr, int(expiration.Seconds()))
	if err != nil {
		return errors.New("failed to logout")
	}

	return nil
}
