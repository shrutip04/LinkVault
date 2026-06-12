package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shrutip04/linkvault/database"
	"github.com/shrutip04/linkvault/routes"
)

func main() {
	fmt.Println("LinkVault starting...")

	database.InitDB()

	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run(":8080")
}