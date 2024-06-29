package entity



// type Hotel struct {
// 	ID        primitive.ObjectID `bson:"_id,omitempty"`
// 	HotelName string             `bson:"hotel_name"`
// 	Cords     string             `bson:"cords"`
// 	Location  string             `bson:"location"`
// 	Phone     string             `bson:"phone"`
// 	CreatedAt time.Time          `bson:"created_at"`
// 	UpdatedAt time.Time          `bson:"updated_at"`
// }

type Hotel struct {
    HotelID  string `json:"hotelid" bson:"_id, omitempty"`
    Name     string `json:"name" bson:"name"`
    Photo    string `json:"photo" bson:"photo"`
    RestaurantAddress  RestaurantAddress `bson:"restaurant_address" json:"restaurant_address"`
    Categories []string `json:"categories" bson:"categories"`
	CreatedAt string          `bson:"created_at"`
	UpdatedAt string         `bson:"updated_at"`
    Vouchers  []Voucher       `json:"vouchers" bson:"vouchers"`

}

// type Category struct{
//     Category string `bson:"category" json:"category"`
// }

type RestaurantAddress struct {
    Lat     string `bson:"lat,omitempty" json:"latitude,omitempty"`
    Long    string `bson:"long,omitempty" json:"longitude,omitempty"`
    Street  string `bson:"street,omitempty" json:"street,omitempty"`
    DoorNo  string `bson:"doorNo,omitempty" json:"doorno,omitempty"`
    Pincode string `bson:"pincode,omitempty" json:"pincode,omitempty"`
}
