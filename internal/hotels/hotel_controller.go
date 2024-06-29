package hotel

import (
	"log"
	"net/http"
	"time"

	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"github.com/gin-gonic/gin"
)

type HotelController struct {
    HotelService HotelService
}

func New(hotelService HotelService) HotelController {
    return HotelController{
        HotelService: hotelService,
    }
}

func (hc *HotelController) CreateHotel(ctx *gin.Context) {
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
    err := hc.HotelService.CreateHotel(&hotel)
    if err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "success", "status_code": http.StatusOK, "hotel_id": hotel.HotelID})
}

func (hc *HotelController) GetHotel(ctx *gin.Context) {
    hotelId := ctx.Param("id")
    hotel, err := hc.HotelService.GetHotel(&hotelId)
    if err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"statuscode": http.StatusOK, "data": hotel})
}

func (hc *HotelController) GetAllHotels(ctx *gin.Context) {
    hotels, err := hc.HotelService.GetAllHotels()
    if err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"statuscode": http.StatusOK, "data": hotels})
}

func (hc *HotelController) UpdateHotel(ctx *gin.Context) {
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
    err := hc.HotelService.UpdateHotel(&hotel)
    if err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "success", "status": http.StatusOK})
}

func (hc *HotelController) DeleteHotel(ctx *gin.Context) {
    hotelId := ctx.Param("id")
    err := hc.HotelService.DeleteHotel(&hotelId)
    if err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "Hotel deleted successfully"})
}

func (hc *HotelController) GetNearbyHotels(c *gin.Context) {
    var userCoordinates struct {
        Lat  float64 `json:"latitude"`
        Long float64 `json:"longitude"`
        Dist float64 `json:"distance"`
    }

    if err := c.ShouldBindJSON(&userCoordinates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    log.Println(userCoordinates.Dist)
    if userCoordinates.Dist == 0{
        c.JSON(http.StatusInternalServerError, gin.H{"error": "please give distance"})
        return
    }
    hotels, err := hc.HotelService.GetNearbyHotels(userCoordinates.Lat, userCoordinates.Long, userCoordinates.Dist)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, hotels)
}


func (hc *HotelController) RegisterHotelRoutes(rg *gin.RouterGroup) {
    hotelRoute := rg.Group("/hotel")
    hotelRoute.POST("/create", hc.CreateHotel)
    hotelRoute.GET("/get/:id", hc.GetHotel)
    hotelRoute.GET("/getall", hc.GetAllHotels)
    hotelRoute.PATCH("/update/:id", hc.UpdateHotel)
    hotelRoute.DELETE("/delete/:id", hc.DeleteHotel)
    hotelRoute.POST("/nearby",hc.GetNearbyHotels)
}
