package entity

type UserAddress struct {
	ID          string `bson:"id" json:"id" binding:"required"`
	Lat         string `bson:"lat,omitempty" json:"latitude,omitempty"`
	Long        string `bson:"long,omitempty" json:"longitude,omitempty"`
	Street      string `bson:"street,omitempty" json:"street,omitempty"`
	DoorNo      string `bson:"doorNo,omitempty" json:"doorno,omitempty"`
	Pincode     string `bson:"pincode,omitempty" json:"pincode,omitempty"`
	Name        string `bson:"name,omitempty" json:"name,omitempty"`
	PhoneNumber string `bson:"phone_number,omitempty" json:"phone_number,omitempty"`
	Landmark    string `bson:"landmark,omitempty" json:"landmark,omitempty"`
	IsActive    bool   `bson:"is_active,omitempty" json:"is_active,omitempty"`
}

// User represents a user with related details.
type User struct {
	Id           string        `bson:"_id,omitempty" json:"id"`
	AuthId       string        `bson:"auth_id" json:"auth_id"`
	Name         string        `bson:"name" json:"name"`
	Email        string        `bson:"email" json:"email"`
	Phone        string        `bson:"phone" json:"phone"`
	ProfilePhoto string        `bson:"profile_photo" json:"profile_photo"`
	UserAddress  []UserAddress `bson:"user_address" json:"user_address"`
	CreatedAt    string        `bson:"created_at" json:"created_at"`
}

// type UserAddress struct {
//     Lat    string `bson:"lat,omitempty"`
//     Long   string `bson:"long,omitempty"`
//     Street string `bson:"street,omitempty"`
//     DoorNo string `bson:"doorNo,omitempty"`
//     Pincode string `bson:"pincode,omitempty"`
// }

// // User represents a user with related details.
// type User struct {
//     ID        string `bson:"id,omitempty"`
//     Name         string        `bson:"name"`
//     Email        string        `bson:"email"`
//     Phone        string        `bson:"phone"`
//     ProfilePhoto string        `bson:"profilePhoto"`
//     UserAddress  []UserAddress `bson:"userAddress"`
//     CreatedAt     string     `bson:"created_at"`
// }
