package services

import (
    "errors"
    "banking_service_project/models"
    "banking_service_project/repositories"
)

type TransferService interface {
    Transfer(fromAccountID, toAccountID int, amount float64) (*models.Transaction, error)
}

type transferService struct {
    accountRepo     repositories.AccountRepository
    transactionRepo repositories.TransactionRepository
}

func NewTransferService(accountRepo repositories.AccountRepository, transactionRepo repositories.TransactionRepository) TransferService {
    return &transferService{accountRepo: accountRepo, transactionRepo: transactionRepo}
}

func (s *transferService) Transfer(fromAccountID, toAccountID int, amount float64) (*models.Transaction, error) {
    fromAcc, err := s.accountRepo.GetByID(fromAccountID)
    if err != nil {
        return nil, errors.New("from account not found")
    }
    if fromAcc.Balance < amount {
        return nil, errors.New("insufficient funds")
    }
    toAcc, err := s.accountRepo.GetByID(toAccountID)
    if err != nil {
        return nil, errors.New("to account not found")
    }

    // Perform debit and credit
    newFromBalance := fromAcc.Balance - amount
    newToBalance := toAcc.Balance + amount

    if err := s.accountRepo.UpdateBalance(fromAccountID, newFromBalance); err != nil {
        return nil, err
    }
    if err := s.accountRepo.UpdateBalance(toAccountID, newToBalance); err != nil {
        return nil, err
    }

    tx := &models.Transaction{
        FromAccountID: fromAccountID,
        ToAccountID:   toAccountID,
        Amount:        amount,
        Type:          "transfer",
    }
    if err := s.transactionRepo.Create(tx); err != nil {
        return nil, err
    }
    return tx, nil
}
