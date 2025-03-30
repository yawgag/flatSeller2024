package userRepo

import (
	"context"
	"time"

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

// func (r *UserRepository) FindUserByLogin(ctx context.Context, login string) (*models.User, error) {
// 	queryCheckLogin := `select login, password, userRole
// 						from users
// 						where login = $1`

// 	user := &models.User{}

// 	err := r.pool.QueryRow(ctx, queryCheckLogin, login).Scan(&user.Login, &user.PasswordHash, &user.UserRole)
// 	if err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}

// 	return user, nil

// }

func (r *UserRepository) GetUser(ctx context.Context, login string) (*models.User, error) {
	query := `select id, login, password, userRole
				from users
				where login = $1`

	user := &models.User{}

	err := r.pool.QueryRow(ctx, query, login).Scan(&user.Id, &user.Login, &user.PasswordHash, &user.UserRole)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) CreateSession(ctx context.Context, user *models.User, sessionId string) (time.Time, error) {
	query := `insert into sessions(sessionId, userId, userRole, expireAt)
				values(
					$1, $2, $3, NOW() + INTERVAL '7 days'
				)
				returning expireAt`
	var expireAt time.Time
	err := r.pool.QueryRow(ctx, query, sessionId, user.Id, user.UserRole).Scan(&expireAt)

	return expireAt, err
}

func (r *UserRepository) GetSession(ctx context.Context, sessionId string) (*models.Session, error) {
	query := `select userId, userRole, expireAt
				from sessions
				where sessionId = $1`

	session := &models.Session{}

	err := r.pool.QueryRow(ctx, query, sessionId).Scan(&session.UserId, &session.UserRole, &session.ExpireAt)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *UserRepository) DeleteSession(ctx context.Context, sessionId string) error {
	query := `delete from sessions
				where sessionId = $1`

	_, err := r.pool.Exec(ctx, query, sessionId)
	return err
}

func (r *UserRepository) UpdateSessionExpireTime(ctx context.Context, sessionId string) (time.Time, error) {
	query := `update sessions
				set expireAt = NOW() + INTERVAL '7 days'
				where sessionId = $1
				returning expireAt`

	var expireAt time.Time
	err := r.pool.QueryRow(ctx, query, sessionId).Scan(&expireAt)

	return expireAt, err
}

/*
rewrite?
*/
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
