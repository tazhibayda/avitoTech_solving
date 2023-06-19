package model

import "time"

type User struct {
	ID      int64   `json:"user_id" db:"user_id"`
	Balance float64 `json:"balance" db:"balance"`
}

type Transaction struct {
	TransactionId int       `json:"transaction_id"`
	UserId        int       `json:"user_id"`
	Amount        int       `json:"amount"`
	Operation     string    `json:"operation"`
	Date          time.Time `json:"date"`
}
