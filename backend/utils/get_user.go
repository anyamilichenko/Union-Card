package utils

import (
	"bilet/backend/code"
	"bilet/backend/models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
)

func GetUserByEmail(email string) (*models.Accounts, *code.ResultCode) {
	var existingUser models.Accounts
	fmt.Println(email) //admin@example.com
	dbError := models.DB.Where("Email = ?", email).First(&existingUser).Error
	if dbError != nil {
		log.Printf("Error fetching user by email: %v", dbError)
	}
	if errors.Is(dbError, gorm.ErrRecordNotFound) {
		return nil, &code.UserDoesNotExist
	}
	fmt.Println(existingUser)
	return &existingUser, nil
}

func GetEmailByCookieToken(c *gin.Context) (string, *code.ResultCode) {
	claims, errCode := GetClaimsByCookieToken(c)
	if errCode != nil {
		return "", errCode
	}
	return claims.Subject, nil
}
