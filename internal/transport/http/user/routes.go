package user

import (
	"github.com/gorilla/mux"
)

func SetupUserRoutes(r *mux.Router) {
	r.HandleFunc("/dummyLogin", DummyLogin)
}
