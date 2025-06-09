package models

import "time"

type Card struct {
    ID             int       `json:"id"`
    AccountID      int       `json:"account_id"`
    EncryptedNumber string   `json:"-"`
    EncryptedCVV   string    `json:"-"`
    ExpiresAt      time.Time `json:"expires_at"`
    CreatedAt      time.Time `json:"created_at"`
}
