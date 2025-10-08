package utils

import (
	"errors"
	"log"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtSecret     []byte
	jwtSecretOnce sync.Once
)

func getJWTSecret() []byte {
	jwtSecretOnce.Do(func() {
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			log.Fatal("JWT_SECRET is not set")
		}
		jwtSecret = []byte(secret)
	})
	return jwtSecret
}

// HashPassword 对密码进行bcrypt加密
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// CheckPassword 验证密码是否正确
func CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// JWT Claims结构
type Claims struct {
	UserID primitive.ObjectID `json:"user_id"`
	Email  string             `json:"email"`
	jwt.RegisteredClaims
}

// GenerateJWT 生成JWT token
func GenerateJWT(userID primitive.ObjectID, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24小时有效期

	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// ParseJWT 解析JWT token
func ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// TemporaryClaims 临时token的Claims结构
type TemporaryClaims struct {
	Email   string `json:"email"`
	Purpose string `json:"purpose"` // "registration" 或 "password_reset"
	jwt.RegisteredClaims
}

// GenerateTemporaryToken 生成临时token（5分钟有效期）
// 用于注册或密码重置验证的两步操作
func GenerateTemporaryToken(email, purpose string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute) // 5分钟有效期

	claims := &TemporaryClaims{
		Email:   email,
		Purpose: purpose,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

// ParseTemporaryToken 解析临时token
func ParseTemporaryToken(tokenString string) (*TemporaryClaims, error) {
	claims := &TemporaryClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid temporary token")
	}

	return claims, nil
}
