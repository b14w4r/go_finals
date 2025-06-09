package repositories

import (
    "database/sql"
    "time"

    "banking_service_project/models"
)

type CardRepository interface {
    Create(card *models.Card) error
    GetByAccountID(accountID int) ([]models.Card, error)
}

type cardRepository struct {
    db *sql.DB
}

func NewCardRepository(db *sql.DB) CardRepository {
    return &cardRepository{db: db}
}

func (r *cardRepository) Create(card *models.Card) error {
    query := `INSERT INTO cards (account_id, encrypted_number, encrypted_cvv, expires_at, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
    card.CreatedAt = time.Now()
    err := r.db.QueryRow(query, card.AccountID, card.EncryptedNumber, card.EncryptedCVV, card.ExpiresAt, card.CreatedAt).Scan(&card.ID)
    if err != nil {
        return err
    }
    return nil
}

func (r *cardRepository) GetByAccountID(accountID int) ([]models.Card, error) {
    query := `SELECT id, account_id, encrypted_number, encrypted_cvv, expires_at, created_at FROM cards WHERE account_id=$1`
    rows, err := r.db.Query(query, accountID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var cards []models.Card
    for rows.Next() {
        var c models.Card
        if err := rows.Scan(&c.ID, &c.AccountID, &c.EncryptedNumber, &c.EncryptedCVV, &c.ExpiresAt, &c.CreatedAt); err != nil {
            return nil, err
        }
        cards = append(cards, c)
    }
    return cards, nil
}
