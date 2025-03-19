package flatRepo

import (
	"context"
	"errors"
	"flatSellerAvito2024/internal/models"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FlatRepository struct {
	pool *pgxpool.Pool
}

func NewFlatRepository(pool *pgxpool.Pool) *FlatRepository {
	return &FlatRepository{pool: pool}
}

func (r *FlatRepository) CreateWithTx(ctx context.Context, tx pgx.Tx, flat *models.Flat) error {
	queryCheckIfExist := `select Count(*) from flats
							where houseid = $1 and flatnumber = $2`

	var counter int
	err := tx.QueryRow(ctx, queryCheckIfExist, flat.HouseId, flat.FlatNumber).Scan(&counter)
	if err != nil {
		return err
	}
	if counter != 0 {
		return errors.New("flat with this number is already exist in this house")
	}

	query := `insert into flats(houseid, price, roomsnumber, flatnumber, moderationstatus)
			values ($1, $2, $3, $4, 1)
			returning id`

	err = tx.QueryRow(ctx, query, flat.HouseId, flat.Price, flat.RoomsNumber, flat.FlatNumber).Scan(&flat.Id)
	if err != nil {
		fmt.Println("somethnig went wrong with database. Id 3", err)
	}
	return err
}

func (r *FlatRepository) FlatInfo(ctx context.Context, id int) (*models.Flat, error) {
	query := `select houseid, price, roomsnumber, flatnumber, moderationstatus 
				from flats
				where id = $1`

	flat := &models.Flat{Id: id}

	err := r.pool.QueryRow(ctx, query, id).Scan(&flat.HouseId, &flat.Price, &flat.RoomsNumber, &flat.FlatNumber, &flat.ModerationStatus)

	if err != nil {
		fmt.Println("something wrong with database. Id 4 ", err)
		return nil, err
	}
	return flat, err
}

func (r *FlatRepository) IsExist(ctx context.Context, id int) (bool, error) {
	/*
		rewrite on "exist" in postgres query
	*/
	query := `select exists(select *
				from flats
				where id = $1)`

	var exist bool

	err := r.pool.QueryRow(ctx, query, id).Scan(&exist)
	if err != nil {
		return false, err
	}
	return exist, nil
}

// func (r *FlatRepository) UpdateWithTx(ctx context.Context, tx pgx.Tx, flat *models.Flat) error {
// 	queryNumberInHouse := `select exsists(select *
// 							from flats
// 							where houseid = $1 and flatnumber = $2)`

// 	queryUpdate := `update into flats
// 				where id = $1
// 				set houseId = $2,
// 					price = $3,
// 					roomsnumber = $4,
// 					flatnumber = $5,
// 					moderationstatus = $6`

// 	var flatIsExist bool
// 	err := tx.QueryRow(ctx, queryNumberInHouse, flat.HouseId, flat.FlatNumber).Scan(&flatIsExist)
// 	if err != nil {
// 		fmt.Println("Something went wrong in database. Id 5", err)
// 		return err
// 	}

// 	if flatIsExist {
// 		return errors.New("flat is already exist in this house")
// 	}

// 	_, err = r.pool.Exec(ctx, queryUpdate, flat.Id, flat.HouseId, flat.Price, flat.RoomsNumber, flat.FlatNumber, flat.ModerationStatus)
// 	if err != nil {
// 		fmt.Println("Something went wrong in database. Id 6", err)
// 		return err
// 	}
// 	return nil

// }
