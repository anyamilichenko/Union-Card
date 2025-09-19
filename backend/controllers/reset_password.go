package controllers

import (
	"bilet/backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func ResetPassword(c *gin.Context) {
	var requestData struct {
		Email string `json:"email"`
	}

	// Получаем email из запроса
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Неверный запрос"})
		return
	}

	// Проверяем, существует ли пользователь с таким email
	var user models.Accounts
	if err := models.DB.Where("email = ?", requestData.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Пользователь не найден"})
		return
	} else {
		_ = models.DB.Where("email = ?", requestData.Email).Delete(&models.Accounts{}).Error
	}

	// Генерация нового пароля
	newPassword := GeneratePassword(10)

	// Хеширование нового пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при хешировании пароля"})
		return
	}

	// Обновляем пароль в базе данных
	user.HashedPassword = string(hashedPassword)
	if err := models.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при сохранении пароля"})
		return
	}

	// Отправляем новый пароль на email
	if err := SendPasswordEmail(user.Email, newPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при отправке письма"})
		return
	}

	// Возвращаем успешный ответ
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Новый пароль отправлен на почту"})
}
