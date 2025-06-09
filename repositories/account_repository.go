package repositories

import (
    "database/sql"
    "errors"
    "time"

    "banking_service_project/models"
)

type AccountRepository interface {
    Create(account *models.Account) error
    GetByUserID(userID int) ([]models.Account, error)
    GetByID(accountID int) (*models.Account, error)
    UpdateBalance(accountID int, newBalance float64) error
}

type accountRepository struct {
    db *sql.DB
}

func NewAccountRepository(db *sql.DB) AccountRepository {
    return &accountRepository{db: db}
}

func (r *accountRepository) Create(account *models.Account) error {
    query := `INSERT INTO accounts (user_id, balance, created_at) VALUES ($1, $2, $3) RETURNING id`
    account.CreatedAt = time.Now()
    err := r.db.QueryRow(query, account.UserID, account.Balance, account.CreatedAt).Scan(&account.ID)
    if err != nil {
        return err
    }
    return nil
}

func (r *accountRepository) GetByUserID(userID int) ([]models.Account, error) {
    query := `SELECT id, user_id, balance, created_at FROM accounts WHERE user_id=$1`
    rows, err := r.db.Query(query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var accounts []models.Account
    for rows.Next() {
        var acc models.Account
        if err := rows.Scan(&acc.ID, &acc.UserID, &acc.Balance, &acc.CreatedAt); err != nil {
            return nil, err
        }
        accounts = append(accounts, acc)
    }
    return accounts, nil
}

func (r *accountRepository) GetByID(accountID int) (*models.Account, error) {
    account := &models.Account{}
    query := `SELECT id, user_id, balance, created_at FROM accounts WHERE id=$1`
    err := r.db.QueryRow(query, accountID).Scan(&account.ID, &account.UserID, &account.Balance, &account.CreatedAt)
    if err == sql.ErrNoRows {
        return nil, errors.New("account not found")
    }
    if err != nil {
        return nil, err
    }
    return account, nil
}

func (r *accountRepository) UpdateBalance(accountID int, newBalance float64) error {
    query := `UPDATE accounts SET balance=$1 WHERE id=$2`
    _, err := r.db.Exec(query, newBalance, accountID)
    return err
}
