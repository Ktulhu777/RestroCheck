package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims — структура, используемая в SSO
type Claims struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
	AppID int    `json:"app_id"`
	jwt.RegisteredClaims
}

// ParseToken парсит и валидирует JWT
func ParseToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что алгоритм подписи совпадает с ожидаемым
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	// Проверяем валидность токена
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Проверяем, не истёк ли токен
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token is expired")
	}

	return claims, nil
}
