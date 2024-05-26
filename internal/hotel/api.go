package hotel

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers(rg *gin.RouterGroup, service Service) {

	res := resource{service: service}

	rg.GET("/hotels", res.get)
	rg.POST("/hotels", res.create)
}

type resource struct {
	service Service
}

func (r resource) get(c *gin.Context) {
	c.JSON(http.StatusOK, "hi")
}

func (r resource) create(c *gin.Context) {

	var data CreateHotelRequest

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	createdHotel, err := r.service.Create(data)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return

	}

	c.JSON(http.StatusCreated, createdHotel)
}
