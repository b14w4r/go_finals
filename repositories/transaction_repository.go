package repositories

import (
    "database/sql"
    "time"

    "banking_service_project/models"
)

type TransactionRepository interface {
    Create(tx *models.Transaction) error
    GetByAccountID(accountID int) ([]models.Transaction, error)
    GetByUserID(userID int) ([]models.Transaction, error)
}

type transactionRepository struct {
    db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
    return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(tx *models.Transaction) error {
    query := `INSERT INTO transactions (from_account_id, to_account_id, amount, type, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
    tx.CreatedAt = time.Now()
    err := r.db.QueryRow(query, tx.FromAccountID, tx.ToAccountID, tx.Amount, tx.Type, tx.CreatedAt).Scan(&tx.ID)
    if err != nil {
        return err
    }
    return nil
}

func (r *transactionRepository) GetByAccountID(accountID int) ([]models.Transaction, error) {
    query := `SELECT id, from_account_id, to_account_id, amount, type, created_at FROM transactions WHERE from_account_id=$1 OR to_account_id=$1`
    rows, err := r.db.Query(query, accountID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var transactions []models.Transaction
    for rows.Next() {
        var t models.Transaction
        if err := rows.Scan(&t.ID, &t.FromAccountID, &t.ToAccountID, &t.Amount, &t.Type, &t.CreatedAt); err != nil {
            return nil, err
        }
        transactions = append(transactions, t)
    }
    return transactions, nil
}

func (r *transactionRepository) GetByUserID(userID int) ([]models.Transaction, error) {
    // TODO: Implement joining accounts to filter by user
    return nil, nil
}
