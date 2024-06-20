package entity 
import (
	// "time"
);

type Restaurant struct{
	RestaurantID string `bson:"restaurant_id" json:"restaurant_id"`
}

type Voucher struct {
	// ID           string             `bson:"_id" json:"voucher_id"`
    VoucherCode    string             `json:"voucher_code" bson:"_id"` // Primary key
    MinCartValue   float64            `json:"min_cart_value" bson:"min_cart_value"`
    Percentage     float64            `json:"percentage" bson:"percentage"`
    MaxDiscountAmt float64            `json:"max_discount_amt" bson:"max_discount_amt"`
    Validity       string          `json:"validity" bson:"validity"`
    Restaurants    []Restaurant `json:"restaurants" bson:"restaurants"` // List of restaurant IDs as ObjectIDs
    Type           string             `json:"type" bson:"type"`                 // "user" or "restaurant"
	CreatedAt    string        `bson:"created_at" json:"created_at"`
}