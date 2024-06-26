package deliveryfee

import (
	"math"
	"net/http"

	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"github.com/gin-gonic/gin"
)

const earthRadiusKm = 6371
func toRadians(angle float64) float64 {
    return angle * math.Pi / 180
}

func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
    dLat := toRadians(lat2 - lat1)
    dLng := toRadians(lng2 - lng1)

    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Cos(toRadians(lat1))*math.Cos(toRadians(lat2))*
            math.Sin(dLng/2)*math.Sin(dLng/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    return earthRadiusKm * c
}

func calculateDeliveryFee(distanceKm float64) float64 {
    if distanceKm < 2 {
        return 30
    }
    return 30 + (distanceKm-2)*10
}

func distanceHandler(c *gin.Context) {
    var req entity.DistanceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    distance := haversineDistance(req.Lat1, req.Lng1, req.Lat2, req.Lng2)
	delivery_fee:= calculateDeliveryFee(distance)
    // res := DistanceResponse{Distance: distance}
    c.JSON(http.StatusOK, gin.H{"status":http.StatusOK, "distance_in_km":distance, "estimated_delivery_fee":delivery_fee})
}



func RegisterDistanceRoutes(rg *gin.RouterGroup) {
	distanceroute := rg.Group("/distance")
	distanceroute.POST("",distanceHandler)
}