package repository

import (
	"context"
	"errors"

	"github.com/ErenKarakus1/Notification-Platform/auth-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrEmailRegistered = errors.New("email is already registered")

var ErrUserNotFound = errors.New("user not found")

const createUserQuery = `
	INSERT INTO users (
		id,
		email,
		password_hash
	)
	VALUES ($1,$2,$3)
	RETURNING
		id,
		email,
		created_at
`

const getUserByEmailQuery = `
	SELECT
		id,
		email,
		password_hash,
		created_at
	FROM users
	WHERE email=$1
`

func CreateUser(ctx context.Context, pool *pgxpool.Pool, user model.User) (model.CreateUserResponse, error) {
	var createdUser model.CreateUserResponse
	err := pool.QueryRow(
		ctx,
		createUserQuery,
		user.ID,
		user.Email,
		user.PasswordHash,
	).Scan(
		&createdUser.ID,
		&createdUser.Email,
		&createdUser.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.CreateUserResponse{}, ErrEmailRegistered
		}
		return model.CreateUserResponse{}, errors.New("couldnt create user")
	}
	return createdUser, nil
}

func GetUserByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (model.User, error) {
	var user model.User
	err := pool.QueryRow(
		ctx,
		getUserByEmailQuery,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, errors.New("internal server error")
	}
	return user, nil
}
