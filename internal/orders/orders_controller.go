package order

import (
    "net/http"
    "time"
    "github.com/deVamshi/golang_food_delivery_api/internal/entity"
    "github.com/gin-gonic/gin"
)

type OrderController struct {
    OrderService OrderService
}

func New(orderService OrderService) OrderController {
    return OrderController{
        OrderService: orderService,
    }
}

func (oc *OrderController) CreateOrder(ctx *gin.Context) {
    var order entity.Order
    order.OrderID = entity.GenerateID()
    istLocation, _ := time.LoadLocation("Asia/Kolkata")
    now := time.Now().In(istLocation)
    formattedTime := now.Format("01-02-2006 15:04:05")
    order.CreatedTime = formattedTime
	order.UpdatedTime = formattedTime
    if err := ctx.ShouldBindJSON(&order); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    err := oc.OrderService.CreateOrder(&order)
    if err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "success", "statuscode": http.StatusOK, "orderid": order.OrderID})
}

func (oc *OrderController) GetOrder(ctx *gin.Context) {
    orderId := ctx.Param("id")
    order, err := oc.OrderService.GetOrder(&orderId)
    if err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"statuscode": http.StatusOK, "data": order})
}

func (oc *OrderController) GetAllOrders(ctx *gin.Context) {
    orders, err := oc.OrderService.GetAllOrders()
    if err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"statuscode": http.StatusOK, "data": orders})
}

func (oc *OrderController) UpdateOrder(ctx *gin.Context) {
    var order entity.Order
    order.OrderID = ctx.Param("id")
    if err := ctx.ShouldBindJSON(&order); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
        return
    }
    istLocation, _ := time.LoadLocation("Asia/Kolkata")
    now := time.Now().In(istLocation)
    formattedTime := now.Format("01-02-2006 15:04:05")
    order.UpdatedTime = formattedTime
    err := oc.OrderService.UpdateOrder(&order)
    if err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{"message": "success", "status": http.StatusOK, "updated": order})
}

func (oc *OrderController) RegisterOrderRoutes(rg *gin.RouterGroup) {
    orderRoute := rg.Group("/order")
    orderRoute.POST("/create", oc.CreateOrder)
    orderRoute.GET("/get/:id", oc.GetOrder)
    orderRoute.GET("/getall", oc.GetAllOrders)
    orderRoute.PATCH("/update/:id", oc.UpdateOrder)
}
