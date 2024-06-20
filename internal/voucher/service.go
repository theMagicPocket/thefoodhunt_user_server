package voucher

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type VoucherServiceImpl struct {
	voucherCollection *mongo.Collection
	ctx            context.Context
}

type VoucherService interface {
	CreateVoucher(*entity.Voucher) error
	GetVoucher(*string) (*entity.Voucher, error)
	GetAll() ([]*entity.Voucher, error)
	UpdateVoucher(*entity.Voucher) error
	DeleteVoucher(*string) error
	ValidateVoucher(*string, *float64) (float64,float64,float64,error)
}

func NewVoucherService(voucherCollection *mongo.Collection, ctx context.Context) VoucherService {
	return &VoucherServiceImpl{
		voucherCollection: voucherCollection,
		ctx:            ctx,
	}
}

var ErrMinimumAmt = errors.New("minimum amount")

func (v *VoucherServiceImpl) CreateVoucher(Voucher *entity.Voucher) error {
	// log.Println(time.Now())
	_, err := v.voucherCollection.InsertOne(v.ctx, Voucher)
	return err
}

func (v *VoucherServiceImpl) GetVoucher(voucherId *string) (*entity.Voucher, error) {
	var Voucher *entity.Voucher
	query := bson.D{bson.E{Key: "_id", Value: voucherId}}
	err := v.voucherCollection.FindOne(v.ctx, query).Decode(&Voucher)
	return Voucher, err
}

