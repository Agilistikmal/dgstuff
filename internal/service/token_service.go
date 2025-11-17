package service

import (
	"time"

	"github.com/agilistikmal/dgstuff/internal/app"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type TokenService struct {
	secretKey string
}

func NewTokenService() *TokenService {
	return &TokenService{secretKey: viper.GetString("token.secret_key")}
}

func (s *TokenService) GenerateToken(data map[string]any, expiration time.Duration) string {
	claims := jwt.MapClaims{
		"data": data,
	}
	if expiration > 0 {
		claims["exp"] = time.Now().Add(expiration).Unix()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		logrus.Errorf("failed to generate token: %v", err)
		return ""
	}
	return tokenString
}

func (s *TokenService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, app.NewForbiddenError("invalid token")
	}
	if !token.Valid {
		return nil, app.NewForbiddenError("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, app.NewForbiddenError("invalid token")
	}
	if claims["exp"] != nil && claims["exp"].(float64) < float64(time.Now().Unix()) {
		return nil, app.NewForbiddenError("token expired")
	}
	return token.Claims.(jwt.MapClaims), nil
}
