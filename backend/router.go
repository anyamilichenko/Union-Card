package backend

import (
	"bilet/backend/controllers"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterTemplates(r *gin.Engine) {
	// Главная страница (с формой авторизации)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "main_page.html", gin.H{})
	})

	r.GET("/admin_main_menu", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_main_menu.html", gin.H{})
	})

	r.GET("/add_member", func(c *gin.Context) {
		c.HTML(http.StatusOK, "add_member.html", gin.H{})
	})

	r.GET("/personal_account", func(c *gin.Context) {
		c.HTML(http.StatusOK, "personal_account.html", gin.H{})
	})

	r.GET("/admin_prof_bilet", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin_prof_bilet.html", gin.H{})
	})

	r.GET("/members_list", func(c *gin.Context) {
		c.HTML(http.StatusOK, "member_list.html", gin.H{})
	})

	r.GET("/userinfo", func(c *gin.Context) {
		c.HTML(http.StatusOK, "prof_bilet.html", gin.H{})
	})
}

func RegisterHandlers(r *gin.Engine) {
	api := r.Group("api")
	api.GET("/personal_account", controllers.PersonalAccount)
	api.GET("/userinfo", controllers.GetUserInfo)
	api.GET("/admin_prof_bilet", controllers.GetUserInfo)
	api.POST("/add_member", controllers.AddMember)
	api.GET("/admin_main_page", controllers.AdminMain)
	api.GET("/members_list", controllers.GetAllUsers)
	api.DELETE("/delete_member", controllers.DeleteUser)

	auth := api.Group("auth")
	auth.POST("/login", controllers.Login)
	auth.POST("/reset_password", controllers.ResetPassword)
	auth.POST("/create_tokens", controllers.CreateTokens)
	auth.POST("/logout", controllers.Logout)
}
