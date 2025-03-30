package service

import (
	"context"
	"flatSellerAvito2024/config"
	"flatSellerAvito2024/internal/models"
	"flatSellerAvito2024/internal/service/flatService"
	"flatSellerAvito2024/internal/service/houseService"
	"flatSellerAvito2024/internal/service/userService"
	"flatSellerAvito2024/internal/storage"
)

type House interface {
	Create(ctx context.Context, house *models.House) error
	// Update(ctx context.Context, house *models.House)
	HouseInfo(ctx context.Context, id int) (*models.House, error)
}

type Flat interface {
	Create(ctx context.Context, flat *models.Flat) error
	// Update(ctx context.Context, flat *Flat)
	FlatInfo(ctx context.Context, id int) (*models.Flat, error)
	// ChangeModerationStatus(ctx context.Context, flat *Flat)
}

type User interface {
	Register(ctx context.Context, user *models.User) error

	Login(ctx context.Context, userData *models.User) (*models.Tokens, error)
	Logout(ctx context.Context, tokens *models.Tokens) error
	UserIsLogin(ctx context.Context, tokens *models.Tokens) (*models.Tokens, error)
	AvaliableForUser(ctx context.Context, tokens *models.Tokens) (bool, error)
}

type Deps struct {
	Repos *storage.Repositories
	Cfg   *config.Config
}

type Services struct {
	Flat  Flat
	House House
	User  User
}

func NewServices(deps *Deps) *Services {
	Flat := flatService.NewFlatService(deps.Repos.Flat, deps.Repos.House, deps.Repos.TxManager)
	House := houseService.NewHouseService(deps.Repos.House)
	User := userService.NewUserService(deps.Repos.User, deps.Cfg)

	return &Services{
		Flat:  Flat,
		House: House,
		User:  User,
	}
}
