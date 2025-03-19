package userService

import (
	"context"
	"errors"
	"flatSellerAvito2024/internal/models"
	"flatSellerAvito2024/internal/storage"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo      storage.User
	SessionsStore *sessions.CookieStore
}

func NewUserService(userRepo storage.User, sessionsStore *sessions.CookieStore) *UserService {
	return &UserService{
		userRepo:      userRepo,
		SessionsStore: sessionsStore,
	}
}

func (s *UserService) Register(ctx context.Context, user *models.User) error {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(passwordHash)

	err = s.userRepo.Register(ctx, user)
	return err
}

func (s *UserService) Login(ctx context.Context, userData *models.User, wr http.ResponseWriter, req *http.Request) error {

	userDb, err := s.userRepo.FindUserByLogin(ctx, userData.Login)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(userDb.PasswordHash), []byte(userData.PasswordHash))
	if err != nil {
		return errors.New("wrong password")
	}

	session, _ := s.SessionsStore.Get(req, "sessionId")
	session.Values["authenticated"] = true
	session.Values["userRole"] = userDb.UserRole
	err = session.Save(req, wr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("welcome, ", userDb.Login, "!")
	return nil
}
