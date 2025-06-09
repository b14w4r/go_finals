package models

import "time"

type PaymentSchedule struct {
    ID        int       `json:"id"`
    CreditID  int       `json:"credit_id"`
    DueDate   time.Time `json:"due_date"`
    Amount    float64   `json:"amount"`
    Paid      bool      `json:"paid"`
}
