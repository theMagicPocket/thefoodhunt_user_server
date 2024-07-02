package tokenverification

import (
	"log"
	"net/http"

	// "github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}
		userId , _, err := ValidateFirebaseToken(token)
		if err != nil {
			log.Println(err)
			
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		log.Println("token verification")
		log.Println(userId)
		c.Set("userId", userId)

		// Token is valid, proceed with the next handler
		c.Next()
	}
}