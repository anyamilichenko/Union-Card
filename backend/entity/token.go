package entity

import "github.com/golang-jwt/jwt"

type Token struct {
	JTI       string `gorm:"primaryKey"`
	Subject   string
	IsRevoked bool `gorm:"default:false"`
	ExpiresAt int64
}

type Claims struct {
	Role   string
	UserId uint
	jwt.StandardClaims
}
