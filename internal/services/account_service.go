package services

import (
    "context"
    
    "github.com/autosave/backend/internal/models"
    "github.com/autosave/backend/internal/repository"
    "github.com/rs/zerolog"
)

type AccountService struct {
    accountRepo     repository.AccountRepository
    transactionRepo repository.TransactionRepository
    logger          *zerolog.Logger
}

func NewAccountService(
    accountRepo repository.AccountRepository,
    transactionRepo repository.TransactionRepository,
    logger *zerolog.Logger,
) *AccountService {
    return &AccountService{
        accountRepo:     accountRepo,
        transactionRepo: transactionRepo,
        logger:          logger,
    }
}

// GetUserAccounts returns all accounts for user
func (s *AccountService) GetUserAccounts(ctx context.Context, userID int) ([]models.AccountResponse, error) {
    accounts, err := s.accountRepo.GetUserAccounts(ctx, userID)
    if err != nil {
        return nil, err
    }
    
    response := make([]models.AccountResponse, 0, len(accounts))
    for _, acc := range accounts {
        response = append(response, acc.ToResponse())
    }
    
    return response, nil
}

// GetAccountTransactions returns transactions for account
func (s *AccountService) GetAccountTransactions(ctx context.Context, accountID int, limit int) ([]models.Transaction, error) {
    filter := models.TransactionFilter{
        AccountID: accountID,
        Limit:     limit,
    }
    
    return s.transactionRepo.GetAccountTransactions(ctx, filter)
}
