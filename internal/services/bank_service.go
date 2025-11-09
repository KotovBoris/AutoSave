package services

import (
    "context"
    "fmt"
    "time"
    
    "github.com/autosave/backend/internal/bankadapter"
    "github.com/autosave/backend/internal/banks"
    "github.com/autosave/backend/internal/models"
    "github.com/autosave/backend/internal/repository"
    "github.com/rs/zerolog"
)

type BankService struct {
    bankRepo       repository.BankRepository
    accountRepo    repository.AccountRepository
    transactionRepo repository.TransactionRepository
    bankFactory    *banks.Factory
    logger         *zerolog.Logger
}

func NewBankService(
    bankRepo repository.BankRepository,
    accountRepo repository.AccountRepository,
    transactionRepo repository.TransactionRepository,
    bankFactory *banks.Factory,
    logger *zerolog.Logger,
) *BankService {
    return &BankService{
        bankRepo:        bankRepo,
        accountRepo:     accountRepo,
        transactionRepo: transactionRepo,
        bankFactory:     bankFactory,
        logger:          logger,
    }
}

// GetAllBanks returns all available banks
func (s *BankService) GetAllBanks(ctx context.Context) ([]models.Bank, error) {
    return s.bankRepo.GetAll(ctx)
}

// ConnectBank connects user to a bank
func (s *BankService) ConnectBank(ctx context.Context, userID int, bankID string) (*models.BankConnection, error) {
    s.logger.Info().Int("userId", userID).Str("bankId", bankID).Msg("Connecting bank")
    
    // Check if already connected
    existing, _ := s.bankRepo.GetConnection(ctx, userID, bankID)
    if existing != nil && existing.Connected {
        return existing, nil
    }
    
    // Create bank adapter
    adapter, err := s.bankFactory.CreateAdapter(bankID)
    if err != nil {
        return nil, fmt.Errorf("failed to create bank adapter: %w", err)
    }
    
    // Get bank info
    bankInfo := adapter.GetBankInfo()
    
    // Get bank token
    tokenResp, err := adapter.GetBankToken(bankInfo.ID, "")
    if err != nil {
        s.logger.Error().Err(err).Str("bankId", bankID).Msg("Failed to get bank token")
        return nil, fmt.Errorf("failed to get bank token: %w", err)
    }
    
    // Format client ID
    externalClientID := fmt.Sprintf("team242-%d", userID)
    
    // Create account consent
    accountConsent, err := adapter.CreateAccountConsent(
        tokenResp.AccessToken,
        externalClientID,
        "team242",
        []string{"ReadAccountsDetail", "ReadBalances", "ReadTransactionsDetail"},
    )
    if err != nil {
        s.logger.Error().Err(err).Msg("Failed to create account consent")
        return nil, fmt.Errorf("failed to create account consent: %w", err)
    }
    
    // Create product consent
    productConsent, err := adapter.CreateProductConsent(
        tokenResp.AccessToken,
        externalClientID,
        "team242",
        []string{"read_product_agreements", "open_product_agreements", "close_product_agreements"},
    )
    if err != nil {
        s.logger.Error().Err(err).Msg("Failed to create product consent")
        return nil, fmt.Errorf("failed to create product consent: %w", err)
    }
    
    // Save connection
    connection := &models.BankConnection{
        UserID:             userID,
        BankID:             bankID,
        ExternalClientID:   externalClientID,
        BankToken:          tokenResp.AccessToken,
        TokenExpiresAt:     &tokenResp.IssuedAt,
        AccountConsentID:   &accountConsent.ConsentID,
        ProductConsentID:   &productConsent.ConsentID,
        Connected:          true,
    }
    
    // Update token expiry
    expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
    connection.TokenExpiresAt = &expiresAt
    
    if err := s.bankRepo.CreateConnection(ctx, connection); err != nil {
        s.logger.Error().Err(err).Msg("Failed to save bank connection")
        return nil, fmt.Errorf("failed to save connection: %w", err)
    }
    
    // Load accounts and transactions
    if err := s.syncBankData(ctx, connection, adapter); err != nil {
        s.logger.Warn().Err(err).Msg("Failed to sync bank data")
        // Don't fail the connection, just log the error
    }
    
    s.logger.Info().Int("connectionId", connection.ID).Msg("Bank connected successfully")
    
    return connection, nil
}

// GetConnectedBanks returns all connected banks for user
func (s *BankService) GetConnectedBanks(ctx context.Context, userID int) ([]models.BankConnection, error) {
    return s.bankRepo.GetUserConnections(ctx, userID)
}

