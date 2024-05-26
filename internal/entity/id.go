package entity

import "github.com/google/uuid"

func GenerateID() string {
	return uuid.New().String()
}
