package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	userCollection *mongo.Collection
	ctx            context.Context
}

type UserService interface {
	CreateUser(*entity.User) error
	GetUser(*string) (*entity.User, error)
	GetAll() ([]*entity.User, error)
	UpdateUser(*entity.User) error
	DeleteUser(*string) error
}

func NewUserService(userCollection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{
		userCollection: userCollection,
		ctx:            ctx,
	}
}
var ErrEmailExists = errors.New("email already in use")

func (u *UserServiceImpl) CreateUser(user *entity.User) error {
	var existingUser entity.User
    err := u.userCollection.FindOne(u.ctx, bson.M{"email": user.Email}).Decode(&existingUser)
    if err == nil {
        // Email already exists
		
        return ErrEmailExists
    } else if err != mongo.ErrNoDocuments {
        // An error occurred other than "no documents found"
        return err
    }
	_, errr := u.userCollection.InsertOne(u.ctx, user)
	return errr
}

func (u *UserServiceImpl) GetUser(userId *string) (*entity.User, error) {
	var user *entity.User
	query := bson.D{bson.E{Key: "user_id", Value: userId}}
	err := u.userCollection.FindOne(u.ctx, query).Decode(&user)
	return user, err
}

func (u *UserServiceImpl) GetAll() ([]*entity.User, error) {
	var users []*entity.User
	cursor, err := u.userCollection.Find(u.ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(u.ctx) {
		var user entity.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(u.ctx)

	if len(users) == 0 {
		return nil, errors.New("documents not found")
	}
	return users, nil
}


func (u *UserServiceImpl) UpdateUser(user *entity.User) error {
	filter := bson.D{bson.E{Key: "user_id", Value: user.ID}}

	// Initialize an empty map for the update document
	updateFields := make(map[string]interface{})

	// Dynamically append fields to the update document if they are not zero values
	if user.Name != "" {
		updateFields["Name"] = user.Name
	}
	if user.Email != "" {
		updateFields["email"] = user.Email
	}
	if user.Phone != "" {
		updateFields["phone"] = user.Phone
	}
	if user.ProfilePhoto != "" {
		updateFields["profile_photo"] = user.ProfilePhoto
	}

	// Handle updates for UserAddress array
	if len(user.UserAddress) > 0 {
		for i, address := range user.UserAddress {
			addressPrefix := fmt.Sprintf("user_address.%d.", i)
			if address.Lat != "" {
				updateFields[addressPrefix+"lat"] = address.Lat
			}
			if address.Long != "" {
				updateFields[addressPrefix+"long"] = address.Long
			}
			if address.Street != "" {
				updateFields[addressPrefix+"street"] = address.Street
			}
			if address.DoorNo != "" {
				updateFields[addressPrefix+"doorNo"] = address.DoorNo
			}
			if address.Pincode != "" {
				updateFields[addressPrefix+"pincode"] = address.Pincode
			}
		}
	}

	// Only proceed with the update if there are fields to update
	if len(updateFields) > 0 {
		updateDoc := bson.D{bson.E{Key: "$set", Value: updateFields}}

		result, err := u.userCollection.UpdateOne(u.ctx, filter, updateDoc)
		if err != nil {
			return err
		}
		if result.MatchedCount != 1 {
			return errors.New("no matched document found for update")
		}
	}

	return nil
}

func (u *UserServiceImpl) DeleteUser(userId *string) error {
	filter := bson.D{bson.E{Key: "user_id", Value: userId}}
	result, _ := u.userCollection.DeleteOne(u.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}
	return nil
}
