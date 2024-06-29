package fooditem

import (
	"log"
	"net/http"
	"time"
	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"github.com/gin-gonic/gin"
)

type FoodItemController struct {
	FoodItemService FoodItemService
}

func New(fooditemService FoodItemService) FoodItemController {
	return FoodItemController{
		FoodItemService: fooditemService,
	}
}

func (fc *FoodItemController) CreateFoodItem(ctx *gin.Context) {
	var fooditem entity.FoodItem
	fooditem.ID = entity.GenerateID()
	istLocation,_:= time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(istLocation)
	formattedTime := now.Format("01-02-2006 15:04:05")
	fooditem.CreatedAt = formattedTime
	if err :=  ctx.ShouldBindJSON(&fooditem); err != nil {
	ctx.JSON(http.StatusBadRequest,gin.H{"message": err.Error()})
	return
	};
 err := fc.FoodItemService.CreateFoodItem(&fooditem)
 if err != nil{
  ctx.JSON(http.StatusBadGateway,gin.H{"message": err.Error()})
  return
 }
  ctx.JSON(http.StatusOK, gin.H{"message":"success","statuscode": http.StatusOK,"fooditemid":fooditem.ID})
}

func (fc *FoodItemController) GetFoodItem(ctx *gin.Context) {
	fooditemId := ctx.Param("id")
	fooditem, err := fc.FoodItemService.GetFoodItem((&fooditemId))
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"statuscode": http.StatusOK,"data": fooditem})
}

func (fc *FoodItemController) GetAllFoodItems(ctx *gin.Context) {
	fooditems, err := fc.FoodItemService.GetAllFoodItems()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"statuscode": http.StatusOK, "data":fooditems})

}

func (fc *FoodItemController) UpdateFoodItem(ctx *gin.Context) {
	var fooditem entity.FoodItem
	fooditem.ID = ctx.Param("id")
	if err := ctx.ShouldBindJSON(&fooditem); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := fc.FoodItemService.UpdateFoodItem(&fooditem)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	istLocation,_:= time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(istLocation)
	formattedTime := now.Format("01-02-2006 15:04:05")
	fooditem.CreatedAt = formattedTime
	ctx.JSON(http.StatusOK, gin.H{"message": "success","status":http.StatusOK,"updated":fooditem})

}

func (fc *FoodItemController) DeleteFoodItem(ctx *gin.Context) {
	fooditemId := ctx.Param("id")
	log.Println(fooditemId)
	err := fc.FoodItemService.DeleteFoodItem(&fooditemId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully Deleted Food Item","status":http.StatusOK})
}

func (fc *FoodItemController) GiveRating(ctx *gin.Context) {
	
	var ratingrequest entity.RatingRequest
	fooditemId := ctx.Param("id")
	if err := ctx.ShouldBindJSON(&ratingrequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	updatedFoodItem, err := fc.FoodItemService.GiveRating(&fooditemId, ratingrequest)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success", "status": http.StatusOK, "updated": updatedFoodItem})
}

func (fc *FoodItemController) GetMenu(ctx *gin.Context) {
	fooditemId := ctx.Param("id")
	fooditems, err := fc.FoodItemService.GetMenu(&fooditemId)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"statuscode": http.StatusOK, "data":fooditems})
}


func (fc *FoodItemController) RegisterFoodItemRoutes(rg *gin.RouterGroup) {
	fooditemroute := rg.Group("/fooditem")
	fooditemroute.POST("/create", fc.CreateFoodItem)
	fooditemroute.GET("/get/:id", fc.GetFoodItem)
	fooditemroute.GET("/getall", fc.GetAllFoodItems)
	fooditemroute.PATCH("/update/:id", fc.UpdateFoodItem)
	fooditemroute.DELETE("/delete/:id", fc.DeleteFoodItem)
	fooditemroute.POST("/rating/:id", fc.GiveRating)
	fooditemroute.GET("/menu/:id", fc.GetMenu)
}
