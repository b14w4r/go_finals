package repositories

import (
    "database/sql"
    "time"

    "banking_service_project/models"
)

type CreditRepository interface {
    Create(credit *models.Credit) error
    GetByID(creditID int) (*models.Credit, error)
}

type creditRepository struct {
    db *sql.DB
}

func NewCreditRepository(db *sql.DB) CreditRepository {
    return &creditRepository{db: db}
}

func (r *creditRepository) Create(credit *models.Credit) error {
    query := `INSERT INTO credits (account_id, principal, interest_rate, term_months, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
    credit.CreatedAt = time.Now()
    err := r.db.QueryRow(query, credit.AccountID, credit.Principal, credit.InterestRate, credit.TermMonths, credit.CreatedAt).Scan(&credit.ID)
    if err != nil {
        return err
    }
    return nil
}

func (r *creditRepository) GetByID(creditID int) (*models.Credit, error) {
    credit := &models.Credit{}
    query := `SELECT id, account_id, principal, interest_rate, term_months, created_at FROM credits WHERE id=$1`
    err := r.db.QueryRow(query, creditID).Scan(&credit.ID, &credit.AccountID, &credit.Principal, &credit.InterestRate, &credit.TermMonths, &credit.CreatedAt)
    if err != nil {
        return nil, err
    }
    return credit, nil
}
