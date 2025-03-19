package houseRepo

import (
	"context"
	"flatSellerAvito2024/internal/models"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HouseRepository struct {
	pool *pgxpool.Pool
}

func NewHouseRepository(pool *pgxpool.Pool) *HouseRepository {
	return &HouseRepository{pool: pool}
}

func (r *HouseRepository) Create(ctx context.Context, house *models.House) error {
	query := `INSERT INTO houses(address, buildyear, developer, dateofcreation)
				VALUES ($1, $2, $3, NOW())
				returning id;`
	err := r.pool.QueryRow(context.Background(), query, house.Address, house.BuildYear, house.Developer).Scan(&house.Id)

	if err != nil {
		fmt.Println("somethnig went wrong with database. Id 1")
	}

	return err
}

func (r *HouseRepository) HouseInfo(ctx context.Context, id int) (*models.House, error) {
	var house models.House

	query := `select * from houses
				where id = $1`

	err := r.pool.QueryRow(ctx, query, id).Scan(&house.Id, &house.Address, &house.BuildYear, &house.Developer, &house.DateOfCreation, &house.Lastflataddeddate)
	if err != nil {
		fmt.Println("somethnig went wrong with database. Id 2")
		return nil, err
	}

	return &house, nil
}

func (r *HouseRepository) ChangeLastFlatAddedDateWithTx(ctx context.Context, tx pgx.Tx, id int) error {
	query := `update houses
				set lastflataddeddate = now()
				where id = $1`

	_, err := tx.Exec(ctx, query, id)
	if err != nil {
		log.Println("something wrong with house value change")
	}

	return err
}

func (r *HouseRepository) IsExist(ctx context.Context, id int) (bool, error) {
	/*
		rewrite on "exist" in postgres query
	*/
	query := `select exists(select *
				from houses
				where id = $1)`

	var count int

	err := r.pool.QueryRow(ctx, query, id).Scan(&count)
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}
