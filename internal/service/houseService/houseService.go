package houseService

import (
	"context"
	"errors"
	"flatSellerAvito2024/internal/models"
	"flatSellerAvito2024/internal/storage"
	"fmt"
)

type HouseService struct {
	houseRepo storage.House
}

func NewHouseService(houseRepo storage.House) *HouseService {
	return &HouseService{houseRepo: houseRepo}
}

func (s *HouseService) Create(ctx context.Context, house *models.House) error {
	/*
		add data validation here
	*/

	err := s.houseRepo.Create(ctx, house)
	if err != nil {
		fmt.Println("something went wring in service: ", err.Error())
	}

	return err
}

func (s *HouseService) HouseInfo(ctx context.Context, id int) (*models.House, error) {
	if id < 0 {
		return nil, errors.New("invalid house ID")
	}

	house, err := s.houseRepo.HouseInfo(ctx, id)
	if err != nil {
		fmt.Println("something went wring in service: ", err.Error())
	}

	return house, nil

}
