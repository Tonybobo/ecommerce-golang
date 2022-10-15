package middleware

import (
	token "ecommerce-golang/tokens"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))

	valid := true
	msg := ""

	if err != nil {
		msg = "Email or Password is incorrect"
		valid = false
	}

	return valid, msg
}

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {

		ClientToken := c.Request.Header.Get("token")

		if ClientToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"errors": "No Authorized"})
			c.Abort()
			return
		}

		claim, err := token.ValidateToken(ClientToken)
		if err != "" {
			c.JSON(http.StatusBadRequest, gin.H{"errors": "Invalid Token"})
			c.Abort()
			return
		}

		c.Set("email", claim.Email)
		c.Set("uid", claim.Uid)
		c.Next()
	}
}
