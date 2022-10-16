package controllers

import (
	"context"
	"ecommerce-golang/models"
	"ecommerce-golang/tokens"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func AddProduct() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var product models.Product
		var user models.User

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		}

		if err := UserCollection.FindOne(ctx, bson.M{"user_id": product.User_id}).Decode(&user); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"errors": "No user with this user_id"})
			return
		}
		defer cancel()

		product.Product_ID = primitive.NewObjectID()

		_, err := ProductCollection.InsertOne(ctx, product)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "error occured while adding product"})
			return
		}

		defer cancel()
		c.JSON(http.StatusOK, "Product Added")
	}
}

func EditProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		//extract token to get user id
		token := c.Request.Header.Get("token")
		claim, msg := tokens.ValidateToken(token)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		if msg != "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Expired Token /Invalid Token"})
		}

		var product models.Product

		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		}

		query := bson.D{{Key: "_id", Value: product.Product_ID}, {Key: "user_id", Value: claim.Uid}}
		fmt.Println(query)
		update := bson.D{{Key: "$set", Value: product}}

		res := ProductCollection.FindOneAndUpdate(
			ctx,
			query,
			update,
			options.FindOneAndUpdate().SetReturnDocument(1))

		var updatedProdcut models.Product
		if err := res.Decode(&updatedProdcut); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		defer cancel()
		c.JSON(http.StatusOK, updatedProdcut)
	}
}

func RemoveProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productlist []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = cursor.All(ctx, &productlist)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(ctx)

		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.JSON(http.StatusNotFound, gin.H{"error": "No Product Found"})
			return
		}

		defer cancel()

		c.JSON(200, productlist)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchProducts []models.Product
		queryParams := c.Query("name")

		if queryParams == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusBadRequest, gin.H{"errors": "Query is empty"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParams}})

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "Something went wrong when fetching the data"})
			return
		}

		err = cursor.All(ctx, &searchProducts)
		if err != nil {
			log.Println(err)
			c.JSON(400, "No Product Matched")
			return
		}

		defer cursor.Close(ctx)

		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.JSON(http.StatusNotFound, "No Product Matched")
			return
		}

		defer cancel()

		c.JSON(http.StatusAccepted, searchProducts)
	}
}
