package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kuthumipepple/shopping-cart/controllers"
	"github.com/kuthumipepple/shopping-cart/database"
	"github.com/kuthumipepple/shopping-cart/middleware"
	"github.com/kuthumipepple/shopping-cart/routes"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app := controllers.NewApplication(
		database.OpenCollection("products"),
		database.OpenCollection("users"),
	)

	router := gin.New()
	router.Use(gin.Logger())

	routes.UserRoutes(router)

	router.Use(middleware.Authenticate())

	router.PATCH("/addtocart", app.AddToCart())
	router.PATCH("/removeitem", app.RemoveItem())
	
	log.Fatal(router.Run(":" + port))
}
