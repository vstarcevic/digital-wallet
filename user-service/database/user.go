package database

import (
	"context"
	"database/sql"
	"errors"

	m "user-service/model"
)

var ErrDuplicate = errors.New("user already exists")
var ErrNotExist = errors.New("user does not exists")

func CreateUserWithTx(ctx context.Context, tx *sql.Tx, conn *sql.DB, email string) (*m.UserResponse, error) {

	var user m.UserResponse
	var existingUser string

	// check if email already exists
	queryExists := `SELECT email FROM "user" WHERE email = $1`
	_ = conn.QueryRow(queryExists, email).Scan(&existingUser)

	if existingUser != "" {
		return nil, ErrDuplicate
	}

	query := `INSERT INTO "user" (email) VALUES ($1) returning id, email, created_at;`

	err := tx.QueryRow(query, email).Scan(&user.UserId, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, errors.New("database error")
	}

	return &user, nil
}

func GetUserByEmail(conn *sql.DB, email string) (*m.UserResponse, error) {
	var user m.UserResponse

	query := `SELECT id, email, created_at FROM "user" WHERE email = $1`
	_ = conn.QueryRow(query, email).Scan(&user.UserId, &user.Email, &user.CreatedAt)

	if user.UserId == 0 {
		return nil, ErrNotExist
	}

	return &user, nil
}