// SyncBanks syncs data from all connected banks
func (s *BankService) SyncBanks(ctx context.Context, userID int) (*models.SyncBankResponse, error) {
    s.logger.Info().Int("userId", userID).Msg("Syncing banks")
    
    connections, err := s.bankRepo.GetUserConnections(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get connections: %w", err)
    }
    
    response := &models.SyncBankResponse{
        SyncedBanks: []string{},
        FailedBanks: []models.BankSyncError{},
        LastSyncAt:  time.Now(),
    }
    
    for _, conn := range connections {
        adapter, err := s.bankFactory.CreateAdapter(conn.BankID)
        if err != nil {
            response.FailedBanks = append(response.FailedBanks, models.BankSyncError{
                BankID: conn.BankID,
                Error:  err.Error(),
            })
            continue
        }
        
        if err := s.syncBankData(ctx, &conn, adapter); err != nil {
            response.FailedBanks = append(response.FailedBanks, models.BankSyncError{
                BankID: conn.BankID,
                Error:  err.Error(),
            })
        } else {
            response.SyncedBanks = append(response.SyncedBanks, conn.BankID)
        }
    }
    
    return response, nil
}

// syncBankData loads accounts and transactions from bank
func (s *BankService) syncBankData(ctx context.Context, conn *models.BankConnection, adapter bankadapter.BankAdapter) error {
    // Get accounts
    accounts, err := adapter.GetAccounts(
        conn.BankToken,
        conn.ExternalClientID,
        *conn.AccountConsentID,
        "team242",
    )
    if err != nil {
        return fmt.Errorf("failed to get accounts: %w", err)
    }
    
    // Save accounts
    for _, acc := range accounts {
        dbAccount := models.Account{
            UserID:         conn.UserID,
            UserBankID:     conn.ID,
            BankID:         conn.BankID,
            ExternalID:     acc.ID,
            Identification: acc.Identification,
            SchemeName:     &acc.SchemeName,
            AccountType:    &acc.AccountType,
            Nickname:       &acc.Nickname,
            Balance:        acc.Balance.Amount,
            Currency:       acc.Balance.Currency,
            ServicerName:   &acc.Servicer.Name,
            IsActive:       true,
        }
        
        // Check if account exists
        existing, _ := s.accountRepo.GetByExternalID(ctx, conn.ID, acc.ID)
        if existing != nil {
            // Update balance
            s.accountRepo.UpdateBalance(ctx, existing.ID, acc.Balance.Amount)
        } else {
            // Create new
            s.accountRepo.Create(ctx, &dbAccount)
        }
    }
    
    // Get transactions for each account (last 3 months)
    fromDate := time.Now().AddDate(0, -3, 0)
    toDate := time.Now()
    
    for _, acc := range accounts {
        transactions, err := adapter.GetTransactions(
            conn.BankToken,
            acc.ID,
            *conn.AccountConsentID,
            "team242",
            fromDate,
            toDate,
            100,
        )
        if err != nil {
            s.logger.Warn().Err(err).Str("accountId", acc.ID).Msg("Failed to get transactions")
            continue
        }
        
        // Get DB account
        dbAccount, err := s.accountRepo.GetByExternalID(ctx, conn.ID, acc.ID)
        if err != nil {
            continue
        }
        
        // Save transactions
        dbTransactions := make([]models.Transaction, 0, len(transactions))
        for _, tx := range transactions {
            dbTx := models.Transaction{
                AccountID:            dbAccount.ID,
                ExternalID:           tx.TransactionID,
                BookingDateTime:      tx.BookingDateTime,
                ValueDateTime:        &tx.ValueDateTime,
                Amount:               tx.Amount,
                Currency:             tx.Currency,
                Description:          &tx.TransactionInfo.Description,
                CreditDebitIndicator: &tx.CreditDebitIndicator,
                CounterpartyName:     &tx.CounterpartyName,
                CounterpartyAccount:  &tx.CounterpartyAccount,
                Category:             &tx.Category,
                IsSalary:             false,
            }
            dbTransactions = append(dbTransactions, dbTx)
        }
        
        if len(dbTransactions) > 0 {
            s.transactionRepo.CreateBatch(ctx, dbTransactions)
        }
    }
    
    // Update last sync time
    now := time.Now()
    conn.LastSyncAt = &now
    s.bankRepo.UpdateConnection(ctx, conn)
    
    return nil
}

// DisconnectBank disconnects user from bank
func (s *BankService) DisconnectBank(ctx context.Context, userID int, bankID string) error {
    s.logger.Info().Int("userId", userID).Str("bankId", bankID).Msg("Disconnecting bank")
    
    if err := s.bankRepo.DeleteConnection(ctx, userID, bankID); err != nil {
        return fmt.Errorf("failed to disconnect bank: %w", err)
    }
    
    return nil
}
