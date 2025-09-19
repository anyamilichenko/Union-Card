package controllers

import (
	"bilet/backend/code"
	"bilet/backend/jsonr"
	"bilet/backend/utils"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func Login(c *gin.Context) {
	var loginJson jsonr.UserLoginJson
	rawData, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawData))

	err := c.ShouldBindJSON(&loginJson)

	if err != nil {
		c.JSON(http.StatusBadRequest, code.BadRequest.SetMessage(err.Error()))
		fmt.Println()
		return
	}

	user, errCode := utils.GetUserByEmail(loginJson.Email)
	if errCode != nil {
		c.JSON(http.StatusUnauthorized, errCode)
		return
	}
	if user.HashedPassword == "" {
		c.JSON(http.StatusUnauthorized, code.UserPasswordIsNotSet)
		return
	}

	if !utils.IsPasswordCorrect(loginJson.Password, user.HashedPassword) {
		c.JSON(http.StatusUnauthorized, code.InvalidPassword)
		return
	}

	refreshTokenString, errCode := utils.NewRefreshToken(user.Email, user.Role, user.Id)
	if errCode != nil {
		c.JSON(errCode.Code, errCode)
		return
	}

	accessTokenString, errCode := utils.NewAccessToken(user.Email, user.Role, user.Id)
	if errCode != nil {
		c.JSON(errCode.Code, errCode)
		return
	}

	c.SetCookie("refresh_token", refreshTokenString, 0, "", "", false, false)
	c.SetCookie("access_token", accessTokenString, 0, "", "", false, false)

	c.JSON(http.StatusOK, gin.H{
		"code":          code.Success.Code,
		"message":       "User logged in",
		"user":          user,
		"refresh_token": refreshTokenString,
		"access_token":  accessTokenString,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"code":    code.Success.Code,
		"message": "Вы успешно вышли из системы",
	})
}
