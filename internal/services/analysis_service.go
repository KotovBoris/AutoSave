package services

import (
    "context"
    "fmt"
    "math"
    "time"
    
    "github.com/autosave/backend/internal/models"
    "github.com/autosave/backend/internal/repository"
    "github.com/rs/zerolog"
)

type AnalysisService struct {
    userRepo        repository.UserRepository
    transactionRepo repository.TransactionRepository
    logger          *zerolog.Logger
}

func NewAnalysisService(
    userRepo repository.UserRepository,
    transactionRepo repository.TransactionRepository,
    logger *zerolog.Logger,
) *AnalysisService {
    return &AnalysisService{
        userRepo:        userRepo,
        transactionRepo: transactionRepo,
        logger:          logger,
    }
}

// DetectSalaries analyzes transactions and detects salary patterns
func (s *AnalysisService) DetectSalaries(ctx context.Context, userID int) ([]models.SalaryDetection, error) {
    s.logger.Info().Int("userId", userID).Msg("Detecting salaries")
    
    // Get transactions for last 3 months
    fromDate := time.Now().AddDate(0, -3, 0)
    toDate := time.Now()
    
    transactions, err := s.transactionRepo.GetUserTransactions(ctx, userID, fromDate, toDate)
    if err != nil {
        return nil, fmt.Errorf("failed to get transactions: %w", err)
    }
    
    // Filter income transactions
    incomeTransactions := make([]models.Transaction, 0)
    for _, tx := range transactions {
        if tx.Amount > 0 {
            incomeTransactions = append(incomeTransactions, tx)
        }
    }
    
    if len(incomeTransactions) == 0 {
        return []models.SalaryDetection{}, nil
    }
    
    // Group by counterparty and find recurring patterns
    type incomePattern struct {
        counterparty string
        amounts      []float64
        dates        []time.Time
        transactions []models.Transaction
    }
    
    patterns := make(map[string]*incomePattern)
    
    for _, tx := range incomeTransactions {
        counterparty := "Unknown"
        if tx.CounterpartyName != nil {
            counterparty = *tx.CounterpartyName
        }
        
        if _, exists := patterns[counterparty]; !exists {
            patterns[counterparty] = &incomePattern{
                counterparty: counterparty,
                amounts:      []float64{},
                dates:        []time.Time{},
                transactions: []models.Transaction{},
            }
        }
        
        patterns[counterparty].amounts = append(patterns[counterparty].amounts, tx.Amount)
        patterns[counterparty].dates = append(patterns[counterparty].dates, tx.BookingDateTime)
        patterns[counterparty].transactions = append(patterns[counterparty].transactions, tx)
    }
    
    // Analyze patterns
    detections := []models.SalaryDetection{}
    
    for _, pattern := range patterns {
        // Need at least 2 transactions to detect pattern
        if len(pattern.transactions) < 2 {
            continue
        }
        
        // Calculate average amount
        avgAmount := 0.0
        for _, amt := range pattern.amounts {
            avgAmount += amt
        }
        avgAmount /= float64(len(pattern.amounts))
        
        // Check if amounts are similar (within 20% variance)
        variance := 0.0
        for _, amt := range pattern.amounts {
            diff := amt - avgAmount
            variance += diff * diff
        }
        variance /= float64(len(pattern.amounts))
        stdDev := math.Sqrt(variance)
        
        // If standard deviation is less than 20% of average, it's likely a salary
        confidence := "low"
        autoSelected := false
        
        if stdDev < avgAmount*0.2 && avgAmount > 10000 {
            confidence = "high"
            autoSelected = true
        } else if stdDev < avgAmount*0.4 && avgAmount > 5000 {
            confidence = "medium"
        }
        
        // Add all transactions from this pattern
        for _, tx := range pattern.transactions {
            detection := models.SalaryDetection{
                TransactionID: tx.ID,
                Date:          tx.BookingDateTime.Format("2006-01-02"),
                Amount:        tx.Amount,
                Counterparty:  pattern.counterparty,
                AccountID:     tx.AccountID,
                Confidence:    confidence,
                AutoSelected:  autoSelected,
            }
            detections = append(detections, detection)
        }
    }
    
    return detections, nil
}

// ConfirmSalaries marks transactions as salaries and calculates financial profile
func (s *AnalysisService) ConfirmSalaries(ctx context.Context, userID int, transactionIDs []int) (*models.SalaryAnalysis, error) {
    s.logger.Info().Int("userId", userID).Int("count", len(transactionIDs)).Msg("Confirming salaries")
    
    if len(transactionIDs) == 0 {
        return nil, fmt.Errorf("no transactions selected")
    }
    
    // Mark transactions as salary
    if err := s.transactionRepo.MarkAsSalary(ctx, transactionIDs); err != nil {
        return nil, fmt.Errorf("failed to mark as salary: %w", err)
    }
    
    // Calculate financial profile
    fromDate := time.Now().AddDate(0, -3, 0)
    toDate := time.Now()
    
    transactions, err := s.transactionRepo.GetUserTransactions(ctx, userID, fromDate, toDate)
    if err != nil {
        return nil, fmt.Errorf("failed to get transactions: %w", err)
    }
    
    totalIncome := 0.0
    totalExpenses := 0.0
    salaryDatesMap := make(map[int]bool)
    
    for _, tx := range transactions {
        if tx.Amount > 0 {
            totalIncome += tx.Amount
            if tx.IsSalary {
                day := tx.BookingDateTime.Day()
                salaryDatesMap[day] = true
            }
        } else {
            totalExpenses += math.Abs(tx.Amount)
        }
    }
    
    periodMonths := 3
    avgSalary := totalIncome / float64(periodMonths)
    avgExpenses := totalExpenses / float64(periodMonths)
    savingsCapacity := avgSalary - avgExpenses
    
    if savingsCapacity < 0 {
        savingsCapacity = 0
    }
    
    // Extract salary dates
    salaryDates := make([]int, 0, len(salaryDatesMap))
    for day := range salaryDatesMap {
        salaryDates = append(salaryDates, day)
    }
    
    // Update user profile
    if err := s.userRepo.UpdateFinancialProfile(ctx, userID, avgSalary, avgExpenses, savingsCapacity, salaryDates); err != nil {
        return nil, fmt.Errorf("failed to update profile: %w", err)
    }
    
    analysis := &models.SalaryAnalysis{
        AvgSalary:       avgSalary,
        AvgExpenses:     avgExpenses,
        SavingsCapacity: savingsCapacity,
        SalaryDates:     salaryDates,
        Analysis: models.AnalysisData{
            TotalIncome:   totalIncome,
            TotalExpenses: totalExpenses,
            PeriodMonths:  periodMonths,
        },
    }
    
    s.logger.Info().
        Int("userId", userID).
        Float64("avgSalary", avgSalary).
        Float64("savingsCapacity", savingsCapacity).
        Msg("Financial profile calculated")
    
    return analysis, nil
}
