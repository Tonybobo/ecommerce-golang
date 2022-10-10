package controllers

import (
	"context"
	"ecommerce-golang/database"
	"ecommerce-golang/models"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryId := c.Query("userID")
		if userQueryId == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}

		c.IndentedJSON(200, "Successfully added to cart")
	}
}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryId := c.Query("userID")
		if userQueryId == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, productID, userQueryId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "Successfully removed item from cart")

	}
}

func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}

		userId, _ := primitive.ObjectIDFromHex(user_id)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		var cart models.User

		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userId}}).Decode(&cart)

		if err != nil {
			log.Println(err)
		}

		filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: userId}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}

		cursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter, unwind, group})

		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Internal Server Error")
			return
		}

		var listing []bson.M

		if err := cursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.IndentedJSON(500, "Internal Server Error")
			return
		}

		for _, json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, cart.UserCart)
		}

		ctx.Done()

	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryId := c.Query("id")

		if userQueryId == "" {
			log.Panicln("User id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User Id is empty"))
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "Successfully buy from cart")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}
		userQueryId := c.Query("userID")
		if userQueryId == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		productID, err := primitive.ObjectIDFromHex(productQueryId)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()

		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productID, userQueryId)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "Successfully placed the order")
	}
}
