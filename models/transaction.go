package models

import (
    "time"
)

type Transaction struct {
    ID            int       `json:"id"`
    FromAccountID int       `json:"from_account_id"`
    ToAccountID   int       `json:"to_account_id"`
    Amount        float64   `json:"amount"`
    CreatedAt     time.Time `json:"created_at"`
    Type          string    `json:"type"` // deposit, withdrawal, transfer
}
