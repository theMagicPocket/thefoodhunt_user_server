package hotel

import (
	"context"

	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Create(entity.Hotel) (primitive.ObjectID, error)
}

type repository struct {
	hotels *mongo.Collection
}

func NewRepository(db *mongo.Database) Repository {
	return repository{db.Collection("hotels")}
}

func (repo repository) Get(_id primitive.ObjectID) {

	repo.hotels.FindOne(context.TODO(), _id)

}

func (repo repository) Create(hotel entity.Hotel) (primitive.ObjectID, error) {
	insertedId, err := repo.hotels.InsertOne(context.TODO(), hotel)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return insertedId.InsertedID.(primitive.ObjectID), nil
}
