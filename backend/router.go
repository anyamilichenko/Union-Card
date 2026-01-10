package backend

import (
	"bilet/backend/handler"
	"bilet/backend/middleware"
	"bilet/backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func RegisterTemplates(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "main_page.html", gin.H{
			"Title":   "Авторизация",
			"IsAdmin": false,
		})
	})

	r.GET("/admin_main_menu", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_main_menu.html", gin.H{
			"Title":   "Главная страница администратора",
			"IsAdmin": true,
		})
	})

	r.GET("/add_member", func(c *gin.Context) {
		c.HTML(http.StatusOK, "add_member.html", gin.H{
			"Title":   "Добавить члена профсоюза",
			"IsAdmin": true,
		})
	})

	r.GET("/personal_account", func(c *gin.Context) {
		c.HTML(http.StatusOK, "personal_account.html", gin.H{
			"Title":   "Личный кабинет",
			"IsAdmin": false,
		})
	})

	r.GET("/admin_prof_bilet", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_prof_bilet.html", gin.H{
			"Title":   "Профсоюзный билет (админ)",
			"IsAdmin": true,
		})
	})

	r.GET("/members_list", func(c *gin.Context) {
		c.HTML(http.StatusOK, "members_list.html", gin.H{
			"Title":   "Список членов профсоюза",
			"IsAdmin": true,
		})
	})

	r.GET("/userinfo", func(c *gin.Context) {
		c.HTML(http.StatusOK, "prof_bilet.html", gin.H{
			"Title":   "Профсоюзный билет",
			"IsAdmin": false,
		})
	})
}

func SetupRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	authService service.AuthService,
) *gin.Engine {
	r := gin.Default()

	// Способ 1: Загружаем все HTML файлы из всех поддиректорий
	r.LoadHTMLGlob("frontend/templates/**/*.html")

	// Способ 2: Загружаем шаблоны из каждой поддиректории отдельно
	// r.LoadHTMLGlob("frontend/templates/*.html")
	// r.LoadHTMLGlob("frontend/templates/admin/*.html")
	// r.LoadHTMLGlob("frontend/templates/user/*.html")
	// r.LoadHTMLGlob("frontend/templates/components/*.html")

	r.Static("/static", "frontend/static")

	r.Use(CORSMiddleware())

	// Регистрируем шаблоны
	RegisterTemplates(r)

	// Публичные маршруты API
	r.POST("/api/auth/login", authHandler.Login)
	r.POST("/api/auth/reset_password", authHandler.ResetPassword)
	r.POST("/api/auth/create_tokens", authHandler.CreateTokens)

	// Защищенные маршруты API
	authRoutes := r.Group("/api")
	authRoutes.Use(middleware.AuthMiddleware(authService))
	{
		// Добавляем API маршрут для admin_main_page
		authRoutes.GET("/admin_main_page", userHandler.GetUserInfo)

		authRoutes.GET("/admin_prof_bilet", userHandler.GetUserInfo)
		authRoutes.POST("/auth/logout", authHandler.Logout)
		authRoutes.GET("/members_list", userHandler.GetAllUsers)
		authRoutes.GET("/personal_account", userHandler.GetUserInfo)
		authRoutes.GET("/userinfo", userHandler.GetUserInfo)
	}

	// Админские маршруты API
	adminRoutes := r.Group("/api/admin")
	adminRoutes.Use(middleware.AuthMiddleware(authService), middleware.AdminMiddleware())
	{
		adminRoutes.GET("/members_list", userHandler.GetAllUsers)
		adminRoutes.POST("/add_member", userHandler.AddMember)
		adminRoutes.DELETE("/delete_member", userHandler.DeleteUser)
	}

	return r
}
