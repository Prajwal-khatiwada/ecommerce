// package middleware

// import (
// 	"net/http"

// 	token "go-api/tokens"

// 	"github.com/gin-gonic/gin"
// )

// func Authentication() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ClientToken := c.Request.Header.Get("token")
// 		if ClientToken == "" {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization Header Provided"})
// 			c.Abort()
// 			return
// 		}

// 		claims, err := token.ValidateToken(ClientToken)
// 		if err != "" {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
// 			c.Abort()
// 			return
// 		}

// 		c.Set("email", claims.Email)
// 		c.Set("uid", claims.Uid)
// 		c.Next()
// 	}
// }

package middleware

import (
	"net/http"

	token "go-api/tokens"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("token")
		if ClientToken == "" {
			// No token provided, assume the user is a guest
			c.Set("isGuest", true)
			c.Next()
			return
		}

		claims, err := token.ValidateToken(ClientToken)
		if err != "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Token"})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Set("isGuest", false) // Not a guest
		c.Next()
	}
}
