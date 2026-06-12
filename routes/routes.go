package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shrutip04/linkvault/handlers"
	"github.com/shrutip04/linkvault/middleware"
)

func SetupRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.GET("/qr/:code", handlers.GenerateQR)
	r.GET("/:code", handlers.RedirectURL)

	// Protected routes (JWT required)
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		auth.POST("/shorten", handlers.ShortenURL)
		auth.GET("/links", handlers.GetAllLinks)
		auth.GET("/links/stats", handlers.GetStats)
	}
}
