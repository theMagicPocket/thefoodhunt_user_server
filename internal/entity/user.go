package entity

type UserAddress struct {
    Lat     string `bson:"lat,omitempty" json:"latitude,omitempty"`
    Long    string `bson:"long,omitempty" json:"longitude,omitempty"`
    Street  string `bson:"street,omitempty" json:"street,omitempty"`
    DoorNo  string `bson:"doorNo,omitempty" json:"doorno,omitempty"`
    Pincode string `bson:"pincode,omitempty" json:"pincode,omitempty"`
}

// User represents a user with related details.
type User struct {
    ID           string        `bson:"user_id,omitempty" json:"user_id,omitempty"`
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




