package flatHandler

import (
	"encoding/json"
	"flatSellerAvito2024/internal/models"
	"flatSellerAvito2024/internal/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type FlatHandler struct {
	flatService service.Flat
}

func NewFlatHandler(flatService service.Flat) *FlatHandler {
	return &FlatHandler{flatService: flatService}
}

func (h *FlatHandler) CreateFlat(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "error form", http.StatusBadRequest)
			return
		}

		houseId, err := strconv.Atoi(r.FormValue("houseId"))
		if err != nil {
			return
		}
		price, err := strconv.Atoi(r.FormValue("price"))
		if err != nil {
			return
		}
		roomsNumber, err := strconv.Atoi(r.FormValue("roomsNumber"))
		if err != nil {
			return
		}
		flatNumber, err := strconv.Atoi(r.FormValue("flatNumber"))
		if err != nil {
			return
		}

		flat := &models.Flat{
			HouseId:     houseId,
			Price:       price,
			RoomsNumber: roomsNumber,
			FlatNumber:  flatNumber,
		}

		err = h.flatService.Create(r.Context(), flat)
		if err != nil {
			fmt.Println("Flat Create error: ", err.Error())
			return
		}
		fmt.Fprint(w, "flat is in database now")

	} else {
		fmt.Println("working")
		http.ServeFile(w, r, "../pkg/template/createFlat.html")
	}
}

func (h *FlatHandler) FlatInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Println("wrong flat id value")
		return
	}

	flat, err := h.flatService.FlatInfo(r.Context(), id)
	if err != nil {
		fmt.Println("FlatInfo error: ", err.Error())
		return
	}

	flatJson, err := json.Marshal(flat)
	if err != nil {
		fmt.Println("something went wrong with json Marshal")
		return
	}

	w.Write(flatJson)

}

func (h *FlatHandler) AllFlatsInfo(w http.ResponseWriter, r *http.Request) {

}
