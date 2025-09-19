package main

import (
	"bilet/backend"
	"bilet/backend/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	models.InitDB()

	r.LoadHTMLGlob("frontend/templates/**/*")

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(backend.CORSMiddleware())

	r.ForwardedByClientIP = true

	backend.RegisterTemplates(r)
	backend.RegisterHandlers(r)

	r.Run(":1460")
}
