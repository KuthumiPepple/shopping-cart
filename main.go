package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kuthumipepple/ecommerce-platform/controllers"
	"github.com/kuthumipepple/ecommerce-platform/database"
	"github.com/kuthumipepple/ecommerce-platform/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	controllers.NewApplication(
		database.OpenCollection("products"),
		database.OpenCollection("users"),
	)

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)
	log.Fatal(router.Run(":" + port))
}
