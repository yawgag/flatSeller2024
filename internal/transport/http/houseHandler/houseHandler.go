package houseHandler

import (
	"encoding/json"
	"flatSellerAvito2024/internal/models"
	"flatSellerAvito2024/internal/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type HouseHandler struct {
	houseService service.House
}

func NewHouseHandler(houseService service.House) *HouseHandler {
	return &HouseHandler{houseService: houseService}
}

func (h *HouseHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "error form", http.StatusBadRequest)
			return
		}

		house := &models.House{
			Address: r.FormValue("address"),
		}
		buildYear, err := strconv.Atoi(r.FormValue("buildYear"))
		if err != nil {
			http.Error(w, "error build year value", http.StatusBadRequest)
			return
		}
		developer := r.FormValue("developer")

		house.BuildYear = buildYear
		house.Developer = &developer

		if err := h.houseService.Create(r.Context(), house); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	} else {
		http.ServeFile(w, r, "../pkg/template/createHouse.html")
	}

}

func (h *HouseHandler) HouseInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	house, err := h.houseService.HouseInfo(r.Context(), id)
	if err != nil {
		fmt.Println("error: ", err.Error())
		return
	}

	houseJson, err := json.Marshal(house)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(houseJson)
}
