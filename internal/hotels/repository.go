package hotels

import (
	"context"
	"errors"

	C "github.com/deVamshi/golang_food_delivery_api/internal/common"
	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(*entity.Hotel) error
	GetHotelById(string) (*entity.Hotel, error)
	Query() ([]*entity.Hotel, error)
	Update(string, map[string]any) error
	Delete(string) error
}

type repository struct {
	db     *mongo.Database
	hotels *mongo.Collection
}

func NewRepository(db *mongo.Database) Repository {
	return repository{db: db, hotels: db.Collection(C.HOTELS_COLLECTION)}
}

func (r repository) Create(newHotel *entity.Hotel) error {
	_, err := r.hotels.InsertOne(context.TODO(), newHotel)
	return err
}

func getVouchersByRestaurantID(ctx context.Context, voucherCollection *mongo.Collection, restaurantID string) ([]entity.Voucher, error) {
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

func (r repository) GetHotelById(hotelID string) (*entity.Hotel, error) {
	var hotel entity.Hotel
	query := bson.D{bson.E{Key: "_id", Value: hotelID}}
	err := r.hotels.FindOne(context.TODO(), query).Decode(&hotel)
	if err != nil {
		return nil, err
	}
	vouchers, err := getVouchersByRestaurantID(context.TODO(), r.db.Collection(C.VOUCHERS_COLLECTION), hotelID) // Assuming VoucherCollection is a field in HotelServiceImpl
	if err != nil {
		return nil, err
	}
	hotel.Vouchers = vouchers
	return &hotel, nil
}

func (r repository) Query() ([]*entity.Hotel, error) {
	var hotels []*entity.Hotel

	ctx := context.TODO()

	cursor, err := r.hotels.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
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

	cursor.Close(ctx)

	if len(hotels) == 0 {
		return nil, errors.New("no hotels found")
	}
	return hotels, nil
}

func (r repository) Update(hotelId string, updateFields map[string]any) error {
	filter := bson.D{bson.E{Key: "_id", Value: updateFields}}

	if len(updateFields) > 0 {
		updateDoc := bson.D{bson.E{Key: "$set", Value: updateFields}}

		result, err := r.hotels.UpdateOne(context.TODO(), filter, updateDoc)
		if err != nil {
			return err
		}

		if result.MatchedCount != 1 {
			return errors.New("no matched document found for update")
		}
	}

	return nil
}

func (r repository) Delete(hotelID string) error {
	filter := bson.D{bson.E{Key: "_id", Value: hotelID}}
	result, err := r.hotels.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("no matched document found for deletion")
	}
	return nil
}
