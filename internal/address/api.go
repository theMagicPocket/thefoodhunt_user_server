package address

import (
	"errors"
	"net/http"

	appErr "github.com/deVamshi/golang_food_delivery_api/internal/errors"
	"github.com/deVamshi/golang_food_delivery_api/pkg/log"
	"github.com/gin-gonic/gin"
)

type resource struct {
	service Service
	logger  *log.AppLogger
}

func RegisterHandlers(rg *gin.RouterGroup, service Service, logger *log.AppLogger) {

	res := resource{service: service, logger: logger}

	r := rg.Group("/address")

	r.POST("/:userId", res.create)
	r.PUT("/:userId", res.update)
	r.DELETE("/:userId/:addressId", res.delete)
}

func (r *resource) create(ctx *gin.Context) {

	userId := ctx.Param("userId")
	if userId == "" {
		ctx.AbortWithError(http.StatusBadRequest, appErr.ErrNoRequiredParam("userId"))
		return
	}

	var addAddressReq AddOrUpdateAddressRequest
	if err := ctx.ShouldBindJSON(&addAddressReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := r.service.Add(ctx, userId, &addAddressReq)

	if err != nil {
		if errors.Is(err, appErr.ErrNoDocuments) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, res)
}

func (r *resource) update(ctx *gin.Context) {

	userId := ctx.Param("userId")
	if userId == "" {
		ctx.AbortWithError(http.StatusBadRequest, appErr.ErrNoRequiredParam("userId"))
		return
	}

	var addAddressReq AddOrUpdateAddressRequest
	if err := ctx.ShouldBindJSON(&addAddressReq); err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	res, err := r.service.Update(ctx, userId, &addAddressReq)

	if err != nil {
		if errors.Is(err, appErr.ErrNoDocuments) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (r *resource) delete(ctx *gin.Context) {

	userId := ctx.Param("userId")
	if userId == "" {
		ctx.AbortWithError(http.StatusBadRequest, appErr.ErrNoRequiredParam("userId"))
		return
	}

	addressId := ctx.Param("addressId")
	if userId == "" {
		ctx.AbortWithError(http.StatusBadRequest, appErr.ErrNoRequiredParam("addressId"))
		return
	}

	isDeleted, err := r.service.Delete(ctx, userId, addressId)

	if err != nil {
		if errors.Is(err, appErr.ErrNoDocuments) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}

	if !isDeleted {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "no address deleted, check the id"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"messsage": "deleted the address"})
}
