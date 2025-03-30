package user

import (
	"flatSellerAvito2024/internal/models"
	"flatSellerAvito2024/internal/service"
	"fmt"
	"net/http"
	"time"
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
			Login:    r.FormValue("login"),
			Password: r.FormValue("password"),
		}

		tokens, err := h.userService.Login(r.Context(), user)
		if err != nil {
			fmt.Println(err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "refreshToken",
			Value:    tokens.RefreshToken,
			HttpOnly: true,
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			Path:     "/",
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "accessToken",
			Value:    tokens.AccessToken,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(7 * 24 * time.Hour),
			Path:     "/",
		})

	} else {
		http.ServeFile(w, r, "../pkg/template/login.html")
	}
}

func (h *userHandler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("refreshToken")
	if err != nil {
		fmt.Println(err)
		return
	}

	tokens := &models.Tokens{RefreshToken: refreshTokenCookie.Value}

	err = h.userService.Logout(r.Context(), tokens)
	if err != nil {
		fmt.Println(err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		HttpOnly: true,
		Expires:  time.Now().Add(-1 * time.Minute),
		Path:     "/",
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(-1 * time.Minute),
		Path:     "/",
	})

	fmt.Fprint(w, "logout")
}

func (h *userHandler) LoginCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenCookie, err := r.Cookie("accessToken")
		if err != nil {
			fmt.Println("need to login")
			return
		}
		refreshTokenCookie, err := r.Cookie("refreshToken")
		if err != nil {
			fmt.Println(err)
			return
		}

		tokens := &models.Tokens{
			AccessToken:  accessTokenCookie.Value,
			RefreshToken: refreshTokenCookie.Value,
		}

		tokens, err = h.userService.UserIsLogin(r.Context(), tokens)
		if err != nil {
			fmt.Println(err)
			return
		}

		if tokens.NewRefreshToken {
			http.SetCookie(w, &http.Cookie{
				Name:     "refreshToken",
				Value:    tokens.RefreshToken,
				HttpOnly: true,
				Expires:  time.Now().Add(7 * 24 * time.Hour),
				Path:     "/",
			})
		}

		if tokens.NewAccessToken {
			http.SetCookie(w, &http.Cookie{
				Name:     "accessToken",
				Value:    tokens.AccessToken,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				Expires:  time.Now().Add(7 * 24 * time.Hour),
				Path:     "/",
			})
		}

		next.ServeHTTP(w, r)
	})
}

func (h *userHandler) RoleCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessTokenCookie, err := r.Cookie("accessToken")
		if err != nil {
			fmt.Println("need to login")
			return
		}
		tokens := &models.Tokens{
			AccessToken: accessTokenCookie.Value,
		}

		avaliable, err := h.userService.AvaliableForUser(r.Context(), tokens)
		if err != nil {
			fmt.Println(err)
			return
		}

		if avaliable {
			next.ServeHTTP(w, r)
		}

	})
}
