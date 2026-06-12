package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shrutip04/linkvault/handlers"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/shorten", handlers.ShortenURL)
	r.GET("/links", handlers.GetAllLinks)
	r.GET("/:code", handlers.RedirectURL)
}