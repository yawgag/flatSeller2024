package user

import (
	"net/http"
	"time"
)

func DummyLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "something went wrong", http.StatusBadRequest)
		}
		role := r.FormValue("role")
		userRoleCookie := http.Cookie{
			Name:    "userRole",
			Value:   role,
			Expires: time.Now().Add(time.Hour * 10),
			Path:    "/",
		}
		http.SetCookie(w, &userRoleCookie)
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		http.ServeFile(w, r, "../pkg/template/dummyLogin.html")
	}
}
