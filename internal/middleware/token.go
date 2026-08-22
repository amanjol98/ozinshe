package middleware

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secretKey []byte
}

func NewTokenService(secretKey string) *TokenService {
	return &TokenService{secretKey: []byte(secretKey)}
}

type TokenClaims struct {
	UserID int
	Role   string
}

func (s *TokenService) GenerateToken(userID int, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *TokenService) ValidateToken(tokenString string) (TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("Неверный метод подписи")
		}
		return s.secretKey, nil
	})

	if err != nil {
		return TokenClaims{}, err
	}

	if !token.Valid {
		return TokenClaims{}, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return TokenClaims{}, errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return TokenClaims{}, errors.New("invalid user_id")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return TokenClaims{}, errors.New("invalid role")
	}

	tokenClaims := TokenClaims{
		UserID: int(userID),
		Role:   role,
	}

	return tokenClaims, nil
}
