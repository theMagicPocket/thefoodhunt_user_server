package utils

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
)

const earthRadiusKm = 6371

func toRadians(angle float64) float64 {
	return angle * math.Pi / 180
}

func HaversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	dLat := toRadians(lat2 - lat1)
	dLng := toRadians(lng2 - lng1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRadians(lat1))*math.Cos(toRadians(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

func GenerateID() string {
	length := 10
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// pass pointers to structs
func UpdateStruct(target interface{}, updates interface{}) error {
	targetVal := reflect.ValueOf(target).Elem()
	updatesVal := reflect.ValueOf(updates).Elem()

	if targetVal.Kind() != reflect.Struct || updatesVal.Kind() != reflect.Struct {
		return fmt.Errorf("both arguments must be pointers to structs")
	}

	for i := 0; i < updatesVal.NumField(); i++ {
		field := updatesVal.Type().Field(i)
		updateField := updatesVal.Field(i)

		if updateField.IsZero() {
			continue
		}

		targetField := targetVal.FieldByName(field.Name)
		if targetField.IsValid() && targetField.CanSet() {
			targetField.Set(updateField)
		}
	}
	return nil
}
