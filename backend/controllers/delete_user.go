package controllers

import (
	"bilet/backend/code"
	"bilet/backend/models"
	"bilet/backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DeleteUser(c *gin.Context) {
	claims, resCode := utils.GetClaimsByCookieToken(c)
	if resCode != nil {
		c.JSON(http.StatusBadRequest, nil)
	}
	var existingUser models.Accounts
	err1 := models.DB.Where("id = ?", claims.UserId).First(&existingUser).Error
	if err1 != nil {
		c.JSON(http.StatusBadRequest, nil)
	}

	var requestData struct {
		ID uint `json:"id"`
	}
	fmt.Println("это id", requestData.ID)
	// Привязка JSON-тела запроса к структуре
	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	// Поиск пользователя по переданному ID
	var deleteUser models.Accounts
	if err := models.DB.Where("id = ?", requestData.ID).First(&deleteUser).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	// Удаление пользователя
	if err := models.DB.Delete(&deleteUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code.Success.Code,
		"message": "Пользователь успешно удалён",
	})
}
