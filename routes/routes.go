package routes

import (
	"ecommerce-golang/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/users/signup", controllers.SignUp())
	incomingRoutes.POST("/users/login", controllers.Login())
	incomingRoutes.POST("/users/logout", controllers.Logout())
	incomingRoutes.GET("/users/user", controllers.GetUser())
	incomingRoutes.PUT("/users/edit", controllers.EditUser())
	// Add function to reset password with email ?????
}

func ProductRoutes(incomingRoutes *gin.Engine) {
	//add product with user id
	incomingRoutes.POST("/product/addproduct", controllers.AddProduct())
	//edit product with user id
	incomingRoutes.PUT("/product/editProduct", controllers.EditProduct())
	//remove product with user id
	incomingRoutes.DELETE("/product/removeProduct", controllers.RemoveProduct())
	incomingRoutes.GET("/product", controllers.AllProduct())
	incomingRoutes.GET("/product/search", controllers.SearchProductByQuery())
}
