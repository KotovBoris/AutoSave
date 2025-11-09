package repository

import (
    "context"
    "database/sql"
    "fmt"
    "time"
    
    "github.com/KotovBoris/AutoSave/backend/internal/models"
    "github.com/jmoiron/sqlx"
)

type depositRepository struct {
    db *sqlx.DB
}

func NewDepositRepository(db *sqlx.DB) DepositRepository {
    return &depositRepository{db: db}
}

func (r *depositRepository) Create(ctx context.Context, deposit *models.Deposit) error {
    query := `
        INSERT INTO deposits (
            goal_id, user_id, bank_id, product_id, agreement_id,
            amount, rate, term_months, status, opened_at,
            matures_at, accrued_interest, error
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        RETURNING id, created_at, updated_at`
    
    err := r.db.QueryRowxContext(ctx, query,
        deposit.GoalID, deposit.UserID, deposit.BankID, deposit.ProductID,
        deposit.AgreementID, deposit.Amount, deposit.Rate, deposit.TermMonths,
        deposit.Status, deposit.OpenedAt, deposit.MaturesAt,
        deposit.AccruedInterest, deposit.Error,
    ).Scan(&deposit.ID, &deposit.CreatedAt, &deposit.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create deposit: %w", err)
    }
    
    return nil
}

func (r *depositRepository) GetByID(ctx context.Context, id int) (*models.Deposit, error) {
    var deposit models.Deposit
    query := `SELECT * FROM deposits WHERE id = $1`
    
    err := r.db.GetContext(ctx, &deposit, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("deposit not found")
        }
        return nil, fmt.Errorf("failed to get deposit: %w", err)
    }
    
    return &deposit, nil
}

func (r *depositRepository) GetByAgreementID(ctx context.Context, agreementID string) (*models.Deposit, error) {
    var deposit models.Deposit
    query := `SELECT * FROM deposits WHERE agreement_id = $1`
    
    err := r.db.GetContext(ctx, &deposit, query, agreementID)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get deposit: %w", err)
    }
    
    return &deposit, nil
}

func (r *depositRepository) GetGoalDeposits(ctx context.Context, goalID int) ([]models.Deposit, error) {
    var deposits []models.Deposit
    query := `
        SELECT * FROM deposits 
        WHERE goal_id = $1 
        ORDER BY opened_at DESC`
    
    err := r.db.SelectContext(ctx, &deposits, query, goalID)
    if err != nil {
        return nil, fmt.Errorf("failed to get goal deposits: %w", err)
    }
    
    return deposits, nil
}

func (r *depositRepository) GetUserDeposits(ctx context.Context, userID int) ([]models.Deposit, error) {
    var deposits []models.Deposit
    query := `
        SELECT * FROM deposits 
        WHERE user_id = $1 
        ORDER BY opened_at DESC`
    
    err := r.db.SelectContext(ctx, &deposits, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user deposits: %w", err)
    }
    
    return deposits, nil
}

func (r *depositRepository) GetActiveDeposits(ctx context.Context, userID int) ([]models.Deposit, error) {
    var deposits []models.Deposit
    query := `
        SELECT * FROM deposits 
        WHERE user_id = $1 AND status = 'active'
        ORDER BY opened_at DESC`
    
    err := r.db.SelectContext(ctx, &deposits, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get active deposits: %w", err)
    }
    
    return deposits, nil
}

func (r *depositRepository) Update(ctx context.Context, deposit *models.Deposit) error {
    query := `
        UPDATE deposits 
        SET status = $2, accrued_interest = $3, error = $4
        WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query,
        deposit.ID, deposit.Status, deposit.AccruedInterest, deposit.Error,
    )
    
    if err != nil {
        return fmt.Errorf("failed to update deposit: %w", err)
    }
    
    return nil
}

func (r *depositRepository) UpdateStatus(ctx context.Context, id int, status string) error {
    query := `UPDATE deposits SET status = $2 WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query, id, status)
    if err != nil {
        return fmt.Errorf("failed to update status: %w", err)
    }
    
    return nil
}

func (r *depositRepository) Close(ctx context.Context, id int, closedAt time.Time) error {
    query := `UPDATE deposits SET status = 'closed', closed_at = $2 WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query, id, closedAt)
    if err != nil {
        return fmt.Errorf("failed to close deposit: %w", err)
    }
    
    return nil
}

