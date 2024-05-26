package main

import (
	"log"
	"net/http"

	"github.com/deVamshi/golang_food_delivery_api/internal/hotel"
	"github.com/deVamshi/golang_food_delivery_api/pkg/dbcontext"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load env vars
	APP_ENV, err := godotenv.Read("../.env")
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
		hotel.RegisterHandlers(v1, hotel.NewService(hotel.NewRepository(dbClient)))
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
