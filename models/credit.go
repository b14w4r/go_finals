package models

import "time"

type Credit struct {
    ID            int       `json:"id"`
    AccountID     int       `json:"account_id"`
    Principal     float64   `json:"principal"`
    InterestRate  float64   `json:"interest_rate"`
    TermMonths    int       `json:"term_months"`
    CreatedAt     time.Time `json:"created_at"`
}
