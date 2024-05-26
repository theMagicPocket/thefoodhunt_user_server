package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hotel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	HotelName string             `bson:"hotel_name"`
	Cords     string             `bson:"cords"`
	Location  string             `bson:"location"`
	Phone     string             `bson:"phone"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
