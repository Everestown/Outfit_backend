package service

import (
	"errors"
	"time"

	"github.com/Everestown/Outfit_backend/internal/models"
	"github.com/Everestown/Outfit_backend/internal/modules/auth/dto"
	"github.com/Everestown/Outfit_backend/internal/modules/auth/repository"
	"github.com/Everestown/Outfit_backend/internal/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(req dto.RegisterRequest) (*dto.TokenResponse, error)
	Login(req dto.LoginRequest) (*dto.TokenResponse, error)
	RefreshToken(refreshToken string) (*dto.TokenResponse, error)
	Logout(userID uint) error
	GetUserByID(userID uint) (*models.User, error)
}

type service struct {
	repo repository.Repository
	jwt  *jwt.JWTManager
}

func NewService(repo repository.Repository, jwtManager *jwt.JWTManager) Service {
	return &service{repo: repo, jwt: jwtManager}
}

func (s *service) Register(req dto.RegisterRequest) (*dto.TokenResponse, error) {
	if user, _ := s.repo.GetUserByEmail(req.Email); user != nil {
		return nil, errors.New("user already exists")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := &models.User{
		Surname:      req.Surname,
		Name:         req.Name,
		Patronymic:   req.Patronymic,
		Username:     req.Username,
		Phone:        req.Phone,
		Email:        req.Email,
		PasswordHash: string(hash),
		RoleID:       1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return s.tokens(user)
}

func (s *service) Login(req dto.LoginRequest) (*dto.TokenResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.tokens(user)
}

func (s *service) RefreshToken(refreshToken string) (*dto.TokenResponse, error) {
	session, err := s.repo.GetUserSession(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	_ = s.repo.DeleteUserSession(session.JTI)

	user, _ := s.repo.GetUserByID(session.UserID)
	return s.tokens(user)
}

func (s *service) Logout(userID uint) error {
	return s.repo.DeleteAllUserSessions(userID)
}

func (s *service) GetUserByID(userID uint) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}

func (s *service) tokens(user *models.User) (*dto.TokenResponse, error) {
	access, _ := s.jwt.GenerateAccessToken(user.ID, user.TokenVersion)
	refresh, _ := s.jwt.GenerateRefreshToken(user.ID)

	session := &models.UserSession{
		UserID:           user.ID,
		RefreshTokenHash: hash(refresh),
		JTI:              uuid.NewString(),
		ExpiresAt:        time.Now().Add(30 * 24 * time.Hour),
	}

	_ = s.repo.CreateUserSession(session)

	return &dto.TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		TokenType:    "Bearer",
		ExpiresIn:    900,
		User:         user,
	}, nil
}

func hash(token string) string {
	h, _ := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	return string(h)
}
