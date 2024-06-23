package order

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	// "time"
)

type OrderServiceImpl struct {
    OrderCollection *mongo.Collection
    ctx             context.Context
}

type OrderService interface {
    CreateOrder(*entity.Order) error
    GetOrder(*string) (*entity.Order, error)
    GetAllOrders() ([]*entity.Order, error)
    UpdateOrder(*entity.Order) error
}

func NewOrderService(OrderCollection *mongo.Collection, ctx context.Context) OrderService {
    return &OrderServiceImpl{
        OrderCollection: OrderCollection,
        ctx:             ctx,
    }
}

func (o *OrderServiceImpl) CreateOrder(order *entity.Order) error {
	log.Println("came")
    _, err := o.OrderCollection.InsertOne(o.ctx, order)
    return err
}

func (o *OrderServiceImpl) GetOrder(orderID *string) (*entity.Order, error) {
    var order *entity.Order
    query := bson.D{bson.E{Key: "_id", Value: orderID}}
    err := o.OrderCollection.FindOne(o.ctx, query).Decode(&order)
    return order, err
}

func (o *OrderServiceImpl) GetAllOrders() ([]*entity.Order, error) {
    var orders []*entity.Order
    cursor, err := o.OrderCollection.Find(o.ctx, bson.D{{}})
    if err != nil {
        return nil, err
    }
    for cursor.Next(o.ctx) {
        var order entity.Order
        err := cursor.Decode(&order)
        if err != nil {
            return nil, err
        }
        orders = append(orders, &order)
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    cursor.Close(o.ctx)

    if len(orders) == 0 {
        return nil, errors.New("no orders found")
    }
    return orders, nil
}

func (o *OrderServiceImpl) UpdateOrder(order *entity.Order) error {
    filter := bson.D{bson.E{Key: "_id", Value: order.OrderID}}

    // Initialize an empty map for the update document
    updateFields := make(map[string]interface{})

    // Dynamically append fields to the update document if they are not zero values or empty
    if order.CartValue != 0 {
        updateFields["cart_value"] = order.CartValue
    }
    if order.Tip != 0 {
        updateFields["tip"] = order.Tip
    }
    if order.VoucherCode != "" {
        updateFields["voucher_code"] = order.VoucherCode
    }
    if order.TotalPaid != 0 {
        updateFields["total_paid"] = order.TotalPaid
    }
    if order.Comments != "" {
        updateFields["comments"] = order.Comments
    }
    if order.PaymentID != "" {
        updateFields["payment_id"] = order.PaymentID
    }
    if order.OrderStatus != "" {
        updateFields["order_status"] = order.OrderStatus
    }
    if order.UpdatedTime != "" {
        updateFields["updated_time"] = order.UpdatedTime
    }

    // Handle updates for FoodItems array
    if len(order.FoodItems) > 0 {
        for i, foodItem := range order.FoodItems {
            itemPrefix := fmt.Sprintf("menu_items.%d.", i)
            if foodItem.ItemName != "" {
                updateFields[itemPrefix+"item_name"] = foodItem.ItemName
            }
            if foodItem.Description != "" {
                updateFields[itemPrefix+"description"] = foodItem.Description
            }
            if foodItem.RestaurantID != "" {
                updateFields[itemPrefix+"restaurant_id"] = foodItem.RestaurantID
            }
            if foodItem.Price != 0 {
                updateFields[itemPrefix+"price"] = foodItem.Price
            }
            if foodItem.Photo != "" {
                updateFields[itemPrefix+"photo"] = foodItem.Photo
            }
            if foodItem.Ratings != 0 {
                updateFields[itemPrefix+"ratings"] = foodItem.Ratings
            }
            if foodItem.NoOfRatings != 0 {
                updateFields[itemPrefix+"no_of_ratings"] = foodItem.NoOfRatings
            }
            updateFields[itemPrefix+"is_veg"] = foodItem.IsVeg
            if foodItem.Category != "" {
                updateFields[itemPrefix+"category"] = foodItem.Category
            }
        }
    }

    // Only proceed with the update if there are fields to update
    if len(updateFields) > 0 {
        updateDoc := bson.D{bson.E{Key: "$set", Value: updateFields}}

        result, err := o.OrderCollection.UpdateOne(o.ctx, filter, updateDoc)
        o.OrderCollection.FindOne(o.ctx, filter).Decode(&order)
        if err != nil {
            return err
        }
        if result.MatchedCount != 1 {
            return errors.New("no matched document found for update")
        }
    }

    return nil
}
