package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/Everestown/Outfit_backend/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user *models.User) error
	GetUserByIdentifier(identifier string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	CreateUserSession(session *models.UserSession) error
	GetUserSession(refreshToken string) (*models.UserSession, error)
	DeleteUserSession(jti string) error
	DeleteAllUserSessions(userID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *repository) GetUserByIdentifier(identifier string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ? OR username = ?", identifier, identifier).First(&user).Error
	return &user, err
}

func (r *repository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *repository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *repository) CreateUserSession(session *models.UserSession) error {
	return r.db.Create(session).Error
}

func (r *repository) GetUserSession(refreshToken string) (*models.UserSession, error) {
	var session models.UserSession

	hash := sha256Hash(refreshToken)

	err := r.db.Where(
		"refresh_token_hash = ? AND revoked = false AND expires_at > ?",
		hash, time.Now(),
	).First(&session).Error

	return &session, err
}

func (r *repository) DeleteUserSession(jti string) error {
	return r.db.Model(&models.UserSession{}).
		Where("jti = ?", jti).
		Update("revoked", true).Error
}

func (r *repository) DeleteAllUserSessions(userID uint) error {
	return r.db.Model(&models.UserSession{}).
		Where("user_id = ?", userID).
		Update("revoked", true).Error
}

func sha256Hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
