package repository

import (
    "context"
    "database/sql"
    "fmt"
    
    "github.com/KotovBoris/AutoSave/backend/internal/models"
    "github.com/jmoiron/sqlx"
)

type goalRepository struct {
    db *sqlx.DB
}

func NewGoalRepository(db *sqlx.DB) GoalRepository {
    return &goalRepository{db: db}
}

func (r *goalRepository) Create(ctx context.Context, goal *models.Goal) error {
    query := `
        INSERT INTO goals (
            user_id, name, target_amount, current_amount,
            monthly_amount, bank_id, deposit_rate, position,
            status, next_deposit_date
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id, created_at, updated_at`
    
    err := r.db.QueryRowxContext(ctx, query,
        goal.UserID, goal.Name, goal.TargetAmount, goal.CurrentAmount,
        goal.MonthlyAmount, goal.BankID, goal.DepositRate, goal.Position,
        goal.Status, goal.NextDepositDate,
    ).Scan(&goal.ID, &goal.CreatedAt, &goal.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create goal: %w", err)
    }
    
    return nil
}

func (r *goalRepository) GetByID(ctx context.Context, id int) (*models.Goal, error) {
    var goal models.Goal
    query := `
        SELECT g.*, b.name as bank_name
        FROM goals g
        JOIN banks b ON g.bank_id = b.id
        WHERE g.id = $1`
    
    err := r.db.GetContext(ctx, &goal, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("goal not found")
        }
        return nil, fmt.Errorf("failed to get goal: %w", err)
    }
    
    return &goal, nil
}

func (r *goalRepository) GetUserGoals(ctx context.Context, userID int) ([]models.Goal, error) {
    var goals []models.Goal
    query := `
        SELECT g.*, b.name as bank_name
        FROM goals g
        JOIN banks b ON g.bank_id = b.id
        WHERE g.user_id = $1 AND g.status != 'cancelled'
        ORDER BY g.position`
    
    err := r.db.SelectContext(ctx, &goals, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user goals: %w", err)
    }
    
    return goals, nil
}

func (r *goalRepository) GetActiveGoal(ctx context.Context, userID int) (*models.Goal, error) {
    var goal models.Goal
    query := `
        SELECT g.*, b.name as bank_name
        FROM goals g
        JOIN banks b ON g.bank_id = b.id
        WHERE g.user_id = $1 AND g.status = 'active'
        ORDER BY g.position
        LIMIT 1`
    
    err := r.db.GetContext(ctx, &goal, query, userID)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get active goal: %w", err)
    }
    
    return &goal, nil
}

func (r *goalRepository) Update(ctx context.Context, goal *models.Goal) error {
    query := `
        UPDATE goals 
        SET name = $2, monthly_amount = $3, current_amount = $4,
            status = $5, next_deposit_date = $6, completed_at = $7
        WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query,
        goal.ID, goal.Name, goal.MonthlyAmount, goal.CurrentAmount,
        goal.Status, goal.NextDepositDate, goal.CompletedAt,
    )
    
    if err != nil {
        return fmt.Errorf("failed to update goal: %w", err)
    }
    
    return nil
}

func (r *goalRepository) UpdatePosition(ctx context.Context, goalID int, position int) error {
    query := `UPDATE goals SET position = $2 WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query, goalID, position)
    if err != nil {
        return fmt.Errorf("failed to update position: %w", err)
    }
    
    return nil
}

func (r *goalRepository) UpdateStatus(ctx context.Context, goalID int, status string) error {
    query := `UPDATE goals SET status = $2 WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query, goalID, status)
    if err != nil {
        return fmt.Errorf("failed to update status: %w", err)
    }
    
    return nil
}

func (r *goalRepository) UpdateCurrentAmount(ctx context.Context, goalID int, amount float64) error {
    query := `UPDATE goals SET current_amount = $2 WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query, goalID, amount)
    if err != nil {
        return fmt.Errorf("failed to update current amount: %w", err)
    }
    
    return nil
}

func (r *goalRepository) Delete(ctx context.Context, id int) error {
    query := `UPDATE goals SET status = 'cancelled' WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete goal: %w", err)
    }
    
    return nil
}

func (r *goalRepository) GetMaxPosition(ctx context.Context, userID int) (int, error) {
    var maxPos sql.NullInt64
    query := `SELECT MAX(position) FROM goals WHERE user_id = $1 AND status != 'cancelled'`
    
    err := r.db.GetContext(ctx, &maxPos, query, userID)
    if err != nil {
        return 0, fmt.Errorf("failed to get max position: %w", err)
    }
    
    if maxPos.Valid {
        return int(maxPos.Int64), nil
    }
    return 0, nil
}

