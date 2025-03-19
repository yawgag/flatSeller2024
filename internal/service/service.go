package service

import (
	"context"
	"flatSellerAvito2024/internal/models"
	"flatSellerAvito2024/internal/service/flatService"
	"flatSellerAvito2024/internal/service/houseService"
	"flatSellerAvito2024/internal/service/userService"
	"flatSellerAvito2024/internal/storage"
	"net/http"

	"github.com/gorilla/sessions"
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

	/*
		Big kludge in functions below.
		ResponseWriter and request in service layer
	*/
	Login(ctx context.Context, userData *models.User, w http.ResponseWriter, r *http.Request) error
}

type Deps struct {
	Repos         *storage.Repositories
	SessionsStore *sessions.CookieStore
}

type Services struct {
	Flat  Flat
	House House
	User  User
}

func NewServices(deps *Deps) *Services {
	Flat := flatService.NewFlatService(deps.Repos.Flat, deps.Repos.House, deps.Repos.TxManager)
	House := houseService.NewHouseService(deps.Repos.House)
	User := userService.NewUserService(deps.Repos.User, deps.SessionsStore)

	return &Services{
		Flat:  Flat,
		House: House,
		User:  User,
	}
}
