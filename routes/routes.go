package routes

import (
	"go-api/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/admin/addproduct", controllers.AddProductAdmin())
	incomingRoutes.GET("/admin/deleteproduct", controllers.DeleteProductAdmin())
	incomingRoutes.PUT("/admin/editproduct", controllers.EditProductAdmin())
	incomingRoutes.GET("/users/productview", controllers.SearchProduct())
	incomingRoutes.GET("/users/search", controllers.SearchProductByQuery())
	incomingRoutes.GET("/users/search-category", controllers.SearchProductByCategory())
}
