package address

import (
	"context"

	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"github.com/deVamshi/golang_food_delivery_api/pkg/utils"
)

type Service interface {
	Add(ctx context.Context, userId string, adrs *AddOrUpdateAddressRequest) (*entity.User, error)
	Update(ctx context.Context, userId string, adrs *AddOrUpdateAddressRequest) (*entity.User, error)
	Delete(ctx context.Context, userId string, addressId string) (bool, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

type AddOrUpdateAddressRequest struct {
	ID          string `json:"id,omitempty"`
	Lat         string `json:"latitude" binding:"required"`
	Long        string `json:"longitude" binding:"required"`
	Street      string `json:"street" binding:"required"`
	DoorNo      string `json:"doorno" binding:"required"`
	Pincode     string `json:"pincode" binding:"required"`
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required,len=10,numeric"`
	Landmark    string `json:"landmark,omitempty" `
	IsActive    bool   `json:"is_active,omitempty"`
}

func (s *service) Add(ctx context.Context, id string, adrs *AddOrUpdateAddressRequest) (*entity.User, error) {

	adrs.ID = utils.GenerateID()
	addrsToAdd := entity.UserAddress(*adrs)

	return s.repo.AddOrUpdate(ctx, (id), addrsToAdd)
}

func (s *service) Update(ctx context.Context, id string, adrs *AddOrUpdateAddressRequest) (*entity.User, error) {

	addrsToUpdate := entity.UserAddress(*adrs)

	return s.repo.AddOrUpdate(ctx, (id), addrsToUpdate)
}

func (s *service) Delete(ctx context.Context, userId string, addressId string) (bool, error) {
	return s.repo.Delete(ctx, userId, addressId)
}
