package routes

import (
	"ecommerce-golang/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	//logout
	//edit information
	// Add function to reset password with email ?????
}

func ProductRoutes(incomingRoutes *gin.Engine) {
	//add product with user id
	incomingRoutes.POST("/product/addproduct", controllers.AddProduct())
	//edit product with user id
	//remove product with user id
	incomingRoutes.GET("/product/productview", controllers.SearchProduct())
	incomingRoutes.GET("/product/search", controllers.SearchProductByQuery())
}
