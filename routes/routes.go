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
	incomingRoutes.POST("/users/forgotpassword", controllers.ForgotPassword())
	incomingRoutes.PATCH("/users/:resetToken", controllers.ResetPassword())
}

func ProductRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/product/addproduct", controllers.AddProduct())
	incomingRoutes.PUT("/product/editProduct", controllers.EditProduct())
	incomingRoutes.DELETE("/product/removeProduct", controllers.RemoveProduct())
	incomingRoutes.GET("/product", controllers.AllProduct())
	incomingRoutes.GET("/product/search", controllers.SearchProductByQuery())
}
