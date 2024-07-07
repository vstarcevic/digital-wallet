package database

import (
	"database/sql"
	"errors"

	m "user-service/model"
)

var ErrDuplicate = errors.New("user already exists")

func CreateUser(conn *sql.DB, email string) (*m.UserResponse, error) {

	var user m.UserResponse
	var existingUser string

	// check if email already exists
	queryExists := `SELECT email FROM "user" WHERE email = $1`
	_ = conn.QueryRow(queryExists, email).Scan(&existingUser)

	if existingUser != "" {
		return nil, ErrDuplicate
	}

	query := `INSERT INTO "user" (email) VALUES ($1) returning id, email, created_at;`

	_ = conn.QueryRow(query, email).Scan(&user.UserId, &user.Email, &user.CreatedAt)

	return &user, nil
}
