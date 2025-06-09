package services

import (
    "errors"
    "banking_service_project/models"
    "banking_service_project/repositories"
)

type AccountService interface {
    CreateAccount(userID int) (*models.Account, error)
    GetUserAccounts(userID int) ([]models.Account, error)
    Deposit(accountID int, amount float64) error
    Withdraw(accountID int, amount float64) error
    GetAccountByID(accountID int) (*models.Account, error)
}

type accountService struct {
    accountRepo     repositories.AccountRepository
    transactionRepo repositories.TransactionRepository
}

func NewAccountService(accountRepo repositories.AccountRepository, transactionRepo repositories.TransactionRepository) AccountService {
    return &accountService{accountRepo: accountRepo, transactionRepo: transactionRepo}
}

func (s *accountService) CreateAccount(userID int) (*models.Account, error) {
    account := &models.Account{
        UserID:  userID,
        Balance: 0,
    }
    if err := s.accountRepo.Create(account); err != nil {
        return nil, err
    }
    return account, nil
}

func (s *accountService) GetUserAccounts(userID int) ([]models.Account, error) {
    return s.accountRepo.GetByUserID(userID)
}

func (s *accountService) Deposit(accountID int, amount float64) error {
    acc, err := s.accountRepo.GetByID(accountID)
    if err != nil {
        return err
    }
    newBalance := acc.Balance + amount
    return s.accountRepo.UpdateBalance(accountID, newBalance)
}

func (s *accountService) Withdraw(accountID int, amount float64) error {
    acc, err := s.accountRepo.GetByID(accountID)
    if err != nil {
        return err
    }
    if acc.Balance < amount {
        return errors.New("insufficient funds")
    }
    newBalance := acc.Balance - amount
    return s.accountRepo.UpdateBalance(accountID, newBalance)
}

func (s *accountService) GetAccountByID(accountID int) (*models.Account, error) {
    return s.accountRepo.GetByID(accountID)
}
