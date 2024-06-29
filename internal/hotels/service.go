package hotel

import (
	"context"
	"errors"
	"math"
	"strconv"

	// "github.com/deVamshi/golang_food_delivery_api/internal/deliveryfee"
	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelServiceImpl struct {
    HotelCollection *mongo.Collection
    VoucherCollection *mongo.Collection
    ctx             context.Context
}

type HotelService interface {
    CreateHotel(*entity.Hotel) error
    GetHotel(*string) (*entity.Hotel, error)
    GetAllHotels() ([]*entity.Hotel, error)
    UpdateHotel(*entity.Hotel) error
    DeleteHotel(*string) error
    GetNearbyHotels(lat, long,dist float64) ([]*entity.Hotel, error) 
}

func NewHotelService(hotelCollection *mongo.Collection,voucherCollection *mongo.Collection, ctx context.Context) HotelService {
    return &HotelServiceImpl{
        HotelCollection: hotelCollection,
        VoucherCollection: voucherCollection,
        ctx:             ctx,
    }
}

func (h *HotelServiceImpl) CreateHotel(hotel *entity.Hotel) error {
    _, err := h.HotelCollection.InsertOne(h.ctx, hotel)
    return err
}


func GetVouchersByRestaurantID(ctx context.Context, voucherCollection *mongo.Collection, restaurantID string) ([]entity.Voucher, error) {
    var vouchers []entity.Voucher
    filter := bson.M{"restaurants.restaurant_id": restaurantID}
    cursor, err := voucherCollection.Find(ctx, filter)
    if err != nil {
        return nil, err
    }
    if err = cursor.All(ctx, &vouchers); err != nil {
        return nil, err
    }
    return vouchers, nil
}

func (h *HotelServiceImpl) GetHotel(hotelID *string) (*entity.Hotel, error) {
    var hotel entity.Hotel
    query := bson.D{bson.E{Key: "_id", Value: hotelID}}
    err := h.HotelCollection.FindOne(h.ctx, query).Decode(&hotel)
    if err != nil {
        return nil, err
    }
    vouchers, err := GetVouchersByRestaurantID(h.ctx, h.VoucherCollection, *hotelID)  // Assuming VoucherCollection is a field in HotelServiceImpl
    if err != nil {
        return nil, err
    }
    hotel.Vouchers = vouchers
    return &hotel, nil
}


func (h *HotelServiceImpl) GetAllHotels() ([]*entity.Hotel, error) {
    var hotels []*entity.Hotel
    cursor, err := h.HotelCollection.Find(h.ctx, bson.D{{}})
    if err != nil {
        return nil, err
    }
    for cursor.Next(h.ctx) {
        var hotel entity.Hotel
        err := cursor.Decode(&hotel)
        if err != nil {
            return nil, err
        }
        hotels = append(hotels, &hotel)
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    cursor.Close(h.ctx)

    if len(hotels) == 0 {
        return nil, errors.New("no hotels found")
    }
    return hotels, nil
}

func (h *HotelServiceImpl) UpdateHotel(hotel *entity.Hotel) error {
    filter := bson.D{bson.E{Key: "_id", Value: hotel.HotelID}}

    updateFields := make(map[string]interface{})

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

    if len(updateFields) > 0 {
        updateDoc := bson.D{bson.E{Key: "$set", Value: updateFields}}

        result, err := h.HotelCollection.UpdateOne(h.ctx, filter, updateDoc)
        if err != nil {
            return err
        }
        if result.MatchedCount != 1 {
            return errors.New("no matched document found for update")
        }
    }

    return nil
}


func (h *HotelServiceImpl) DeleteHotel(hotelID *string) error {
    filter := bson.D{bson.E{Key: "_id", Value: hotelID}}
    result, err := h.HotelCollection.DeleteOne(h.ctx, filter)
    if err != nil {
        return err
    }
    if result.DeletedCount == 0 {
        return errors.New("no matched document found for deletion")
    }
    return nil
}


func (h *HotelServiceImpl) GetNearbyHotels(userLat, userLong,dist float64) ([]*entity.Hotel, error) {
    var hotels []*entity.Hotel
    cursor, err := h.HotelCollection.Find(h.ctx, bson.D{{}})
    if err != nil {
        return nil, err
    }
    for cursor.Next(h.ctx) {
        var hotel entity.Hotel
        err := cursor.Decode(&hotel)
        if err != nil {
            return nil, err
        }

        // Parse the hotel's latitude and longitude
        hotelLat, err := strconv.ParseFloat(hotel.RestaurantAddress.Lat, 64)
        if err != nil {
            return nil, err
        }
        hotelLong, err := strconv.ParseFloat(hotel.RestaurantAddress.Long, 64)
        if err != nil {
            return nil, err
        }

        // Calculate the distance to the user's coordinates
        distance := haversineDistance(userLat, userLong, hotelLat, hotelLong)
        if distance <= dist {
            hotels = append(hotels, &hotel)
        }
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    cursor.Close(h.ctx)

    if len(hotels) == 0 {
        return nil, errors.New("no nearby hotels found")
    }
    return hotels, nil
}


const earthRadiusKm = 6371
func toRadians(angle float64) float64 {
    return angle * math.Pi / 180
}

func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
    dLat := toRadians(lat2 - lat1)
    dLng := toRadians(lng2 - lng1)

    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Cos(toRadians(lat1))*math.Cos(toRadians(lat2))*
            math.Sin(dLng/2)*math.Sin(dLng/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    return earthRadiusKm * c
}