func (v *VoucherServiceImpl) GetAll() ([]*entity.Voucher, error) {
	var vouchers []*entity.Voucher
	cursor, err := v.voucherCollection.Find(v.ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(v.ctx) {
		var Voucher entity.Voucher
		err := cursor.Decode(&Voucher)
		if err != nil {
			return nil, err
		}
		vouchers = append(vouchers, &Voucher)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(v.ctx)

	if len(vouchers) == 0 {
		return nil, errors.New("documents not found")
	}
	return vouchers, nil
}


// func (v *VoucherServiceImpl) UpdateVoucher(Voucher *entity.Voucher) error {
//     filter := bson.D{bson.E{Key: "_id", Value: Voucher.VoucherCode}}

//     // Initialize an empty map for the update document
//     updateFields := make(map[string]interface{})

//     // Dynamically append fields to the update document if they are not zero values
//     if Voucher.VoucherCode != "" {
//         updateFields["voucher_code"] = Voucher.VoucherCode
//     }
//     if Voucher.MinCartValue !=  0{
//         updateFields["min_cart_value"] = Voucher.MinCartValue
//     }
//     if Voucher.Percentage != 0 {
//         updateFields["percentage"] = Voucher.Percentage
//     }
//     if Voucher.MaxDiscountAmt != 0 {
//         updateFields["max_discount_amt"] = Voucher.MaxDiscountAmt
//     }
//     // if Voucher.Validity != {
//         updateFields["validity"] = Voucher.Validity
//     // }
    

//     // Handle updates for Addons array
//     if len(Voucher.Restaurants) > 0 {
//         for i, restaurant := range Voucher.Restaurants {
//             restaurantPrefix := fmt.Sprintf("restaurants.%d.", i)
//             if restaurant.RestaurantID != "" {
//                 updateFields[restaurantPrefix+"id"] = restaurant.RestaurantID
//             }
//         }
//     }

//     // Only proceed with the update if there are fields to update
//     if len(updateFields) > 0 {
//         updateDoc := bson.D{bson.E{Key: "$set", Value: updateFields}}

//         result, err := v.voucherCollection.UpdateOne(v.ctx, filter, updateDoc)
//         v.voucherCollection.FindOne(v.ctx, filter).Decode(&Voucher)
//         if err != nil {
//             return err
//         }
//         if result.MatchedCount != 1 {
//             return errors.New("no matched document found for update")
//         }
//     }

//     return nil
// }

func (v *VoucherServiceImpl) DeleteVoucher(voucherId *string) error {
	filter := bson.D{bson.E{Key: "_id", Value: voucherId}}
	result, _ := v.voucherCollection.DeleteOne(v.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched document found for delete")
	}
	return nil
}


func (v *VoucherServiceImpl) UpdateVoucher(Voucher *entity.Voucher) error {
	filter := bson.D{{Key: "_id", Value: Voucher.VoucherCode}}
	// Create updates variable to hold all the update fields
	updates := bson.D{}
	// newVoucher := *Voucher
	// Get the type of struct == Voucher
	typeData := reflect.TypeOf(*Voucher)

	// Get the values from the provided object: newVoucher
	values := reflect.ValueOf(*Voucher)

	// Starting from index 0 to include all fields
	for i := 0; i < typeData.NumField(); i++ {
		field := typeData.Field(i)   // Get the field from the struct definition
		val := values.Field(i)       // Get the value from the specified field position
		tag := field.Tag.Get("json") // From the field, get the json struct tag

		// We want to avoid zero values, as the omitted fields from newVoucher
		// correspond to their zero values, and we only want provided fields
		if tag != "" && !isZeroType(val) {
			log.Println("Updating field:", tag, "with value:", val)
			update := bson.E{Key: tag, Value: val.Interface()}
			updates = append(updates, update)
		}
	}

	if len(updates) == 0 {
		log.Println("No updates to be made, all fields are zero values.")
		return nil
	}

	updateFilter := bson.D{{Key: "$set", Value: updates}}
	log.Println("Update filter:", updateFilter)
	_, updateErr := v.voucherCollection.UpdateOne(v.ctx, filter, updateFilter)
	if updateErr != nil {
		log.Fatalf("Error updating voucher: %v", updateErr)
		return updateErr
	}

	return nil
}

// isZeroType checks if the value from the struct is the zero value of its type
func isZeroType(value reflect.Value) bool {
	zero := reflect.Zero(value.Type()).Interface()

	switch value.Kind() {
	case reflect.Slice, reflect.Array, reflect.Chan, reflect.Map:
		return value.Len() == 0
	default:
		return reflect.DeepEqual(zero, value.Interface())
}
}


func (v *VoucherServiceImpl) ValidateVoucher(voucherCode *string, cartvalue *float64) (float64,float64,float64, error) {
	
	// Filter to find the voucher by voucher code
	filter := bson.D{bson.E{Key: "_id", Value: *voucherCode}}

	var voucher entity.Voucher
	err := v.voucherCollection.FindOne(v.ctx, filter).Decode(&voucher)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return 0,*cartvalue,0, errors.New("voucher code does not exist")
		}
		return 0,*cartvalue,0, err
	}
	layout := "2006-01-02 15:04:05"

	// Time input string
	timeStr1 := voucher.Validity

	// // Parse the time string
	// time1, err1 := time.Parse(layout, timeStr1)
	// if err1 != nil {
	// 	fmt.Println("Error parsing time1:", err1)
	// 	return 0,err1
	// }
	loc, err2 := time.LoadLocation("Asia/Kolkata")
    if err2 != nil {
        fmt.Println("Error loading location:", err2)
        return 0,*cartvalue,0, err2
    }

    // Get the current time in IST
    now := time.Now().In(loc)
	timeInIST, err := time.ParseInLocation(layout, timeStr1, loc)
    if err != nil {
        fmt.Println("Error parsing time in IST:", err)
        return 0,*cartvalue,0, err
    }

	log.Println(now)
	log.Println(timeInIST)
	
	if timeInIST.Before(now) {
		return 0,*cartvalue,0, errors.New("voucher has been expired")
	}
	// else if time1.After(now) {
	// 	fmt.Println("voucher validity is there")
	// } else {
	// 	fmt.Println("time1 is equal to the current time")
	// }

	if(voucher.MinCartValue > *cartvalue){
		return voucher.MinCartValue,*cartvalue,0, ErrMinimumAmt
	}
	discount_amt := *cartvalue * voucher.Percentage / 100
	if(discount_amt > voucher.MaxDiscountAmt){
		return voucher.MinCartValue,*cartvalue-voucher.MaxDiscountAmt,voucher.MaxDiscountAmt, nil
	}
	return voucher.MinCartValue,*cartvalue-discount_amt,discount_amt, nil
}


