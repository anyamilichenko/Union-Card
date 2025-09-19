package utils

import (
	"bilet/backend/config"
	"bilet/backend/models"
	"github.com/golang-jwt/jwt"
)

func ParseToken(tokenString string) (claims *models.Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return config.JwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
