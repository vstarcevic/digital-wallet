package database

import (
	"database/sql"
	"errors"
	"fmt"

	m "transaction-service/model"
)

var ErrDuplicate = errors.New("user balance already exists")

func InsertBalance(conn *sql.DB, user m.User) error {

	var existingUserId int

	// check if balance already exists for the user
	queryExists := `SELECT 1 FROM "user" WHERE userId = $1`
	_ = conn.QueryRow(queryExists, user.UserId).Scan(&existingUserId)

	if existingUserId > 0 {
		return ErrDuplicate
	}

	query := `INSERT INTO "balance" (userId) VALUES ($1);`

	_, err := conn.Exec(query, user.UserId)
	if err != nil {
		return errors.New("database error")
	}
	fmt.Printf("User Balance added: %d", user.UserId)

	return nil
}
