package mainPage

import (
	"github.com/gorilla/mux"
)

func SetupMainPageRoutes(r *mux.Router) {
	r.HandleFunc("/", mainPageHandler)
}
