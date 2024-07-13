package model

import (
	"time"
)

type UserResponse struct {
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
	UserId  int    `json:"-"`
	Balance string `json:"balance"`
	Email   string `json:"email"`
	Error   string `json:"error"`
}
