package userService

import (
	"context"
	"errors"
	"flatSellerAvito2024/config"
	"flatSellerAvito2024/internal/models"
	"flatSellerAvito2024/internal/storage"
	"fmt"
	"slices"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo storage.User
	cfg      *config.Config
}

func NewUserService(userRepo storage.User, cfg *config.Config) *UserService {
	return &UserService{
		userRepo: userRepo,
		cfg:      cfg,
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

func (s *UserService) Login(ctx context.Context, userData *models.User) (*models.Tokens, error) {

	user, err := s.userRepo.GetUser(ctx, userData.Login)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(userData.Password))
	if err != nil {
		fmt.Println("wrong password")
		return nil, err
	}

	sessionId := uuid.New().String()
	refreshTokenExpireTime, err := s.userRepo.CreateSession(ctx, user, sessionId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":     user.Id,
		"userRole":   user.UserRole,
		"expireTime": time.Now().Add(time.Second * 30), // changeTIME
	})

	signedAccessToken, err := accessToken.SignedString([]byte(s.cfg.SecretWord))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sessionId":  sessionId,
		"expireTime": refreshTokenExpireTime,
	})

	signedRefreshToken, err := refreshToken.SignedString([]byte(s.cfg.SecretWord))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &models.Tokens{
		AccessToken:  signedAccessToken,
		RefreshToken: signedRefreshToken,
	}, nil

}

func (s *UserService) Logout(ctx context.Context, tokens *models.Tokens) error {
	refreshTokenClaims, err := s.ParseJwt(tokens.RefreshToken)
	if err != nil {
		return err
	}

	err = s.userRepo.DeleteSession(ctx, (*refreshTokenClaims)["sessionId"].(string))
	return err
}

func (s *UserService) UserIsLogin(ctx context.Context, tokens *models.Tokens) (*models.Tokens, error) {
	accessTokenClaims, err := s.ParseJwt(tokens.AccessToken)
	if err != nil {
		return nil, err
	}
	refreshTokenClaims, err := s.ParseJwt(tokens.RefreshToken)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	ATexpireTime, err := time.Parse(time.RFC3339, (*accessTokenClaims)["expireTime"].(string))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if ATexpireTime.Before(time.Now()) {
		tokens.NewAccessToken = true
		fmt.Println("new access token")
		RTexpireTime, err := time.Parse(time.RFC3339, (*refreshTokenClaims)["expireTime"].(string))
		if err != nil {
			return nil, err
		}

		if RTexpireTime.Before(time.Now()) {
			tokens.NewRefreshToken = true
			newExpireTime, err := s.userRepo.UpdateSessionExpireTime(ctx, (*refreshTokenClaims)["sessionId"].(string))
			if err != nil {
				return nil, err
			}
			(*refreshTokenClaims)["expireTime"] = newExpireTime

			newRefreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(s.cfg.SecretWord))
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			tokens.RefreshToken = newRefreshToken

		}

		session, err := s.userRepo.GetSession(ctx, (*refreshTokenClaims)["sessionId"].(string))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		(*accessTokenClaims)["userRole"] = session.UserRole
		(*accessTokenClaims)["expireTime"] = time.Now().Add(time.Second * 30)

		newAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString([]byte(s.cfg.SecretWord))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		tokens.AccessToken = newAccessToken
	}

	return tokens, nil
}

func (s *UserService) AvaliableForUser(ctx context.Context, tokens *models.Tokens) (bool, error) {
	avaliableRoles := []string{"moderator", "admin"}

	accessTokenClaims, err := s.ParseJwt(tokens.AccessToken)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	userRole := (*accessTokenClaims)["userRole"].(string)

	if slices.Contains(avaliableRoles, userRole) {
		return true, nil
	}
	return false, nil

}

func (s *UserService) ParseJwt(token string) (*jwt.MapClaims, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("wrong token format")
		}
		return []byte(s.cfg.SecretWord), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if ok && jwtToken.Valid {
		return &claims, nil
	}

	return nil, errors.New("wrong token")
}
