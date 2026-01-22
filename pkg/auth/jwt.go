package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

// Claims представляет данные, хранимые в JWT токене
type Claims struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTManager управляет созданием и валидацией JWT токенов
type JWTManager struct {
	secretKey []byte
	ttl       time.Duration
}

// NewJWTManager создает новый экземпляр JWT менеджера
func NewJWTManager(secretKey string, ttlHours int) *JWTManager {
	return &JWTManager{
		secretKey: []byte(secretKey),
		ttl:       time.Duration(ttlHours) * time.Hour,
	}
}

// GenerateToken создает новый JWT токен для пользователя
func (m *JWTManager) GenerateToken(userID int, email, username string) (string, time.Time, error) {
	// 1. Создать Claims с данными пользователя
	expiredAt := time.Now().Add(m.ttl)
	claims := &Claims{
		UserID:   userID,
		Email:    email,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "blog-api",
			Subject:   "user",
		},
	}

	// 2. Создать токен используя алгоритм подписи (HS256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 3. Подписать токен секретным ключом
	tokenString, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	// 4. Вернуть подписанную строку токена и время истечения
	return tokenString, expiredAt, nil
}

// ValidateToken проверяет и парсит JWT токен
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	// 1. Распарсить токен с проверкой подписи
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return m.secretKey, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	// 2. Извлечь claims из токена
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// 3. Проверить время истечения токена
	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, ErrExpiredToken
	}

	// 4. Вернуть claims, если токен валидный
	return claims, nil
}
