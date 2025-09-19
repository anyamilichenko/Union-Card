package utils

import (
	"bilet/backend/code"
	"github.com/gin-gonic/gin"
)

func GetRoleByToken(c *gin.Context) (string, *code.ResultCode) {
	claims, errCode := GetClaimsByCookieToken(c)
	if errCode != nil {
		return "", errCode
	}
	return claims.Role, nil
}
