package repositories

import (
    "database/sql"
    "time"

    "banking_service_project/models"
)

type PaymentScheduleRepository interface {
    Create(schedule *models.PaymentSchedule) error
    GetByCreditID(creditID int) ([]models.PaymentSchedule, error)
}

type paymentScheduleRepository struct {
    db *sql.DB
}

func NewPaymentScheduleRepository(db *sql.DB) PaymentScheduleRepository {
    return &paymentScheduleRepository{db: db}
}

func (r *paymentScheduleRepository) Create(schedule *models.PaymentSchedule) error {
    query := `INSERT INTO payment_schedules (credit_id, due_date, amount, paid) VALUES ($1, $2, $3, $4) RETURNING id`
    err := r.db.QueryRow(query, schedule.CreditID, schedule.DueDate, schedule.Amount, schedule.Paid).Scan(&schedule.ID)
    if err != nil {
        return err
    }
    return nil
}

func (r *paymentScheduleRepository) GetByCreditID(creditID int) ([]models.PaymentSchedule, error) {
    query := `SELECT id, credit_id, due_date, amount, paid FROM payment_schedules WHERE credit_id=$1`
    rows, err := r.db.Query(query, creditID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var schedules []models.PaymentSchedule
    for rows.Next() {
        var ps models.PaymentSchedule
        if err := rows.Scan(&ps.ID, &ps.CreditID, &ps.DueDate, &ps.Amount, &ps.Paid); err != nil {
            return nil, err
        }
        schedules = append(schedules, ps)
    }
    return schedules, nil
}
