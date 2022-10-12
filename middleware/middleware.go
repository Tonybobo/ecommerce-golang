package middleware

import (
	token "ecommerce-golang/tokens"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
