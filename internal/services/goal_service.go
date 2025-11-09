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

type GoalService struct {
    goalRepo    repository.GoalRepository
    depositRepo repository.DepositRepository
    userRepo    repository.UserRepository
    bankRepo    repository.BankRepository
    logger      *zerolog.Logger
}

func NewGoalService(
    goalRepo repository.GoalRepository,
    depositRepo repository.DepositRepository,
    userRepo repository.UserRepository,
    bankRepo repository.BankRepository,
    logger *zerolog.Logger,
) *GoalService {
    return &GoalService{
        goalRepo:    goalRepo,
        depositRepo: depositRepo,
        userRepo:    userRepo,
        bankRepo:    bankRepo,
        logger:      logger,
    }
}

// GetUserGoals returns all goals for user
func (s *GoalService) GetUserGoals(ctx context.Context, userID int) ([]models.GoalResponse, error) {
    goals, err := s.goalRepo.GetUserGoals(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get goals: %w", err)
    }
    
    response := make([]models.GoalResponse, 0, len(goals))
    
    for _, goal := range goals {
        // Get deposits
        deposits, _ := s.depositRepo.GetGoalDeposits(ctx, goal.ID)
        
        resp := s.buildGoalResponse(&goal, deposits)
        response = append(response, resp)
    }
    
    return response, nil
}

// CreateGoal creates new goal
func (s *GoalService) CreateGoal(ctx context.Context, userID int, req models.CreateGoalRequest) (*models.GoalResponse, error) {
    s.logger.Info().
        Int("userId", userID).
        Str("name", req.Name).
        Float64("targetAmount", req.TargetAmount).
        Msg("Creating goal")
    
    // Check if bank is connected
    conn, err := s.bankRepo.GetConnection(ctx, userID, req.BankID)
    if err != nil || conn == nil {
        return nil, fmt.Errorf("bank %s is not connected", req.BankID)
    }
    
    // Get bank info
    bank, err := s.bankRepo.GetByID(ctx, req.BankID)
    if err != nil {
        return nil, fmt.Errorf("bank not found: %w", err)
    }
    
    // Get user
    user, err := s.userRepo.GetByID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }
    
    // Get max position
    maxPos, _ := s.goalRepo.GetMaxPosition(ctx, userID)
    position := maxPos + 1
    
    // Determine status
    status := "waiting"
    var nextDepositDate *time.Time
    
    if position == 1 {
        status = "active"
        // Calculate next deposit date
        if len(user.SalaryDates) > 0 {
            next := calculateNextSalaryDate(user.SalaryDates)
            nextDepositDate = &next
        }
    }
    
    // Create goal
    goal := &models.Goal{
        UserID:          userID,
        Name:            req.Name,
        TargetAmount:    req.TargetAmount,
        CurrentAmount:   0,
        MonthlyAmount:   req.MonthlyAmount,
        BankID:          req.BankID,
        DepositRate:     bank.DepositRate,
        Position:        position,
        Status:          status,
        NextDepositDate: nextDepositDate,
    }
    
    if err := s.goalRepo.Create(ctx, goal); err != nil {
        return nil, fmt.Errorf("failed to create goal: %w", err)
    }
    
    // Build response with plan
    resp := s.buildGoalResponse(goal, []models.Deposit{})
    
    s.logger.Info().Int("goalId", goal.ID).Msg("Goal created successfully")
    
    return &resp, nil
}

// UpdateGoal updates goal
func (s *GoalService) UpdateGoal(ctx context.Context, userID, goalID int, req models.UpdateGoalRequest) error {
    s.logger.Info().Int("goalId", goalID).Msg("Updating goal")
    
    goal, err := s.goalRepo.GetByID(ctx, goalID)
    if err != nil {
        return fmt.Errorf("goal not found: %w", err)
    }
    
    if goal.UserID != userID {
        return fmt.Errorf("goal does not belong to user")
    }
    
    if req.Name != nil {
        goal.Name = *req.Name
    }
    
    if req.MonthlyAmount != nil {
        goal.MonthlyAmount = *req.MonthlyAmount
    }
    
    if err := s.goalRepo.Update(ctx, goal); err != nil {
        return fmt.Errorf("failed to update goal: %w", err)
    }
    
    return nil
}

