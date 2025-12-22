package dto

import "github.com/Everestown/Outfit_backend/internal/models"

type RegisterRequest struct {
	Surname    string `json:"surname" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Patronymic string `json:"patronymic"`
	Username   string `json:"username" binding:"required"`
	Phone      string `json:"phone"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"`
	User         *models.User `json:"user"`
}
