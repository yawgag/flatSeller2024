package user

import (
	"flatSellerAvito2024/internal/models"
	"flatSellerAvito2024/internal/service"
	"fmt"
	"net/http"
)

type userHandler struct {
	userService service.User
}

func NewUserHandler(userService service.User) *userHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "wrong data", http.StatusForbidden)
		}

		login := r.FormValue("login")
		password := r.FormValue("password")
		passwordRpt := r.FormValue("passwordRpt")
		userRole := r.FormValue("userRole")

		if password != passwordRpt {
			fmt.Println("password are different")
		}

		user := &models.User{
			Login:        login,
			PasswordHash: password,
			UserRole:     userRole,
		}

		h.userService.Register(r.Context(), user)

	} else {
		http.ServeFile(w, r, "../pkg/template/register.html")
	}
}

func (h *userHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			return
		}

		user := &models.User{
			Login:        r.FormValue("login"),
			PasswordHash: r.FormValue("password"),
		}

		err = h.userService.Login(r.Context(), user, w, r)
		if err != nil {
			fmt.Println(err)
			return
		}

	} else {
		http.ServeFile(w, r, "../pkg/template/login.html")
	}
}
