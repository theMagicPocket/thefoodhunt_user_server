package hotel

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterHandlers(rg *gin.RouterGroup, service Service) {

	res := resource{service: service}

	rg.GET("/hotels", res.get)

}

type resource struct {
	service Service
}

func (r resource) get(c *gin.Context) {
	c.JSON(http.StatusOK, "hi")
}
