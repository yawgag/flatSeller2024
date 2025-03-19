package app

import (
	"flatSellerAvito2024/config"
	"flatSellerAvito2024/internal/service"
	"flatSellerAvito2024/internal/storage"
	"flatSellerAvito2024/internal/storage/postgres"
	"flatSellerAvito2024/internal/transport"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	dbConnPool, err := postgres.InitDb()
	if err != nil {
		log.Fatal("something wrong with database")
	}

	sessionsStore := sessions.NewCookieStore([]byte(cfg.SecretWord))
	sessionsStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,
	}

	repos := storage.NewRepositories(dbConnPool)
	services := service.NewServices(&service.Deps{
		Repos:         repos,
		SessionsStore: sessionsStore,
	})

	handler := transport.NewHandler(services)
	router := handler.InitRouter()

	http.ListenAndServe(cfg.ServerAddress, router)

}
