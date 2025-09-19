package controllers

import (
	"bilet/backend/models"
	"bilet/backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AdminMain(c *gin.Context) {
	claims, resCode := utils.GetClaimsByCookieToken(c)

	if resCode != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Неверный или отсутствующий токен",
		})
		return
	}
	var existingUser models.Accounts
	err := models.DB.Where("id = ?", claims.UserId).First(&existingUser).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Неверный токен",
		})
		return
	}
	fmt.Println(existingUser)

	c.JSON(http.StatusOK, gin.H{
		"user": existingUser,
		"code": 200,
	})
}
