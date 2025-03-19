package transport

import (
	"flatSellerAvito2024/internal/service"
	"flatSellerAvito2024/internal/transport/http/flatHandler"
	"flatSellerAvito2024/internal/transport/http/houseHandler"
	"flatSellerAvito2024/internal/transport/http/mainPage"
	"flatSellerAvito2024/internal/transport/http/user"
	"flatSellerAvito2024/internal/transport/middleware"

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

	houseRouter := router.PathPrefix("/house").Subrouter()
	h.SetupHomeRoutes(houseRouter)

	flatRouter := router.PathPrefix("/flat").Subrouter()
	h.SetupFlatRouter(flatRouter)

	userRouter := router.PathPrefix("/user").Subrouter()
	h.SetupUserRoutes(userRouter)

	mainPageRouter := router.PathPrefix("/").Subrouter()
	mainPage.SetupMainPageRoutes(mainPageRouter)

	return router
}

func (h *Handler) SetupFlatRouter(r *mux.Router) {
	flatHandler := flatHandler.NewFlatHandler(h.Services.Flat)

	r.HandleFunc("/create", flatHandler.CreateFlat)
	r.HandleFunc("/{id:[0-9]+}", flatHandler.FlatInfo)
}

func (h *Handler) SetupHomeRoutes(r *mux.Router) {
	houseHandler := houseHandler.NewHouseHandler(h.Services.House)

	r.Use(middleware.ModeratorChecker)
	r.HandleFunc("/{id:[0-9]+}", houseHandler.HouseInfo)
	r.HandleFunc("/create", houseHandler.Create)

}

func (h *Handler) SetupUserRoutes(r *mux.Router) {
	userHandler := user.NewUserHandler(h.Services.User)

	r.HandleFunc("/register", userHandler.Register)
	r.HandleFunc("/login", userHandler.Login)
}
