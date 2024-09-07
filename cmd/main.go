package main

import (
	"context"
	"log"
	"net/http"

	"github.com/deVamshi/golang_food_delivery_api/internal/deliveryfee"
	"github.com/deVamshi/golang_food_delivery_api/internal/fooditem"
	hotel "github.com/deVamshi/golang_food_delivery_api/internal/hotels"

	// "github.com/deVamshi/golang_food_delivery_api/internal/hotel"
	"github.com/deVamshi/golang_food_delivery_api/internal/matrixapi"
	order "github.com/deVamshi/golang_food_delivery_api/internal/orders"
	"github.com/deVamshi/golang_food_delivery_api/internal/payments"
	"github.com/deVamshi/golang_food_delivery_api/internal/tokenverification"
	"github.com/deVamshi/golang_food_delivery_api/internal/user"
	"github.com/deVamshi/golang_food_delivery_api/internal/voucher"
	"github.com/deVamshi/golang_food_delivery_api/pkg/dbcontext"
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
	// hotels
	hotelservice    hotel.HotelService
	hotelcontroller hotel.HotelController
)

func main() {
	// load env vars
	ctx = context.Background()
	APP_ENV, err := godotenv.Read(".env")
	if err != nil {
		log.Fatal(err)
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	dbClient, err := dbcontext.ConnectDB(APP_ENV["MONGOURI"])
	defer func() {
		if err := dbcontext.DisconnectDB(); err != nil {
			log.Fatal(err)
		}
		log.Print("Disconnected to DB\nBye!")
	}()

	if err != nil {
		log.Fatal(err)
		return
	}

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
		var hotelscollection = dbClient.Collection("hotels")
		hotelservice = hotel.NewHotelService(hotelscollection, voucherscollection, ctx)
		hotelcontroller = hotel.New(hotelservice)
		hotelcontroller.RegisterHotelRoutes(v1)
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
		Addr:    ":8080",
		Handler: r,
	}

	// func() {
	// 	time.Sleep(3 * time.Second)
	// 	server.Shutdown(context.Background())
	// }()

	server.ListenAndServe()
}
