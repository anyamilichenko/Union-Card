package utils

import (
	"bilet/backend/code"
	"bilet/backend/config"
	"bilet/backend/entity"
	"bilet/backend/repository"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"time"
)

type JWTUtil interface {
	NewRefreshToken(email, role string, userId uint) (string, *code.ResultCode)
	NewAccessToken(email, role string, userId uint) (string, *code.ResultCode)
	ParseToken(tokenString string) (*entity.Claims, error)
}

type jwtUtil struct {
	jwtKey      []byte
	accountRepo repository.AccountRepository
	tokenRepo   repository.TokenRepository
}

func NewJWTUtil(jwtKey []byte, accountRepo repository.AccountRepository, tokenRepo repository.TokenRepository) JWTUtil {
	return &jwtUtil{
		jwtKey:      jwtKey,
		accountRepo: accountRepo,
		tokenRepo:   tokenRepo,
	}
}

func (j *jwtUtil) NewRefreshToken(email, role string, userId uint) (string, *code.ResultCode) {
	tokenId := "refresh-" + uuid.NewString()
	return j.newToken(role, tokenId, email, config.RefreshTokenExpiration, userId)
}

func (j *jwtUtil) NewAccessToken(email, role string, userId uint) (string, *code.ResultCode) {
	tokenId := "access-" + uuid.NewString()
	return j.newToken(role, tokenId, email, config.AccessTokenExpiration, userId)
}

func (j *jwtUtil) newToken(role, jti, email string, expiration time.Duration, userId uint) (string, *code.ResultCode) {
	tokenExpirationTime := time.Now().Add(expiration)

	tokenClaims := &entity.Claims{
		Role:   role,
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			Id:        jti,
			Subject:   email,
			ExpiresAt: tokenExpirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenString, err := token.SignedString(j.jwtKey)
	if err != nil {
		return "", code.InternalServerError.SetMessage("Could not generate token")
	}

	dbToken := &entity.Token{
		JTI:       jti,
		Subject:   email,
		IsRevoked: false,
		ExpiresAt: tokenExpirationTime.Unix(),
	}

	if err := j.tokenRepo.Create(dbToken); err != nil {
		return "", code.InternalServerError.SetMessage("Could not save token")
	}

	return tokenString, nil
}

func (j *jwtUtil) ParseToken(tokenString string) (*entity.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &entity.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return j.jwtKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*entity.Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
