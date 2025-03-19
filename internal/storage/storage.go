package storage

import (
	"context"
	"flatSellerAvito2024/internal/models"
	"flatSellerAvito2024/internal/storage/postgres/flatRepo"
	"flatSellerAvito2024/internal/storage/postgres/houseRepo"
	"flatSellerAvito2024/internal/storage/postgres/transactionManager"
	"flatSellerAvito2024/internal/storage/postgres/userRepo"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type House interface {
	Create(ctx context.Context, house *models.House) error
	IsExist(ctx context.Context, id int) (bool, error)
	// Update(ctx context.Context, house *models.House)
	HouseInfo(ctx context.Context, id int) (*models.House, error)
	ChangeLastFlatAddedDateWithTx(ctx context.Context, tx pgx.Tx, id int) error
}

type Flat interface {
	CreateWithTx(ctx context.Context, tx pgx.Tx, flat *models.Flat) error
	IsExist(ctx context.Context, id int) (bool, error)
	// UpdateWithTx(ctx context.Context, tx pgx.Tx, flat *models.Flat) error
	FlatInfo(ctx context.Context, id int) (*models.Flat, error)
	// ChangeModerationStatus(ctx context.Context, flat *Flat)
}

type User interface {
	Register(ctx context.Context, user *models.User) error
	FindUserByLogin(ctx context.Context, login string) (*models.User, error)
}

type TransactionManager interface {
	TxBegin(ctx context.Context) (pgx.Tx, error)
	TxRollback(ctx context.Context, tx pgx.Tx) error
	TxCommit(ctx context.Context, tx pgx.Tx) error
}

type Repositories struct {
	Flat      Flat
	House     House
	User      User
	TxManager TransactionManager
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		Flat:      flatRepo.NewFlatRepository(db),
		House:     houseRepo.NewHouseRepository(db),
		User:      userRepo.NewUserRepository(db),
		TxManager: transactionManager.NewTransactionManager(db),
	}
}
