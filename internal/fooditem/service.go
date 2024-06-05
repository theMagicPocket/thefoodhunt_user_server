package fooditem

import (
	"context"
	"errors"
	"fmt"
	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
    // "log"
)

type FoodItemServiceImpl struct {
	FoodItemCollection *mongo.Collection
	ctx            context.Context
}

type FoodItemService interface {
	CreateFoodItem(*entity.FoodItem) error
	GetFoodItem(*string) (*entity.FoodItem, error)
	GetAllFoodItems() ([]*entity.FoodItem, error)
	UpdateFoodItem(*entity.FoodItem) error
	DeleteFoodItem(*string) error
    GiveRating(*string,entity.RatingRequest) (*entity.FoodItem,error) 
}

func NewFoodItemService(FoodItemCollection *mongo.Collection, ctx context.Context) FoodItemService {
	return &FoodItemServiceImpl{
		FoodItemCollection: FoodItemCollection,
		ctx:            ctx,
	}
}


func (f *FoodItemServiceImpl) CreateFoodItem(fooditem *entity.FoodItem) error {
	_, err := f.FoodItemCollection.InsertOne(f.ctx, fooditem)
	return err
}

func (f *FoodItemServiceImpl) GetFoodItem(fooditemId *string) (*entity.FoodItem, error) {
	var fooditem *entity.FoodItem
	query := bson.D{bson.E{Key: "_id", Value: fooditemId}}
	err := f.FoodItemCollection.FindOne(f.ctx, query).Decode(&fooditem)
	return fooditem, err
}

func (f *FoodItemServiceImpl) GetAllFoodItems() ([]*entity.FoodItem, error) {
	var fooditems []*entity.FoodItem
	cursor, err := f.FoodItemCollection.Find(f.ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(f.ctx) {
		var fooditem entity.FoodItem
		err := cursor.Decode(&fooditem)
		if err != nil {
			return nil, err
		}
		fooditems = append(fooditems, &fooditem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(f.ctx)

	if len(fooditems) == 0 {
		return nil, errors.New("no Fooditems found")
	}
	return fooditems, nil
}


func (f *FoodItemServiceImpl) UpdateFoodItem(foodItem *entity.FoodItem) error {
    filter := bson.D{bson.E{Key: "_id", Value: foodItem.ID}}

    // Initialize an empty map for the update document
    updateFields := make(map[string]interface{})

    // Dynamically append fields to the update document if they are not zero values
    if foodItem.ItemName != "" {
        updateFields["item_name"] = foodItem.ItemName
    }
    if foodItem.Description != "" {
        updateFields["description"] = foodItem.Description
    }
    if foodItem.RestaurantID != "" {
        updateFields["restaurant_id"] = foodItem.RestaurantID
    }
    if foodItem.Price != 0 {
        updateFields["price"] = foodItem.Price
    }
    if foodItem.Photo != "" {
        updateFields["photo"] = foodItem.Photo
    }
    if foodItem.Ratings != 0 {
        updateFields["ratings"] = foodItem.Ratings
    }
    if foodItem.NoOfRatings != 0 {
        updateFields["no_of_ratings"] = foodItem.NoOfRatings
    }
    updateFields["is_veg"] = foodItem.IsVeg
    if foodItem.Category != "" {
        updateFields["category"] = foodItem.Category
    }

    // Handle updates for Addons array
    if len(foodItem.Addons) > 0 {
        for i, addon := range foodItem.Addons {
            addonPrefix := fmt.Sprintf("addons.%d.", i)
            if addon.Name != "" {
                updateFields[addonPrefix+"name"] = addon.Name
            }
            if addon.Price != 0 {
                updateFields[addonPrefix+"price"] = addon.Price
            }
        }
    }

    // Only proceed with the update if there are fields to update
    if len(updateFields) > 0 {
        updateDoc := bson.D{bson.E{Key: "$set", Value: updateFields}}

        result, err := f.FoodItemCollection.UpdateOne(f.ctx, filter, updateDoc)
        f.FoodItemCollection.FindOne(f.ctx, filter).Decode(&foodItem)
        if err != nil {
            return err
        }
        if result.MatchedCount != 1 {
            return errors.New("no matched document found for update")
        }
    }

    return nil
}

func (f *FoodItemServiceImpl) DeleteFoodItem(fooditemId *string) error {
	filter := bson.D{bson.E{Key: "_id", Value: fooditemId}}
	result, _ := f.FoodItemCollection.DeleteOne(f.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}
	return nil
}



func (f *FoodItemServiceImpl) GiveRating(fooditemId *string, ratingRequest entity.RatingRequest) (*entity.FoodItem, error) {
    filter := bson.D{bson.E{Key: "_id", Value: fooditemId}}
    var foodItem entity.FoodItem

    // Fetch the current food item from the database
    err := f.FoodItemCollection.FindOne(f.ctx, filter).Decode(&foodItem)
    if err != nil {
        return nil, err
    }

    // Increment the number of ratings
    foodItem.NoOfRatings += 1
    updateFields := bson.D{bson.E{Key: "no_of_ratings", Value: foodItem.NoOfRatings}}
    if(ratingRequest.Star == 1.0){
        foodItem.CountRatings.Onestar += 1
        updateFields = append(updateFields, bson.E{Key: "countratings.onestar", Value: foodItem.CountRatings.Onestar})
    }else if(ratingRequest.Star == 2.0){
        foodItem.CountRatings.Twostar += 1
        updateFields = append(updateFields, bson.E{Key: "countratings.twostar", Value: foodItem.CountRatings.Twostar})
    }else if(ratingRequest.Star == 3.0){
        foodItem.CountRatings.Threestar += 1
        updateFields = append(updateFields, bson.E{Key: "countratings.threestar", Value: foodItem.CountRatings.Threestar})
    }else if(ratingRequest.Star == 4.0){
        foodItem.CountRatings.FourStar += 1
        updateFields = append(updateFields, bson.E{Key: "countratings.fourstar", Value: foodItem.CountRatings.FourStar})
    }else if(ratingRequest.Star == 5.0){
        foodItem.CountRatings.FiveStar += 1
        updateFields = append(updateFields, bson.E{Key: "countratings.fivestar", Value: foodItem.CountRatings.FiveStar})
    }
    foodItem.Ratings = 1 * foodItem.CountRatings.Onestar + 2 * foodItem.CountRatings.Twostar + 3*foodItem.CountRatings.Threestar + 4*foodItem.CountRatings.FourStar + 5* foodItem.CountRatings.FiveStar
    foodItem.Ratings = foodItem.Ratings / float64(foodItem.NoOfRatings)
    updateFields = append(updateFields, bson.E{Key: "ratings", Value: foodItem.Ratings})
    

    // Update the document in the database
    updateDoc := bson.D{bson.E{Key: "$set", Value: updateFields}}
    _, err = f.FoodItemCollection.UpdateOne(f.ctx, filter, updateDoc)
    if err != nil {
        return nil, err
    }

    // Return the updated food item
    return &foodItem, nil
}


