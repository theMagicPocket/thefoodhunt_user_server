package main

import (
	"context"
	"net/http"
	"os"
	"time"
	"github.com/deVamshi/golang_food_delivery_api/internal/deliveryfee"
	"github.com/deVamshi/golang_food_delivery_api/internal/fooditem"
	"github.com/deVamshi/golang_food_delivery_api/internal/hotels"
	"github.com/deVamshi/golang_food_delivery_api/internal/tokenverification"

	"github.com/deVamshi/golang_food_delivery_api/internal/matrixapi"
	order "github.com/deVamshi/golang_food_delivery_api/internal/orders"
	"github.com/deVamshi/golang_food_delivery_api/internal/payments"
	"github.com/deVamshi/golang_food_delivery_api/internal/user"
	"github.com/deVamshi/golang_food_delivery_api/internal/voucher"
	"github.com/deVamshi/golang_food_delivery_api/pkg/dbcontext"
	"github.com/deVamshi/golang_food_delivery_api/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	userservice    user.UserService
	usercontroller user.UserController
	ctx            context.Context
)

var (
	fooditemservice    fooditem.FoodItemService
	fooditemcontroller fooditem.FoodItemController
	voucherservice     voucher.VoucherService
	vouchercontroller  voucher.VoucherController
	orderservice       order.OrderService
	ordercontroller    order.OrderController
)

func main() {

	logger := log.New()

	ctx = context.Background()
	APP_ENV, err := godotenv.Read(".env")
	if err != nil {
		logger.Fatal(err)
		logger.Fatal("Error loading .env file")
	}

	r := gin.Default()

	dbClient, err := dbcontext.ConnectDB(APP_ENV["MONGOURI"])
	defer func() {
		if err := dbcontext.DisconnectDB(); err != nil {
			logger.Fatal(err)
		}
		logger.Info("Disconnected to DB\nBye!")
	}()

	if err != nil {
		logger.Fatal(err)
	}

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message": "welcome to yummy-foods"})
	})

	v1 := r.Group("/v1")
	{
		v1.Use(tokenverification.AuthMiddleware())
		// hotel.RegisterHandlers(v1, hotel.NewService(hotel.NewRepository(dbClient)))
		// vouchers
		var voucherscollection = dbClient.Collection("vouchers")
		voucherservice = voucher.NewVoucherService(voucherscollection, ctx)
		vouchercontroller = voucher.New(voucherservice)
		vouchercontroller.RegisterVoucherRoutes(v1)
		// hotels
		hotels.RegisterHandlers(v1, hotels.NewService(hotels.NewRepository(dbClient)), logger)
		// users
		var usercollection = dbClient.Collection("users")
		userservice = user.NewUserService(usercollection, ctx)
		usercontroller = user.New(userservice)
		usercontroller.RegisterUserRoutes(v1)
		// fooditems
		var itemscollection = dbClient.Collection("fooditems")
		fooditemservice = fooditem.NewFoodItemService(itemscollection, ctx)
		fooditemcontroller = fooditem.New(fooditemservice)
		fooditemcontroller.RegisterFoodItemRoutes(v1)
		// orders
		var orderscollection = dbClient.Collection("orders")
		orderservice = order.NewOrderService(orderscollection, ctx)
		ordercontroller = order.New(orderservice)
		ordercontroller.RegisterOrderRoutes(v1)
		// payments
		payments.RegisterPaymentRoutes(v1)
		//distance calculator and estimated delivery fee
		deliveryfee.RegisterDistanceRoutes(v1)
		// exact distance calculation using google matrix api
		matrixapi.RegisterDistanceMatrixRoutes(v1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,

	}

	// logger.Fatal("unknown error occured")

	err = server.ListenAndServe()
	logger.Fatal(err)
}
