package main

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/shrutip04/linkvault/database"
	"github.com/shrutip04/linkvault/routes"
)

func main() {
	fmt.Println("LinkVault starting...")

	database.InitDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://your-frontend-domain.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	routes.SetupRoutes(r)

	r.Run(":8080")
}