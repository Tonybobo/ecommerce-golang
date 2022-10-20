package controllers

import (
	"context"
	"ecommerce-golang/database"
	"ecommerce-golang/middleware"
	"ecommerce-golang/models"
	generate "ecommerce-golang/tokens"
	"ecommerce-golang/utils"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	UserCollection    *mongo.Collection = database.UserData(database.Client, "Users")
	ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
	Validate                            = validator.New()
)

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		}
		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone is already in use"})
			return
		}
		password := middleware.HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_ID = user.ID.Hex()
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		_, inserterr := UserCollection.InsertOne(ctx, user)
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}
		defer cancel()
		c.JSON(http.StatusCreated, "Successfully Signed Up!!")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var (
			user      models.User
			foundUser models.User
		)
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password incorrect"})
			return
		}

		PasswordIsValid, msg := middleware.VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !PasswordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		token, refreshToken, err := generate.TokenGenerator(*foundUser.Email, *foundUser.First_Name, *foundUser.Last_Name, foundUser.User_ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		defer cancel()

		generate.UpdateAllTokens(token, refreshToken, foundUser.User_ID)
		c.Header("token", token)
		c.JSON(http.StatusOK, foundUser)
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("token", "")
		c.JSON(http.StatusOK, "Successfully logout")
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		claim, msg := generate.ValidateToken(token)

		if msg != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid/Expired Token"})
			return
		}

		var user models.User
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		userId, err := primitive.ObjectIDFromHex(claim.Uid)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := UserCollection.FindOne(ctx, bson.M{"_id": userId}, options.FindOne())

		if err := result.Decode(&user); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		defer cancel()

		c.JSON(http.StatusOK, user)

	}
}

func EditUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		token := c.Request.Header.Get("token")
		claim, msg := generate.ValidateToken(token)
		if msg != "" {
			c.JSON(http.StatusForbidden, gin.H{"errors": "Invalid/Expired Token"})
			return
		}

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}

		var updatedUser models.User
		userId, _ := primitive.ObjectIDFromHex(claim.Uid)
		user.Token = &token
		user.Refresh_Token = &token
		query := bson.D{{Key: "_id", Value: userId}}
		update := bson.D{{Key: "$set", Value: user}}

		result := UserCollection.FindOneAndUpdate(ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

		defer cancel()
		if err := result.Decode(&updatedUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updatedUser)
	}
}

func ForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.ForgotPasswordInput
		var temp = template.Must(template.ParseGlob("templates/*.html"))
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user models.User
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		query := bson.M{"email": strings.ToLower(input.Email)}
		err := UserCollection.FindOne(ctx, query).Decode(&user)
		defer cancel()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusOK, gin.H{"data": "Please Check your email"})
				return
			}
			c.JSON(http.StatusBadGateway, gin.H{"errors": err.Error()})
			return
		}

		emailData := utils.EmailData{
			URL:       "http://localhost:3000/resetPassword/token=" + *user.Refresh_Token,
			FirstName: *user.First_Name,
			Subject:   "Reset Password",
		}

		err = utils.SendEmail(&user, &emailData, temp, "resetPassword.html")
		if err != nil {
			fmt.Println(err.Error())
		}

		c.JSON(http.StatusOK, gin.H{"data": "Please Check your email"})
	}
}
