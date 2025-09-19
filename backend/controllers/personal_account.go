package controllers

import (
	"bilet/backend/code"
	"bilet/backend/jsonr"
	"bilet/backend/models"
	"bilet/backend/utils"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func PersonalAccount(c *gin.Context) {
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
		"code": code.Success.Code,
	})
}
func AddMember(c *gin.Context) {
	claims, resCode := utils.GetClaimsByCookieToken(c)
	if resCode != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid token"})
		return
	}

	var existingUser models.Accounts
	err1 := models.DB.Where("id = ?", claims.UserId).First(&existingUser).Error
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	var addMember jsonr.AddMember

	addMember.LastName = c.PostForm("lastName")
	addMember.FirstName = c.PostForm("firstName")
	addMember.MiddleName = c.PostForm("middleName")
	addMember.DateBirth = c.PostForm("dateBirth")
	addMember.PhoneNumber = c.PostForm("phoneNumber")
	addMember.Email = c.PostForm("email")
	addMember.MembershipStatus = c.PostForm("membershipStatus")
	addMember.Role = c.PostForm("role")
	addMember.Password = c.PostForm("password")

	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error: " + err.Error()})
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error opening photo file", "error": err.Error()})
		return
	}
	defer fileContent.Close()

	photo, err := io.ReadAll(fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error reading photo data", "error": err.Error()})
		return
	}

	account := &models.Accounts{
		LastName:         addMember.LastName,
		FirstName:        addMember.FirstName,
		MiddleName:       addMember.MiddleName,
		DateBirth:        addMember.DateBirth,
		PhoneNumber:      addMember.PhoneNumber,
		Email:            addMember.Email,
		Photo:            photo,
		MembershipStatus: addMember.MembershipStatus,
		Role:             addMember.Role,
	}

	if err := models.DB.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error saving user to database", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code.Success.Code,
		"message": "Created user",
		"user":    account,
	})
}

func GetUserInfo(c *gin.Context) {
	claims, resCode := utils.GetClaimsByCookieToken(c)
	if resCode != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid token"})
		return
	}

	var existingUser models.Accounts
	err := models.DB.Where("id = ?", claims.UserId).First(&existingUser).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	// Конвертируем фото в base64
	var photoBase64 string
	if len(existingUser.Photo) > 0 {
		photoBase64 = base64.StdEncoding.EncodeToString(existingUser.Photo)
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"last_name":         existingUser.LastName,
			"first_name":        existingUser.FirstName,
			"middle_name":       existingUser.MiddleName,
			"date_birth":        existingUser.DateBirth,
			"phone_number":      existingUser.PhoneNumber,
			"email":             existingUser.Email,
			"membership_status": existingUser.MembershipStatus,
			"role":              existingUser.Role,
			"photo":             photoBase64,
		},
		"code": code.Success.Code,
	})
}

func GetAllUsers(c *gin.Context) {

	_, resCode := utils.GetClaimsByCookieToken(c)
	if resCode != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
		return
	}

	var users []models.Accounts
	err := models.DB.Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"code":  code.Success.Code,
	})
}
