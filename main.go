package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleHome(c *gin.Context) {
	c.JSON(http.StatusOK ,"Hello Go!")
}


func main() {
	fmt.Println("HI")
	r := gin.Default();

	r.GET("/", handleHome)

	r.Run("localhost:8080")
}