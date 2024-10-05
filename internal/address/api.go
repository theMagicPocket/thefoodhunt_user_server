package address

import (
	"errors"
	"net/http"

	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
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
	r.PUT("/:userId/", res.update)
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, entity.NotOk(err.Error()))
		return
	}

	res, err := r.service.Add(ctx, userId, &addAddressReq)

	if err != nil {
		if errors.Is(err, appErr.ErrNoDocuments) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, entity.NotOk(err.Error()))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, entity.NotOk(err.Error()))
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
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, entity.NotOk(err.Error()))
		return
	}

	res, err := r.service.Update(ctx, userId, &addAddressReq)

	if err != nil {
		if errors.Is(err, appErr.ErrNoDocuments) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, entity.NotOk(err.Error()))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, entity.NotOk(err.Error()))
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

	_, err := r.service.Delete(ctx, userId, addressId)

	if err != nil {
		if errors.Is(err, appErr.ErrNoDocuments) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, entity.NotOk(err.Error()))
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, entity.NotOk(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, entity.Ok("address deleted", nil))
}
