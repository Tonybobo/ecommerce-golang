package controllers

import (
	"context"
	"ecommerce-golang/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			log.Println("User id is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid User ID"})
			c.Abort()
			return
		}

		address := make([]models.Address, 0)
		userId, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		filter := bson.D{primitive.E{Key: "_id", Value: userId}}

		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "address", Value: address}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.IndentedJSON(404, "Something Went Wrong While updating address")
			return
		}

		defer cancel()

		ctx.Done()
		c.IndentedJSON(200, "Successfully Updated ")
	}

}
