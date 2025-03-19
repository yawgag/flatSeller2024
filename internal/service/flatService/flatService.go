package flatService

import (
	"context"
	"errors"
	"flatSellerAvito2024/internal/models"
	"flatSellerAvito2024/internal/storage"
	"fmt"
	// "flatSellerAvito2024/internal/storage/postgres/flatRepo"
	// "flatSellerAvito2024/internal/storage/postgres/houseRepo"
)

type FlatService struct {
	flatRepository     storage.Flat
	houseRepository    storage.House
	transactionManager storage.TransactionManager
}

func NewFlatService(flatRepository storage.Flat, houseRepository storage.House, txManager storage.TransactionManager) *FlatService {
	return &FlatService{
		flatRepository:     flatRepository,
		houseRepository:    houseRepository,
		transactionManager: txManager,
	}
}

func (s *FlatService) Create(ctx context.Context, flat *models.Flat) error {

	houseIsExist, err := s.houseRepository.IsExist(ctx, flat.HouseId)
	if err != nil {
		return err
	} else if !houseIsExist {
		return errors.New("house is not exist")
	}

	tx, err := s.transactionManager.TxBegin(ctx)
	if err != nil {
		return err
	}
	// start transaction processing
	err = s.flatRepository.CreateWithTx(ctx, tx, flat)
	if err != nil {
		_ = s.transactionManager.TxRollback(ctx, tx)
		return err
	}

	err = s.houseRepository.ChangeLastFlatAddedDateWithTx(ctx, tx, flat.HouseId)
	if err != nil {
		err = s.transactionManager.TxRollback(ctx, tx)
		fmt.Println(err)
		return err
	}

	// end transaction processing

	err = s.transactionManager.TxCommit(ctx, tx)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *FlatService) FlatInfo(ctx context.Context, id int) (*models.Flat, error) {

	flatExist, err := s.flatRepository.IsExist(ctx, id)
	fmt.Println(err)
	fmt.Println(flatExist)

	flat, err := s.flatRepository.FlatInfo(ctx, id)
	if err != nil {
		fmt.Println("something wrong with db. Id 6")
		return nil, err
	}
	return flat, nil
}
