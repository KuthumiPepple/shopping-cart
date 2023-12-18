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

	router.POST("/addtocart", app.AddToCart())
	router.DELETE("/removeitem", app.RemoveItem())
	router.POST("/instantbuy", app.InstantBuy())
	router.POST("/cartcheckout", app.BuyFromCart())

	log.Fatal(router.Run(":" + port))
}
