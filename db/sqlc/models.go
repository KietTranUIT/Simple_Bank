package db

import (
	"database/sql"
	"time"
)

type Account struct {
	ID        int64     `json:"id"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"createdAt"`
}

type Entry struct {
	ID        int64         `json:"id"`
	AccountID sql.NullInt64 `json:"accountId"`
	// can be negative or positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"createdAt"`
}

type Transfer struct {
	ID            int64         `json:"id"`
	FromAccountID sql.NullInt64 `json:"fromAccountId"`
	ToAccountID   sql.NullInt64 `json:"toAccountId"`
	// must be positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"createdAt"`
}
