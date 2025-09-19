package utils

import (
	"bilet/backend/code"
	"bilet/backend/models"
	"github.com/gin-gonic/gin"
)

func GetClaimsByCookieToken(c *gin.Context) (*models.Claims, *code.ResultCode) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return nil, &code.Unauthorized
	}
	return GetClaimsByToken(token)
}

func GetClaimsByToken(token string) (*models.Claims, *code.ResultCode) {
	claims, err := ParseToken(token)
	if err != nil {
		return nil, &code.Unauthorized
	}
	return claims, nil
}
