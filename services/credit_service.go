package services

import (
    "errors"
    "math"
    "time"

    "banking_service_project/models"
    "banking_service_project/repositories"
)

type CreditService interface {
    ApplyCredit(accountID int, principal, annualRate float64, termMonths int) (*models.Credit, []models.PaymentSchedule, error)
    GetSchedule(creditID int) ([]models.PaymentSchedule, error)
}

type creditService struct {
    creditRepo    repositories.CreditRepository
    scheduleRepo  repositories.PaymentScheduleRepository
    accountRepo   repositories.AccountRepository
}

func NewCreditService(creditRepo repositories.CreditRepository, scheduleRepo repositories.PaymentScheduleRepository, accountRepo repositories.AccountRepository) CreditService {
    return &creditService{creditRepo: creditRepo, scheduleRepo: scheduleRepo, accountRepo: accountRepo}
}

func (s *creditService) ApplyCredit(accountID int, principal, annualRate float64, termMonths int) (*models.Credit, []models.PaymentSchedule, error) {
    acc, err := s.accountRepo.GetByID(accountID)
    if err != nil {
        return nil, nil, errors.New("account not found")
    }

    monthlyRate := annualRate / 12 / 100
    annuity := (principal * monthlyRate) / (1 - math.Pow(1+monthlyRate, float64(-termMonths)))

    credit := &models.Credit{
        AccountID:    acc.ID,
        Principal:    principal,
        InterestRate: annualRate,
        TermMonths:   termMonths,
        CreatedAt:    time.Now(),
    }
    if err := s.creditRepo.Create(credit); err != nil {
        return nil, nil, err
    }

    var schedules []models.PaymentSchedule
    for i := 1; i <= termMonths; i++ {
        dueDate := time.Now().AddDate(0, i, 0)
        schedule := models.PaymentSchedule{
            CreditID: credit.ID,
            DueDate:  dueDate,
            Amount:   annuity,
            Paid:     false,
        }
        if err := s.scheduleRepo.Create(&schedule); err != nil {
            return nil, nil, err
        }
        schedules = append(schedules, schedule)
    }
    return credit, schedules, nil
}

func (s *creditService) GetSchedule(creditID int) ([]models.PaymentSchedule, error) {
    return s.scheduleRepo.GetByCreditID(creditID)
}
