package utils

import (
	"bilet/backend/code"
	"bilet/backend/config"
	"bilet/backend/models"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

func NewRefreshToken(Email string, role string, userId int) (string, *code.ResultCode) {
	tokenId := "refresh-" + uuid.NewString()
	return NewToken(role, tokenId, Email, config.RefreshTokenExpiration, userId)
}

func NewAccessToken(Email string, role string, userId int) (string, *code.ResultCode) {
	tokenId := "access-" + uuid.NewString()
	return NewToken(role, tokenId, Email, config.AccessTokenExpiration, userId)
}

func NewToken(role string, jti string, Email string, expiration time.Duration, userId int) (string, *code.ResultCode) {
	tokenExpirationTime := time.Now().Add(expiration)
	tokenClaims := &models.Claims{
		Role:   role,
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			Id:        jti,
			Subject:   Email,
			ExpiresAt: tokenExpirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenString, err := token.SignedString(config.JwtKey)
	if err != nil {
		return "", code.InternalServerError.SetMessage("Could not generate token")
	}

	dbToken := &models.Token{
		JTI:       jti,
		Subject:   Email,
		IsRevoked: false,
		ExpiresAt: tokenExpirationTime.Unix(),
	}
	models.DB.Create(dbToken)

	return tokenString, nil
}
