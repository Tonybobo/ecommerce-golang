package controllers

import "github.com/gin-gonic/gin"

func HashPassword(password string) string {
	return ""
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	return true, ""
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func ProductViewerAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
