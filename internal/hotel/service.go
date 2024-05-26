package hotel

import (
	"time"

	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service interface {
	Create(input CreateHotelRequest) (primitive.ObjectID, error)
}

type service struct {
	repo Repository
}

type Hotel struct {
	entity.Hotel
}

type CreateHotelRequest struct {
	HotelName string `json:"hotel_name"`
	Cords     string `json:"cords"`
	Location  string `json:"location"`
	Phone     string `json:"phone"`
}

func NewService(repo Repository) Service {
	return service{repo}
}

func (s service) Create(req CreateHotelRequest) (primitive.ObjectID, error) {

	// validate the incoming data
	// TODO: call validate on the request
	// add any neccesary fiels
	now := time.Now()

	newHotel := entity.Hotel{
		HotelName: req.HotelName,
		Cords:     req.Cords,
		Location:  req.Location,
		Phone:     req.Phone,
		CreatedAt: now,
		UpdatedAt: now,
	}

	id, err := s.repo.Create(newHotel)

	if err != nil {
		return primitive.ObjectID{}, nil
	}

	return id, nil
}
