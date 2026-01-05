package handler

import (
	"bilet/backend/code"
	"bilet/backend/service"
	"bilet/backend/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AuthHandler struct {
	authService service.AuthService
	jwtUtil     utils.JWTUtil
}

func NewAuthHandler(authService service.AuthService, jwtUtil utils.JWTUtil) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		jwtUtil:     jwtUtil,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, code.BadRequest.SetMessage(err.Error()))
		return
	}

	user, errCode := h.authService.Login(loginData.Email, loginData.Password)
	if errCode != nil {
		c.JSON(http.StatusUnauthorized, errCode)
		return
	}

	refreshToken, errCode := h.jwtUtil.NewRefreshToken(user.Email, user.Role, user.ID)
	if errCode != nil {
		c.JSON(errCode.Code, errCode)
		return
	}

	accessToken, errCode := h.jwtUtil.NewAccessToken(user.Email, user.Role, user.ID)
	if errCode != nil {
		c.JSON(errCode.Code, errCode)
		return
	}

	// Возвращаем токены только в JSON (без установки кук)
	// Фронтенд будет сохранять их в localStorage
	c.JSON(http.StatusOK, gin.H{
		"code":          code.Success.Code,
		"message":       "User logged in",
		"user":          user,
		"refresh_token": refreshToken,
		"access_token":  accessToken,
	})
}
func (h *AuthHandler) CreateTokens(c *gin.Context) {
	// Пробуем получить refresh token из JSON тела или заголовка
	var requestData struct {
		RefreshToken string `json:"refresh_token"`
	}

	// Сначала пробуем из JSON
	if err := c.ShouldBindJSON(&requestData); err != nil {
		// Если нет в JSON, пробуем из заголовка
		token := extractTokenFromHeader(c)
		if token == "" {
			c.JSON(http.StatusUnauthorized, code.Unauthorized.SetMessage("Refresh token required"))
			return
		}
		requestData.RefreshToken = token
	}

	refreshToken, accessToken, errCode := h.authService.CreateTokens(requestData.RefreshToken)
	if errCode != nil {
		c.JSON(errCode.Code, errCode)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":          code.Success.Code,
		"refresh_token": refreshToken,
		"access_token":  accessToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Пробуем получить токен из заголовка Authorization
	token := extractTokenFromHeader(c)

	if token != "" {
		// Отзываем токен на сервере
		h.authService.Logout(token)
	}

	// Также очищаем куки на всякий случай
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.SetCookie("access_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"code":    code.Success.Code,
		"message": "Вы успешно вышли из системы",
	})
}

func extractTokenFromHeader(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 && parts[0] == "Bearer" {
		return parts[1]
	}

	return ""
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var requestData struct {
		Email string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	errCode := h.authService.ResetPassword(requestData.Email)
	if errCode != nil {
		c.JSON(errCode.Code, errCode)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "New password sent to email",
	})
}
