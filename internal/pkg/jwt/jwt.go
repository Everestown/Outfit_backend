package jwt

import (
	"time"

	"github.com/Everestown/Outfit_backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JWTManager struct {
	secret string
	db     *gorm.DB // Для проверки token_version
}

type Claims struct {
	UserID uint `json:"user_id"`
	Ver    uint `json:"ver"`
	jwt.RegisteredClaims
}

func NewJWTManager(secret string, db *gorm.DB) *JWTManager {
	return &JWTManager{
		secret: secret,
		db:     db,
	}
}

func (j *JWTManager) GenerateAccessToken(userID uint, tokenVersion uint) (string, error) {
	claims := Claims{
		UserID: userID,
		Ver:    tokenVersion,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWTManager) GenerateRefreshToken(userID uint) (string, error) {
	jti := uuid.New().String()
	claims := jwt.MapClaims{
		"sub": userID,
		"jti": jti,
		"exp": time.Now().Add(30 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	var user models.User
	if err := j.db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		return nil, err
	}

	if claims.Ver != user.TokenVersion {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
