package hotels

import (
	"net/http"
	"time"

	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"github.com/deVamshi/golang_food_delivery_api/pkg/log"
	"github.com/gin-gonic/gin"
)

type resource struct {
	service Service
	logger  log.AppLogger
}

func RegisterHandlers(rg *gin.RouterGroup, service Service, logger log.AppLogger) {

	res := resource{service: service, logger: logger}

	hotelRoute := rg.Group("/hotel")
	hotelRoute.POST("/create", res.CreateHotel)
	hotelRoute.GET("/get/:id", res.GetHotel)
	hotelRoute.GET("/getall", res.GetAllHotels)
	hotelRoute.PATCH("/update/:id", res.UpdateHotel)
	hotelRoute.DELETE("/delete/:id", res.DeleteHotel)
	hotelRoute.POST("/nearby", res.GetNearbyHotels)
}

func (r resource) CreateHotel(ctx *gin.Context) {
	var hotel entity.Hotel
	hotel.HotelID = entity.GenerateID()
	istLocation, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(istLocation)
	formattedTime := now.Format("01-02-2006 15:04:05")
	hotel.CreatedAt = formattedTime
	hotel.UpdatedAt = formattedTime
	if err := ctx.ShouldBindJSON(&hotel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := r.service.CreateHotel(&hotel)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success", "status_code": http.StatusOK, "hotel_id": hotel.HotelID})
}

func (r resource) GetHotel(ctx *gin.Context) {

	hotelId := ctx.Param("id")
	hotel, err := r.service.GetHotel(hotelId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"statuscode": http.StatusOK, "data": hotel})
}

func (r resource) GetAllHotels(ctx *gin.Context) {
	hotels, err := r.service.GetAllHotels()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"statuscode": http.StatusOK, "data": hotels})
}

func (r resource) UpdateHotel(ctx *gin.Context) {
	var hotel entity.Hotel
	hotel.HotelID = ctx.Param("id")
	if err := ctx.ShouldBindJSON(&hotel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	istLocation, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(istLocation)
	formattedTime := now.Format("01-02-2006 15:04:05")
	hotel.UpdatedAt = formattedTime
	err := r.service.UpdateHotel(&hotel)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success", "status": http.StatusOK})
}

func (r resource) DeleteHotel(ctx *gin.Context) {
	hotelId := ctx.Param("id")
	err := r.service.DeleteHotel(hotelId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Hotel deleted successfully"})
}

func (r resource) GetNearbyHotels(c *gin.Context) {
	var userCoordinates struct {
		Lat  float64 `json:"latitude"`
		Long float64 `json:"longitude"`
		Dist float64 `json:"distance"`
	}

	if err := c.ShouldBindJSON(&userCoordinates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	r.logger.Info(userCoordinates.Dist)
	if userCoordinates.Dist == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "please give distance"})
		return
	}
	hotels, err := r.service.GetNearbyHotels(userCoordinates.Lat, userCoordinates.Long, userCoordinates.Dist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, hotels)
}
