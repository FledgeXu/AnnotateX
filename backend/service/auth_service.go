package service

import (
	"annotate-x/models"
	"annotate-x/repo"
	"annotate-x/utils/security"
	"context"
	"errors"
	"strconv"
	"time"
)

type IAuthService interface {
	Login(ctx context.Context, username, password string) (*models.User, string, error)
	Logout(ctx context.Context, userId int64, rawExpiration time.Time) error
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

func (s *AuthService) Login(ctx context.Context, username, password string) (*models.User, string, error) {
	user, err := s.UserRepo.GetUserByUsername(ctx, username)
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
			s.UserRepo.UpdateUserPassword(ctx, user.ID, newHash)
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

func (s *AuthService) Logout(ctx context.Context, userId int64, rawExpiration time.Time) error {
	expiration := time.Until(rawExpiration)
	if expiration <= 0 {
		return nil
	}

	err := s.CacheService.BlacklistToken(ctx, strconv.FormatInt(userId, 10), int(expiration.Seconds()))
	if err != nil {
		return errors.New("failed to logout")
	}

	return nil
}
