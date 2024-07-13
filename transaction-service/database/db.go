package database

import (
	"database/sql"
	"errors"
	"fmt"

	m "transaction-service/model"

	"github.com/shopspring/decimal"
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

func GetBalance(conn *sql.DB, userId int) string {

	var balance string

	queryExists := `select CAST(balance AS varchar) from balance where userid = $1`
	_ = conn.QueryRow(queryExists, userId).Scan(&balance)

	return balance
}

func TryUpdateBalanceWLock(tx *sql.Tx, userId int, amount decimal.Decimal) (*decimal.Decimal, error) {

	var currentBalance decimal.Decimal
	queryBalanceWLock := "SELECT balance FROM balance WHERE userid = $1 LIMIT 1 FOR NO KEY UPDATE;"
	tx.QueryRow(queryBalanceWLock, userId).Scan(&currentBalance)

	if currentBalance.Add(amount).LessThan(decimal.NewFromInt(0)) {
		return nil, errors.New("balance cannot be negative")
	}

	newAmount := currentBalance.Add(amount)

	var newBalance decimal.Decimal
	queryBalanceUpdate := "UPDATE balance set balance = $2 where userid = $1 returning balance;"
	err := tx.QueryRow(queryBalanceUpdate, userId, newAmount).Scan(&newBalance)
	if err != nil {
		return nil, err
	}

	if !newAmount.Equal(newBalance) {
		return nil, errors.New("error updating balance")
	}

	return &newAmount, nil
}

func AddTransaction(tx *sql.Tx, userId int, amount decimal.Decimal) error {

	query := `INSERT INTO "transaction" (userid, amount) VALUES ($1, $2) returning userid;`
	err := tx.QueryRow(query, userId, amount).Scan(&userId)
	if err != nil {
		return err
	}

	return nil
}
