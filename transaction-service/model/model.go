package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	UserId    int       `json:"user_id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type JsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type JsonRequest struct {
	Email string `json:"email"`
}

type UserBalanceResponse struct {
	UserId  *int            `json:"user_id"`
	Balance decimal.Decimal `json:"balance"`
	Error   *string         `json:"error"`
}

type AddMoneyRequest struct {
	UserId int             `json:"user_id"`
	Amount decimal.Decimal `json:"amount"`
}

type AddMoneyResponse struct {
	UpdatedBalance decimal.Decimal `json:"updated_balance"`
}

type TransferMoneyRequest struct {
	FromUserId int             `json:"from_user_id"`
	ToUserId   int             `json:"to_user_id"`
	Amount     decimal.Decimal `json:"amount_to_transfer"`
}

type TransferMoneyResponse struct{}
