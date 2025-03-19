package userRepo

import (
	"context"

	"flatSellerAvito2024/internal/models"

	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) FindUserByLogin(ctx context.Context, login string) (*models.User, error) {
	queryCheckLogin := `select login, password, userRole
						from users
						where login = $1`

	user := &models.User{}

	err := r.pool.QueryRow(ctx, queryCheckLogin, login).Scan(&user.Login, &user.PasswordHash, &user.UserRole)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return user, nil

}

func (r *UserRepository) Register(ctx context.Context, user *models.User) error {
	queryAddNewUser := `insert into users(login, password, userRole)
						values ($1, $2, $3)`

	_, err := r.pool.Exec(ctx, queryAddNewUser, user.Login, user.PasswordHash, user.UserRole)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
