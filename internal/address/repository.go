package address

import (
	"context"
	"errors"

	K "github.com/deVamshi/golang_food_delivery_api/internal/constants"
	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	appErr "github.com/deVamshi/golang_food_delivery_api/internal/errors"
	"github.com/deVamshi/golang_food_delivery_api/pkg/dbcontext"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	AddOrUpdate(ctx context.Context, userId string, address entity.UserAddress) (*entity.User, error)
	Delete(ctx context.Context, userId string, addressId string) (bool, error)
}

type repository struct {
	db *dbcontext.DB
}

func NewRepository(db *dbcontext.DB) Repository {
	return repository{db: db}
}

func (r repository) AddOrUpdate(ctx context.Context, userId string, address entity.UserAddress) (*entity.User, error) {

	user, err := getUser(ctx, r.db, userId)
	if err != nil {
		return nil, err
	}

	if user.UserAddress == nil {
		user.UserAddress = []entity.UserAddress{}
	}

	foundIdx := -1
	for id, adrs := range user.UserAddress {
		if adrs.ID == address.ID {
			foundIdx = id
			break
		}
	}

	if foundIdx == -1 {
		// address not found, adding
		user.UserAddress = append(user.UserAddress, address)
	} else {
		// address found, updating
		user.UserAddress[foundIdx] = address
	}

	update := bson.D{bson.E{Key: "$set", Value: user}}
	_, err = r.db.DB().Collection(K.USERS_COLLECTION).UpdateByID(ctx, userId, update)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r repository) Delete(ctx context.Context, userId string, addressId string) (bool, error) {

	user, err := getUser(ctx, r.db, userId)
	if err != nil {
		return false, err
	}

	if user.UserAddress == nil {
		user.UserAddress = []entity.UserAddress{}
	}

	foundIdx := -1
	for id, adrs := range user.UserAddress {
		if adrs.ID == addressId {
			foundIdx = id
			break
		}
	}
	if foundIdx != -1 {
		user.UserAddress = append(user.UserAddress[:foundIdx], user.UserAddress[foundIdx+1:]...)
	}

	update := bson.D{bson.E{Key: "$set", Value: user}}
	_, err = r.db.DB().Collection(K.USERS_COLLECTION).UpdateByID(ctx, userId, update)
	if err != nil {
		return false, err
	}

	return true, nil
}

// helpers
func getUser(ctx context.Context, db *dbcontext.DB, userId string) (*entity.User, error) {

	var user entity.User
	usersColl := db.DB().Collection(K.USERS_COLLECTION)

	filter := bson.D{bson.E{Key: "_id", Value: userId}}

	err := usersColl.FindOne(ctx, filter).Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, appErr.ErrNoDocuments
		}
		return nil, err
	}

	return &user, nil
}

// func (r repository) Add(ctx context.Context, id string, address entity.UserAddress) (*entity.User, error) {

// 	var user entity.User
// 	usersColl := r.db.DB().Collection("users")

// 	filter := bson.D{bson.E{Key: "_id", Value: id}}
// 	update := bson.D{{"$push", bson.E{"user_address", address}}}
// 	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

// 	err := usersColl.FindOneAndUpdate(ctx, filter, update, opts).Decode(&user)
// 	if err != nil {
// 		if errors.Is(err, mongo.ErrNoDocuments) {
// 			return nil, appErr.ErrNoDocuments
// 		}
// 		return nil, err
// 	}

// 	fmt.Println(user)

// 	return &user, nil
// }

// func (r repository) Add(ctx context.Context, userId string, address entity.UserAddress) (*entity.User, error) {

// 	user, err := getUser(ctx, r.db, userId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if user.UserAddress == nil {
// 		user.UserAddress = []entity.UserAddress{}
// 	}

// 	user.UserAddress = append(user.UserAddress, address)

// 	update := bson.D{bson.E{Key: "$set", Value: user}}
// 	_, err = r.db.DB().Collection(K.USERS_COLLECTION).UpdateByID(ctx, userId, update)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return user, nil
// }
