package entity

type Addon struct {
	Name  string `bson:"name" json:"name"`
	Price int    `bson:"price" json:"price"`
}

type RatingRequest struct {
	Star float64 `json:"star"`
}

type CountRating struct {
	Onestar   float64 `bson:"onestar" json:"onestar"`
	Twostar   float64 `bson:"twostar" json:"twostar"`
	Threestar float64 `bson:"threestar" json:"threestar"`
	FourStar  float64 `bson:"fourstar" json:"fourstar"`
	FiveStar  float64 `bson:"fivestar" json:"fivestar"`
}

type FoodItem struct {
	ID           string      `bson:"_id,omitempty" json:"item_id"`
	ItemName     string      `bson:"item_name" json:"item_name"`
	Description  string      `bson:"description" json:"description"`
	RestaurantID string      `bson:"restaurant_id" json:"restaurant_id"`
	Addons       []Addon     `bson:"addons" json:"addons"`
	Price        int         `bson:"price" json:"price"`
	Photo        string      `bson:"photo" json:"photo"`
	Ratings      float64     `bson:"ratings" json:"ratings"`
	CountRatings CountRating `bson:"countratings" json:"countratings"`
	NoOfRatings  int         `bson:"no_of_ratings" json:"no_of_ratings"`
	IsVeg        bool        `bson:"is_veg" json:"is_veg"`
	Category     string      `bson:"category" json:"category"`
	CreatedAt    string      `bson:"created_at" json:"created_at"`
	Active       bool        `bson:"is_active" json:"is_active"`
}
