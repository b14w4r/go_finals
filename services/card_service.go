package services

import (
    "errors"
    "time"

    "banking_service_project/models"
    "banking_service_project/repositories"
    "banking_service_project/utils"
)

type CardService interface {
    CreateCard(accountID int) (*models.Card, error)
    GetCards(accountID int) ([]models.Card, error)
}

type cardService struct {
    cardRepo     repositories.CardRepository
    accountRepo  repositories.AccountRepository
    pgpPublicKey string
    pgpPrivateKey string
}

func NewCardService(cardRepo repositories.CardRepository, accountRepo repositories.AccountRepository, pgpPublicKey, pgpPrivateKey string) CardService {
    return &cardService{cardRepo: cardRepo, accountRepo: accountRepo, pgpPublicKey: pgpPublicKey, pgpPrivateKey: pgpPrivateKey}
}

func (s *cardService) CreateCard(accountID int) (*models.Card, error) {
    acc, err := s.accountRepo.GetByID(accountID)
    if err != nil {
        return nil, errors.New("account not found")
    }

    cardNumber := utils.GenerateCardNumber()
    encryptedNumber, err := utils.EncryptPGP(cardNumber, s.pgpPublicKey)
    if err != nil {
        return nil, err
    }

    cvv := utils.GenerateCVV()
    encryptedCVV, err := utils.EncryptPGP(cvv, s.pgpPublicKey)
    if err != nil {
        return nil, err
    }

    card := &models.Card{
        AccountID:       acc.ID,
        EncryptedNumber: encryptedNumber,
        EncryptedCVV:    encryptedCVV,
        ExpiresAt:       time.Now().AddDate(3, 0, 0), // 3 years validity
    }

    if err := s.cardRepo.Create(card); err != nil {
        return nil, err
    }
    return card, nil
}

func (s *cardService) GetCards(accountID int) ([]models.Card, error) {
    return s.cardRepo.GetByAccountID(accountID)
}
