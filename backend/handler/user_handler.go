package handler

import (
	"bilet/backend/code"
	"bilet/backend/entity"
	"bilet/backend/service"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	_ "strconv"
)

type UserHandler struct {
	userService service.UserService
	authService service.AuthService
}

func NewUserHandler(userService service.UserService, authService service.AuthService) *UserHandler {
	return &UserHandler{
		userService: userService,
		authService: authService,
	}
}

func (h *UserHandler) GetUserInfo(c *gin.Context) {
	log.Println("GetUserInfo handler called")

	// ВАЖНО: Используем данные из контекста, которые установил middleware
	claims, exists := c.Get("claims")
	if !exists {
		log.Println("Claims not found in context")
		c.JSON(http.StatusUnauthorized, code.Unauthorized.SetMessage("Authentication required"))
		return
	}

	userClaims := claims.(*entity.Claims)
	log.Printf("Getting user by ID: %d", userClaims.UserId)

	user, errCode := h.userService.GetUserByID(userClaims.UserId)
	if errCode != nil {
		log.Printf("Service error: %v", errCode)
		c.JSON(errCode.Code, errCode)
		return
	}

	// Конвертируем фото в base64
	var photoBase64 string
	if len(user.Photo) > 0 {
		photoBase64 = base64.StdEncoding.EncodeToString(user.Photo)
	}

	log.Printf("Returning user data for: %s", user.Email)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"last_name":         user.LastName,
			"first_name":        user.FirstName,
			"middle_name":       user.MiddleName,
			"date_birth":        user.DateBirth,
			"phone_number":      user.PhoneNumber,
			"email":             user.Email,
			"membership_status": user.MembershipStatus,
			"role":              user.Role,
			"photo":             photoBase64,
		},
		"code": code.Success.Code,
	})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, errCode := h.userService.GetAllUsers()
	if errCode != nil {
		c.JSON(errCode.Code, errCode)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"code":  code.Success.Code,
	})
}

func (h *UserHandler) AddMember(c *gin.Context) {
	log.Println("AddMember called")

	// ВАЖНО: Используем c.PostForm вместо ShouldBind для multipart/form-data
	lastName := c.PostForm("lastName")
	firstName := c.PostForm("firstName")
	middleName := c.PostForm("middleName")
	dateBirth := c.PostForm("dateBirth")
	phoneNumber := c.PostForm("phoneNumber")
	email := c.PostForm("email")
	membershipStatus := c.PostForm("membershipStatus")
	role := c.PostForm("role")

	log.Printf("Received form data: %s %s %s %s %s %s %s %s",
		lastName, firstName, middleName, dateBirth, phoneNumber, email, membershipStatus, role)

	// Валидация обязательных полей
	if lastName == "" || firstName == "" || email == "" || phoneNumber == "" || dateBirth == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Заполните обязательные поля: фамилия, имя, email, телефон, дата рождения",
		})
		return
	}

	// Получаем фото
	var photo []byte
	file, err := c.FormFile("photo")
	if err == nil && file != nil {
		fileContent, err := file.Open()
		if err != nil {
			log.Printf("Error opening file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при открытии файла",
			})
			return
		}
		defer fileContent.Close()

		photo, err = io.ReadAll(fileContent)
		if err != nil {
			log.Printf("Error reading file: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка при чтении файла",
			})
			return
		}
		log.Printf("Photo size: %d bytes", len(photo))
	} else {
		log.Printf("No photo uploaded: %v", err)
	}

	// Создаем аккаунт
	account := &entity.Account{
		LastName:         lastName,
		FirstName:        firstName,
		MiddleName:       middleName,
		DateBirth:        dateBirth,
		PhoneNumber:      phoneNumber,
		Email:            email,
		Photo:            photo,
		MembershipStatus: membershipStatus,
		Role:             role,
	}

	log.Printf("Creating account: %+v", account)

	createdUser, errCode := h.userService.CreateUser(account)
	if errCode != nil {
		log.Printf("Service error: %v", errCode)
		c.JSON(errCode.Code, errCode)
		return
	}

	log.Println("Member created successfully")

	c.JSON(http.StatusOK, gin.H{
		"code":    code.Success.Code,
		"message": "Член профсоюза успешно добавлен",
		"user":    createdUser,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	token, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, code.Unauthorized.SetMessage("Token required"))
		return
	}

	claims, errCode := h.authService.ValidateToken(token)
	if errCode != nil {
		c.JSON(errCode.Code, errCode)
		return
	}

	var requestData struct {
		ID uint `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	errCode = h.userService.DeleteUser(claims.UserId, requestData.ID)
	if errCode != nil {
		c.JSON(errCode.Code, errCode)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code.Success.Code,
		"message": "User deleted successfully",
	})
}
