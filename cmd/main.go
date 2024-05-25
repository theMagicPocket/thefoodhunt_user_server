package main

import (
	"github.com/deVamshi/golang_food_delivery_api/internal/hotel"
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	v1 := r.Group("/v1")
	{
		hotel.RegisterHandlers(v1, hotel.NewService(hotel.NewRepository("db")))
	}

	r.Run("localhost:8080")
}
