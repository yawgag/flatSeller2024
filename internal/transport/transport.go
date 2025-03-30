package transport

import (
	"flatSellerAvito2024/internal/service"
	"flatSellerAvito2024/internal/transport/http/flatHandler"
	"flatSellerAvito2024/internal/transport/http/houseHandler"
	"flatSellerAvito2024/internal/transport/http/mainPage"
	"flatSellerAvito2024/internal/transport/http/user"

	"github.com/gorilla/mux"
)

type Handler struct {
	Services *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		Services: services,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	router := mux.NewRouter()

	userHandl := user.NewUserHandler(h.Services.User)
	userRouter := router.PathPrefix("/user").Subrouter()
	userRouter.HandleFunc("/register", userHandl.Register)
	userRouter.HandleFunc("/login", userHandl.Login)
	userRouter.HandleFunc("/logout", userHandl.Logout)

	houseHandl := houseHandler.NewHouseHandler(h.Services.House)
	houseRouter := router.PathPrefix("/house").Subrouter()
	houseRouter.Use(userHandl.LoginCheckMiddleware)
	houseRouter.Use(userHandl.RoleCheckMiddleware)
	houseRouter.HandleFunc("/{id:[0-9]+}", houseHandl.HouseInfo)
	houseRouter.HandleFunc("/create", houseHandl.Create)

	flatHandl := flatHandler.NewFlatHandler(h.Services.Flat)
	flatRouter := router.PathPrefix("/flat").Subrouter()
	flatRouter.Use(userHandl.LoginCheckMiddleware)
	flatRouter.HandleFunc("/create", flatHandl.CreateFlat)
	flatRouter.HandleFunc("/{id:[0-9]+}", flatHandl.FlatInfo)

	mainPageRouter := router.PathPrefix("/").Subrouter()
	mainPage.SetupMainPageRoutes(mainPageRouter)

	return router
}
