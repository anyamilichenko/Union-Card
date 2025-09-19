package models

import "github.com/golang-jwt/jwt"

type Claims struct {
	Role   string
	UserId int
	jwt.StandardClaims
}
