package middleware

import (
	"bilet/backend/code"
	"bilet/backend/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var token string

		// 1. Пробуем получить токен из заголовка Authorization
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("Auth header: %s, Path: %s\n", authHeader, c.Request.URL.Path)

		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
				fmt.Printf("Token extracted from header: %s...\n", token[:10])
			}
		}

		// 2. Если нет в заголовке, пробуем из куки
		if token == "" {
			cookieToken, err := c.Cookie("access_token")
			if err == nil {
				token = cookieToken
				fmt.Printf("Token extracted from cookie\n")
			}
		}

		if token == "" {
			fmt.Printf("No token found for path: %s\n", c.Request.URL.Path)
			// Для API запросов возвращаем JSON ошибку
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(http.StatusUnauthorized, code.Unauthorized.SetMessage("Authentication required"))
				c.Abort()
				return
			}
			// Для HTML страниц - редирект на главную
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		claims, errCode := authService.ValidateToken(token)
		if errCode != nil {
			fmt.Printf("Token validation failed: %v\n", errCode)
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(http.StatusUnauthorized, errCode)
				c.Abort()
				return
			}
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}

		fmt.Printf("Token validated for user: %s, role: %s\n", claims.Subject, claims.Role)

		c.Set("claims", claims)
		c.Set("userId", claims.UserId)
		c.Set("userEmail", claims.Subject)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists || role != "admin" {
			if strings.HasPrefix(c.Request.URL.Path, "/api/") {
				c.JSON(http.StatusForbidden, code.Forbidden.SetMessage("Admin access required"))
				c.Abort()
				return
			}
			c.Redirect(http.StatusFound, "/")
			c.Abort()
			return
		}
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	// 1. Попробовать получить из куки
	token, err := c.Cookie("access_token")
	if err == nil && token != "" {
		return token
	}

	// 2. Попробовать получить из заголовка Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	return ""
}
