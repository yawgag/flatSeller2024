package middleware

import (
	"net/http"
)

/*
add tokens for check in 2 functions below
*/
func AuthChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("userRole")
		if err != nil {
			http.Error(w, "something went wrong. ", http.StatusForbidden)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func ModeratorChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roleCookie, err := r.Cookie("userRole")
		if err != nil {
			http.Error(w, "something went wrong", http.StatusForbidden)
		} else if roleCookie.Value == "user" {
			http.Error(w, "wrong user role", http.StatusForbidden)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
