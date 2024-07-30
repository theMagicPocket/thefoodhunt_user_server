package hotels

import (
	"errors"
	"strconv"

	U "github.com/deVamshi/golang_food_delivery_api/internal"
	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
)

type Service interface {
	CreateHotel(*entity.Hotel) error
	GetHotel(string) (*entity.Hotel, error)
	GetAllHotels() ([]*entity.Hotel, error)
	UpdateHotel(*entity.Hotel) error
	DeleteHotel(string) error
	GetNearbyHotels(lat, long, dist float64) ([]*entity.Hotel, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return service{repo: repo}
}

func (s service) CreateHotel(hotel *entity.Hotel) error {
	return s.repo.Create(hotel)
}

func (s service) GetHotel(hotelID string) (*entity.Hotel, error) {
	return s.repo.GetHotelById(hotelID)
}

func (s service) GetAllHotels() ([]*entity.Hotel, error) {
	return s.repo.Query()
}

func (s service) UpdateHotel(hotel *entity.Hotel) error {

	updateFields := make(map[string]any)

	if hotel.Name != "" {
		updateFields["name"] = hotel.Name
	}
	if hotel.Photo != "" {
		updateFields["photo"] = hotel.Photo
	}
	if hotel.RestaurantAddress.Lat != "" {
		updateFields["restaurant_address.lat"] = hotel.RestaurantAddress.Lat
	}
	if hotel.RestaurantAddress.Long != "" {
		updateFields["restaurant_address.long"] = hotel.RestaurantAddress.Long
	}
	if hotel.RestaurantAddress.Street != "" {
		updateFields["restaurant_address.street"] = hotel.RestaurantAddress.Street
	}
	if hotel.RestaurantAddress.DoorNo != "" {
		updateFields["restaurant_address.doorNo"] = hotel.RestaurantAddress.DoorNo
	}
	if hotel.RestaurantAddress.Pincode != "" {
		updateFields["restaurant_address.pincode"] = hotel.RestaurantAddress.Pincode
	}
	if hotel.UpdatedAt != "" {
		updateFields["updated_at"] = hotel.UpdatedAt
	}

	return s.repo.Update(hotel.HotelID, updateFields)
}

func (s service) DeleteHotel(hotelID string) error {
	return s.repo.Delete(hotelID)
}

func (s service) GetNearbyHotels(userLat, userLong, dist float64) ([]*entity.Hotel, error) {

	result := []*entity.Hotel{}

	hotels, err := s.repo.Query()

	if err != nil {
		return result, err
	}

	for _, h := range hotels {

		// Parse the hotel's latitude and longitude
		hotelLat, err := strconv.ParseFloat(h.RestaurantAddress.Lat, 64)
		if err != nil {
			return nil, err
		}
		hotelLong, err := strconv.ParseFloat(h.RestaurantAddress.Long, 64)
		if err != nil {
			return nil, err
		}

		// Calculate the distance to the user's coordinates
		distance := U.HaversineDistance(userLat, userLong, hotelLat, hotelLong)
		if distance <= dist {
			result = append(result, h)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no nearby hotels found")
	}

	return result, nil
}
