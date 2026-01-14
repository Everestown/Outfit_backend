package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UUID         string         `gorm:"type:uuid;uniqueIndex;default:gen_random_uuid()" json:"uuid"`
	Surname      string         `gorm:"type:varchar(50);not null" json:"surname" validate:"required,min=2,max=50"`
	Name         string         `gorm:"type:varchar(50);not null" json:"name" validate:"required,min=2,max=50"`
	Patronymic   string         `gorm:"type:varchar(50)" json:"patronymic" validate:"min=2,max=50"`
	Username     string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username" validate:"required,min=2,max=50"`
	Phone        string         `gorm:"type:varchar(25);uniqueIndex" json:"phone" validate:"min=5,max=25,numeric"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" validate:"required,email"`
	PasswordHash string         `gorm:"type:varchar(512);not null" json:"-"`
	RoleID       uint           `gorm:"not null" json:"role_id"`
	TokenVersion uint           `gorm:"not null;default:0" json:"token_version"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (User) TableName() string {
	return "store.users"
}

type UserSession struct {
	ID               string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID           uint      `gorm:"not null;index"`
	RefreshTokenHash string    `gorm:"type:varchar(1024);not null"`
	JTI              string    `gorm:"type:varchar(256);uniqueIndex;not null"`
	IP               *string   `gorm:"type:inet"`
	DeviceInfo       *string   `gorm:"type:varchar(256)"`
	CreatedAt        time.Time `gorm:"not null;default:current_timestamp"`
	LastUsedAt       time.Time `gorm:"not null;default:current_timestamp"`
	ExpiresAt        time.Time `gorm:"not null"`
	Revoked          bool      `gorm:"not null;default:false"`
}

func (UserSession) TableName() string {
	return "store.user_sessions"
}
