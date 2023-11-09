package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kuthumipepple/shopping-cart/controllers"
)

func UserRoutes(routes *gin.Engine) {
	routes.POST("/users/signup", controllers.Signup())
	routes.POST("/users/login", controllers.Login())
	routes.POST("/admin/products", controllers.AddProductAdmin())
	routes.GET("/users/products", controllers.GetProducts())
	routes.GET("/users/search", controllers.SearchProductByQuery())
}
