package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/deVamshi/golang_food_delivery_api/internal/deliveryfee"
	"github.com/deVamshi/golang_food_delivery_api/internal/fooditem"
	"github.com/deVamshi/golang_food_delivery_api/internal/hotels"
	"github.com/deVamshi/golang_food_delivery_api/internal/tokenverification"
	"github.com/joho/godotenv"

	"github.com/deVamshi/golang_food_delivery_api/internal/matrixapi"
	order "github.com/deVamshi/golang_food_delivery_api/internal/orders"
	"github.com/deVamshi/golang_food_delivery_api/internal/payments"
	"github.com/deVamshi/golang_food_delivery_api/internal/user"
	"github.com/deVamshi/golang_food_delivery_api/internal/voucher"
	"github.com/deVamshi/golang_food_delivery_api/pkg/dbcontext"
	"github.com/deVamshi/golang_food_delivery_api/pkg/log"
	"github.com/gin-gonic/gin"
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

type Config struct {
	MONGODB_URI string
	SECRET_KEY  string
	MATRIX_KEY  string
	PORT        string
}

func main() {

	logger := log.New()

	err := godotenv.Load()
	if err != nil {
		logger.Error("Not loading from .env file:\n", err)
	}

	cfg := Config{}

	flag.StringVar(&cfg.MONGODB_URI, "MONGODB_URI", os.Getenv("MONGODB_URI"), "URI for mongodb")
	flag.StringVar(&cfg.SECRET_KEY, "SECRET_KEY", os.Getenv("SECRET_KEY"), "secret key")
	flag.StringVar(&cfg.MATRIX_KEY, "MATRIX_KEY", os.Getenv("MATRIX_KEY"), "matrix key")
	flag.StringVar(&cfg.PORT, "PORT", os.Getenv("PORT"), "port on which to run the server")

	flag.Parse()

	ctx = context.Background()
	r := gin.Default()

	MONGO_URI := cfg.MONGODB_URI
	if MONGO_URI == "" {
		logger.Fatal("mongouri is empty")
	}
	logger.Info("MONGOURI", MONGO_URI)
	dbClient, err := dbcontext.ConnectDB(MONGO_URI)
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
		ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message": "hi"})
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

	server := &http.Server{
		Addr:         ":" + cfg.PORT,
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
	}

	// logger.Fatal("unknown error occured")

	logger.Info(fmt.Sprintf("listening on http://localhost:%s", cfg.PORT))
	err = server.ListenAndServe()
	logger.Fatal(err)
}
