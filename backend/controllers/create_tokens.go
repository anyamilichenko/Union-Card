package controllers

import (
	"bilet/backend/code"
	"bilet/backend/models"
	"bilet/backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func CreateTokens(c *gin.Context) {
	token, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, code.Unauthorized.SetMessage(err.Error()))
		return
	}

	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, code.Unauthorized.SetMessage(err.Error()))
		return
	}

	if !strings.HasPrefix(claims.Id, "refresh-") {
		c.JSON(http.StatusUnauthorized, code.Unauthorized.SetMessage("Bad token provided"))
		return
	}

	if time.Now().After(time.Unix(claims.ExpiresAt, 0)) {
		c.JSON(http.StatusUnauthorized, code.Unauthorized.SetMessage("Token has expired"))
		return
	}

	var refreshToken models.Token
	dbError := models.DB.Where("jti = ?", claims.Id).First(&refreshToken).Error
	if dbError != nil {
		c.JSON(http.StatusUnauthorized, code.Unauthorized.SetMessage("Token not found"))
		return
	}

	if refreshToken.IsRevoked {
		errCode := DeleteUserTokens(claims.Subject)
		if errCode != nil {
			c.JSON(errCode.Code, errCode)
			return
		}
		c.JSON(http.StatusUnauthorized, code.Unauthorized.SetMessage("Token has already been used"))
		return
	}

	refreshTokenString, errCode := utils.NewRefreshToken(claims.Subject, claims.Role, claims.UserId)
	if errCode != nil {
		c.JSON(errCode.Code, errCode)
		return
	}

	accessTokenString, errCode := utils.NewAccessToken(claims.Subject, claims.Role, claims.UserId)
	if errCode != nil {
		c.JSON(errCode.Code, errCode)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":          code.Success.Code,
		"refresh_token": refreshTokenString,
		"access_token":  accessTokenString,
	})
}

func DeleteUserTokens(Email string) *code.ResultCode {
	dbError := models.DB.Where("subject = ?", Email).Delete(&models.Token{}).Error
	if dbError != nil {
		return code.InternalServerError.SetMessage("Failed to delete user tokens")
	}
	return nil
}

func RevokeDeviceTokens(deviceId string) *code.ResultCode {
	dbError := models.DB.Model(&models.Token{}).Where("device_id = ?", deviceId).Update("is_revoked", true).Error
	if dbError != nil {
		return code.InternalServerError.SetMessage("Failed to revoke device=" + deviceId + " tokens")
	}
	return nil
}