// DeleteGoal deletes goal and closes all deposits
func (s *GoalService) DeleteGoal(ctx context.Context, userID, goalID int) error {
    s.logger.Info().Int("goalId", goalID).Msg("Deleting goal")
    
    goal, err := s.goalRepo.GetByID(ctx, goalID)
    if err != nil {
        return fmt.Errorf("goal not found: %w", err)
    }
    
    if goal.UserID != userID {
        return fmt.Errorf("goal does not belong to user")
    }
    
    // TODO: Close all deposits via bank API
    
    // Delete goal
    if err := s.goalRepo.Delete(ctx, goalID); err != nil {
        return fmt.Errorf("failed to delete goal: %w", err)
    }
    
    // Reorder remaining goals
    s.reorderGoalsAfterDelete(ctx, userID, goal.Position)
    
    return nil
}

// ReorderGoals changes goal priorities
func (s *GoalService) ReorderGoals(ctx context.Context, userID int, goalIDs []int) error {
    s.logger.Info().Int("userId", userID).Ints("goalIds", goalIDs).Msg("Reordering goals")
    
    // Verify all goals belong to user
    for i, goalID := range goalIDs {
        goal, err := s.goalRepo.GetByID(ctx, goalID)
        if err != nil {
            return fmt.Errorf("goal %d not found", goalID)
        }
        if goal.UserID != userID {
            return fmt.Errorf("goal %d does not belong to user", goalID)
        }
        
        // Update position
        newPosition := i + 1
        goal.Position = newPosition
        
        // Update status
        if newPosition == 1 {
            goal.Status = "active"
        } else if goal.Status == "active" {
            goal.Status = "waiting"
        }
        
        s.goalRepo.Update(ctx, goal)
    }
    
    return nil
}

// Helper functions

func (s *GoalService) buildGoalResponse(goal *models.Goal, deposits []models.Deposit) models.GoalResponse {
    // Calculate total interest
    totalInterest := 0.0
    for _, dep := range deposits {
        totalInterest += dep.AccruedInterest
    }
    
    // Calculate estimated completion
    var estimatedCompletion *time.Time
    monthsRemaining := 0
    
    if goal.MonthlyAmount > 0 {
        remaining := goal.TargetAmount - goal.CurrentAmount
        monthsRemaining = int(math.Ceil(remaining / goal.MonthlyAmount))
        
        if goal.NextDepositDate != nil {
            completion := goal.NextDepositDate.AddDate(0, monthsRemaining, 0)
            estimatedCompletion = &completion
        }
    }
    
    // Calculate progress
    progress := 0.0
    if goal.TargetAmount > 0 {
        progress = (goal.CurrentAmount / goal.TargetAmount) * 100
        if progress > 100 {
            progress = 100
        }
    }
    
    return models.GoalResponse{
        ID:                  goal.ID,
        Name:                goal.Name,
        TargetAmount:        goal.TargetAmount,
        CurrentAmount:       goal.CurrentAmount,
        MonthlyAmount:       goal.MonthlyAmount,
        BankID:              goal.BankID,
        BankName:            goal.BankName,
        DepositRate:         goal.DepositRate,
        Position:            goal.Position,
        Status:              goal.Status,
        NextDepositDate:     goal.NextDepositDate,
        CreatedAt:           goal.CreatedAt,
        CompletedAt:         goal.CompletedAt,
        Deposits:            deposits,
        EstimatedCompletion: estimatedCompletion,
        EstimatedInterest:   totalInterest,
        ProgressPercentage:  progress,
    }
}

func (s *GoalService) reorderGoalsAfterDelete(ctx context.Context, userID, deletedPosition int) {
    goals, _ := s.goalRepo.GetUserGoals(ctx, userID)
    
    for _, goal := range goals {
        if goal.Position > deletedPosition {
            goal.Position--
            s.goalRepo.UpdatePosition(ctx, goal.ID, goal.Position)
            
            if goal.Position == 1 {
                goal.Status = "active"
                s.goalRepo.UpdateStatus(ctx, goal.ID, "active")
            }
        }
    }
}

func calculateNextSalaryDate(salaryDates []int) time.Time {
    now := time.Now()
    year := now.Year()
    month := now.Month()
    
    // Find nearest salary date
    for _, day := range salaryDates {
        next := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
        if next.After(now) {
            return next
        }
    }
    
    // If no date this month, use next month
    nextMonth := month + 1
    nextYear := year
    if nextMonth > 12 {
        nextMonth = 1
        nextYear++
    }
    
    return time.Date(nextYear, nextMonth, salaryDates[0], 0, 0, 0, 0, time.UTC)
}
