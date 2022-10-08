package controllers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func NewApplication(prodCollection , userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection : prodCollection,
		userCollection : userCollection
	}
}