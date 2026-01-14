package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/Everestown/Outfit_backend/internal/models"
	"github.com/Everestown/Outfit_backend/internal/modules/auth/dto"
	"github.com/Everestown/Outfit_backend/internal/modules/auth/repository"
	"github.com/Everestown/Outfit_backend/internal/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(req dto.RegisterRequest, ctx dto.SessionContext) (*dto.TokenResponse, error)
	Login(req dto.LoginRequest, ctx dto.SessionContext) (*dto.TokenResponse, error)
	RefreshToken(refreshToken string, ctx dto.SessionContext) (*dto.TokenResponse, error)
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

func (s *service) Register(req dto.RegisterRequest, ctx dto.SessionContext) (*dto.TokenResponse, error) {
	user, err := s.repo.GetUserByEmail(req.Email)
	if err == nil && user.ID != 0 {
		return nil, errors.New("user already exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	passHash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user = &models.User{
		Surname:      req.Surname,
		Name:         req.Name,
		Patronymic:   req.Patronymic,
		Username:     req.Username,
		Phone:        req.Phone,
		Email:        req.Email,
		PasswordHash: string(passHash),
		RoleID:       1,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	return s.tokens(user, ctx)
}

func (s *service) Login(req dto.LoginRequest, ctx dto.SessionContext) (*dto.TokenResponse, error) {
	user, err := s.repo.GetUserByIdentifier(req.Identifier)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	return s.tokens(user, ctx)
}

func (s *service) RefreshToken(refreshToken string, ctx dto.SessionContext) (*dto.TokenResponse, error) {
	session, err := s.repo.GetUserSession(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Rotation — старая сессия закрывается
	_ = s.repo.DeleteUserSession(session.JTI)

	user, err := s.repo.GetUserByID(session.UserID)
	if err != nil {
		return nil, err
	}

	return s.tokens(user, ctx)
}

func (s *service) Logout(userID uint) error {
	return s.repo.DeleteAllUserSessions(userID)
}

func (s *service) GetUserByID(userID uint) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}

func (s *service) tokens(user *models.User, ctx dto.SessionContext) (*dto.TokenResponse, error) {
	access, _ := s.jwt.GenerateAccessToken(user.ID, user.TokenVersion)
	refresh, _ := s.jwt.GenerateRefreshToken(user.ID)

	refreshHash := sha256Hash(refresh)

	session := &models.UserSession{
		UserID:           user.ID,
		RefreshTokenHash: refreshHash,
		JTI:              uuid.NewString(),
		IP:               &ctx.IP,
		DeviceInfo:       &ctx.DeviceInfo,
		ExpiresAt:        time.Now().Add(30 * 24 * time.Hour),
		Revoked:          false,
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

func sha256Hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
