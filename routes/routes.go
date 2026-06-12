package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shrutip04/linkvault/handlers"
	"github.com/shrutip04/linkvault/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// Public
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.GET("/qr/:code", handlers.GenerateQR)
	r.GET("/:code", handlers.RedirectURL)

	// Protected
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		auth.POST("/shorten", handlers.ShortenURL)
		auth.GET("/links", handlers.GetAllLinks)
		auth.GET("/links/stats", handlers.GetStats)
		auth.DELETE("/links/:id", handlers.DeleteLink)
		auth.PUT("/links/:id", handlers.EditLink)
	}
}
