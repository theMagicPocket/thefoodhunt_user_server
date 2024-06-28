package matrixapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	
)

func distance(c *gin.Context) {
    var req entity.DistanceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	origin := fmt.Sprintf("%f,%f", req.Lat1, req.Lng1)
	destination := fmt.Sprintf("%f,%f", req.Lat2, req.Lng2)
	APP_ENV, err := godotenv.Read("../.env")
	if err != nil {
		log.Fatal(err)
		log.Fatal("Error loading .env file")
	}
	apiKey := APP_ENV["MATRIX_KEY"]

	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/distancematrix/json?origins=%s&destinations=%s&key=%s", origin, destination, apiKey)
	client := resty.New()
	resp, err := client.R().Get(url)
	if err != nil {
		log.Printf("Error fetching data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching data"})
		return
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		log.Printf("Error unmarshaling response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error unmarshaling response"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func RegisterDistanceMatrixRoutes(rg *gin.RouterGroup) {
	distanceroute := rg.Group("/matrix/distance")
	distanceroute.POST("", distance)
}
