package services

import (
    "time"

    "banking_service_project/repositories"
)

type AnalyticsService interface {
    GetMonthlyStats(userID int, month, year int) (map[string]float64, error)
    PredictBalance(accountID int, days int) (float64, error)
}

type analyticsService struct {
    transactionRepo repositories.TransactionRepository
}

func NewAnalyticsService(transactionRepo repositories.TransactionRepository) AnalyticsService {
    return &analyticsService{transactionRepo: transactionRepo}
}

func (s *analyticsService) GetMonthlyStats(userID int, month, year int) (map[string]float64, error) {
    // TODO: implement SQL aggregation for income/expenses
    return map[string]float64{
        "income":  0,
        "expense": 0,
    }, nil
}

func (s *analyticsService) PredictBalance(accountID int, days int) (float64, error) {
    // TODO: implement projection logic based on scheduled payments
    return 0, nil
}
