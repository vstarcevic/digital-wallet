package database

import (
	"context"
	"database/sql"
	"errors"

	m "user-service/model"
)

var ErrDuplicate = errors.New("user already exists")

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
