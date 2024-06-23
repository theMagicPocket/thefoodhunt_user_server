package entity

type Order struct {
    OrderID     string  `bson:"_id,omitempty" json:"order_id"`
    CartValue   float64 `bson:"cart_value" json:"cart_value"`
    FoodItems   []FoodItem `bson:"menu_items" json:"menu_items"`
    Tip         float64 `bson:"tip,omitempty" json:"tip,omitempty"`
    VoucherCode string  `bson:"voucher_code,omitempty" json:"voucher_code,omitempty"`
    TotalPaid   float64 `bson:"total_paid" json:"total_paid"`
    Comments    string  `bson:"comments,omitempty" json:"comments,omitempty"`
    PaymentID   string  `bson:"payment_id" json:"payment_id"`
    CreatedTime string  `bson:"created_time" json:"created_time"`
    OrderStatus string  `bson:"order_status" json:"order_status"`
    UpdatedTime string  `bson:"updated_time,omitempty" json:"updated_time,omitempty"`
}

// type MenuItem struct {
//     ItemID   string `bson:"item_id" json:"item_id"`
//     Quantity int    `bson:"quantity" json:"quantity"`
// }
