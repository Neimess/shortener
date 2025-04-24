package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	cachepkg "github.com/Neimess/shortener/internal/infrastructure/cache"
	"github.com/Neimess/shortener/internal/infrastructure/db"
	jwtutil "github.com/Neimess/shortener/internal/util/jwt"
	userValidator "github.com/Neimess/shortener/internal/validators/user"
	"golang.org/x/crypto/bcrypt"
)

type serviceImpl struct {
	repos      db.RepositorySet
	cache      cachepkg.ExpirableCache
	jwtManager *jwtutil.JWTManager
}

func NewService(
	repos db.RepositorySet,
	cache cachepkg.ExpirableCache,
	jwtSecret []byte,
	jwtExpiry time.Duration,
	refreshExpiry time.Duration,

) Service {
	return &serviceImpl{
		repos: repos,
		cache: cache,
		jwtManager: jwtutil.New(
			jwtSecret,
			jwtExpiry,
			refreshExpiry),
	}
}

func (s *serviceImpl) Register(ctx context.Context, email, password string) (int, error) {

	input := &userValidator.RegisterInput{Email: email, Password: password}
	if err := input.Validate(); err != nil {
		return 0, fmt.Errorf("validation failed: %w", err)
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	userID, err := s.repos.User.Create(ctx, email, string(passwordHash))
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (s *serviceImpl) Login(ctx context.Context, email, password string) (string, string, error) {
	input := &userValidator.LoginInput{Email: email, Password: password}
	if err := input.Validate(); err != nil {
		return "", "", fmt.Errorf("validation failed: %w", err)
	}

	user, err := s.repos.User.GetByEmail(ctx, email)
	if err != nil {
		return "", "", fmt.Errorf("user not found: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", fmt.Errorf("invalid password")
	}
	userID := strconv.Itoa(user.ID)
	accessToken, err := s.jwtManager.GenerateAccessToken(userID, user.Email)
	if err != nil {
		return "", "", fmt.Errorf("access token error: %w", err)
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(userID, user.Email)
	if err != nil {
		return "", "", fmt.Errorf("access token error: %w", err)
	}

	key := "refresh:" + refreshToken
	if err := s.cache.SetWithTTL(ctx, key, userID, s.jwtManager.RefreshTTL()); err != nil {
		return "", "", fmt.Errorf("cannot persist refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *serviceImpl) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := s.jwtManager.Decode(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token: %w", err)
	}

	keyOld := "refresh:" + refreshToken
	storedUserID, err := s.cache.Get(ctx, keyOld)
	if err != nil || storedUserID == "" {
		return "", "", errors.New("refresh token revoked or expired")
	}
	_ = s.cache.Del(ctx, keyOld)

	userID := claims.Subject
	email := claims.Email

	newAt, err := s.jwtManager.GenerateAccessToken(userID, email)
	if err != nil {
		return "", "", fmt.Errorf("cannot generate new access token: %w", err)
	}
	newRt, err := s.jwtManager.GenerateRefreshToken(userID, email)
	if err != nil {
		return "", "", fmt.Errorf("cannot generate new refresh token: %w", err)
	}

	keyNew := "refresh:" + newRt
	if err := s.cache.SetWithTTL(ctx, keyNew, userID, s.jwtManager.RefreshTTL()); err != nil {
		return "", "", fmt.Errorf("cannot persist new refresh token: %w", err)
	}

	return newAt, newRt, nil
}
