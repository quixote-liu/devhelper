package service

import (
	"errors"
	"time"

	"devhelper/internal/config"
	"devhelper/internal/models"
	"devhelper/internal/repository"
	"devhelper/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repository.UserRepo
	cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepo, cfg *config.Config) *AuthService {
	return &AuthService{userRepo: userRepo, cfg: cfg}
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s *AuthService) Register(username, email, password string) (*models.User, *TokenPair, error) {
	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return nil, nil, errors.New("email already registered")
	}
	if _, err := s.userRepo.FindByUsername(username); err == nil {
		return nil, nil, errors.New("username already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, nil, err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		Role:         "user",
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, err
	}

	tokens, err := s.generateTokens(user)
	return user, tokens, err
}

func (s *AuthService) Login(email, password string) (*models.User, *TokenPair, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("invalid credentials")
		}
		return nil, nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	now := time.Now()
	user.LastLogin = &now
	_ = s.userRepo.Update(user)

	tokens, err := s.generateTokens(user)
	return user, tokens, err
}

func (s *AuthService) Refresh(refreshToken string) (*TokenPair, error) {
	claims, err := utils.ParseToken(refreshToken, s.cfg.JWTSecret)
	if err != nil || claims.Type != "refresh" {
		return nil, errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.generateTokens(user)
}

func (s *AuthService) generateTokens(user *models.User) (*TokenPair, error) {
	access, err := utils.GenerateAccessToken(user.ID, user.Email, user.Role, s.cfg.JWTSecret, s.cfg.JWTAccessExpiry)
	if err != nil {
		return nil, err
	}
	refresh, err := utils.GenerateRefreshToken(user.ID, user.Email, user.Role, s.cfg.JWTSecret, s.cfg.JWTRefreshExpiry)
	if err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (s *AuthService) GetUser(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}
