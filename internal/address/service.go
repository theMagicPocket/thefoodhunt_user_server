package address

import (
	"context"
	"errors"

	"github.com/deVamshi/golang_food_delivery_api/internal/entity"
	"github.com/deVamshi/golang_food_delivery_api/pkg/utils"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	Add(ctx context.Context, userId string, adrs *AddOrUpdateAddressRequest) (*entity.User, error)
	Update(ctx context.Context, userId string, adrs *AddOrUpdateAddressRequest) (*entity.User, error)
	Delete(ctx context.Context, userId string, addressId string) (bool, error)
}

type service struct {
	repo      Repository
	validator *validator.Validate
}

func NewService(repo Repository, v *validator.Validate) Service {
	return &service{repo: repo, validator: v}
}

type AddOrUpdateAddressRequest struct {
	ID          string `json:"id,omitempty" validate:"required"`
	Lat         string `json:"latitude" validate:"required"`
	Long        string `json:"longitude" validate:"required"`
	Street      string `json:"street" validate:"required"`
	DoorNo      string `json:"doorno" validate:"required"`
	Pincode     string `json:"pincode" validate:"required"`
	Name        string `json:"name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required,len=10,numeric"`
	Landmark    string `json:"landmark,omitempty" validate:"required"`
	IsActive    bool   `json:"is_active"  validate:"boolean"`
}

func (s *service) Add(ctx context.Context, id string, adrs *AddOrUpdateAddressRequest) (*entity.User, error) {

	adrs.ID = utils.GenerateID()
	err := s.validator.Struct(adrs)
	if err != nil {
		return nil, err
	}
	return s.repo.AddOrUpdate(ctx, id, entity.UserAddress(*adrs))
}

func (s *service) Update(ctx context.Context, id string, adrs *AddOrUpdateAddressRequest) (*entity.User, error) {

	// TODO: need to add validations for this, add binding field in the struct tag

	if adrs.ID == "" {
		return nil, errors.New("address id missing")
	}

	addrsToUpdate := entity.UserAddress(*adrs)

	return s.repo.AddOrUpdate(ctx, id, addrsToUpdate)
}

func (s *service) Delete(ctx context.Context, userId string, addressId string) (bool, error) {
	return s.repo.Delete(ctx, userId, addressId)
}
