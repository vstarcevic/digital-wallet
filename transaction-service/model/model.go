package model

import (
	"time"
)

type User struct {
	UserId    int       `json:"user-id"`
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
	UserId  *int    `json:"user-id"`
	Balance string  `json:"balance"`
	Error   *string `json:"error"`
}
