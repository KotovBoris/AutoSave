package repository

import (
    "context"
    "database/sql"
    "fmt"
    
    "github.com/KotovBoris/AutoSave/backend/internal/models"
    "github.com/jmoiron/sqlx"
)

type operationRepository struct {
    db *sqlx.DB
}

func NewOperationRepository(db *sqlx.DB) OperationRepository {
    return &operationRepository{db: db}
}

func (r *operationRepository) Create(ctx context.Context, operation *models.Operation) error {
    query := `
        INSERT INTO operations (
            user_id, type, amount, related_goal_id, related_loan_id,
            related_deposit_id, status, error, metadata
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at`
    
    err := r.db.QueryRowxContext(ctx, query,
        operation.UserID, operation.Type, operation.Amount,
        operation.RelatedGoalID, operation.RelatedLoanID,
        operation.RelatedDepositID, operation.Status,
        operation.Error, operation.Metadata,
    ).Scan(&operation.ID, &operation.CreatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create operation: %w", err)
    }
    
    return nil
}

func (r *operationRepository) GetByID(ctx context.Context, id int) (*models.Operation, error) {
    var operation models.Operation
    query := `SELECT * FROM operations WHERE id = $1`
    
    err := r.db.GetContext(ctx, &operation, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("operation not found")
        }
        return nil, fmt.Errorf("failed to get operation: %w", err)
    }
    
    return &operation, nil
}

func (r *operationRepository) GetUserOperations(ctx context.Context, userID int, limit int) ([]models.Operation, error) {
    var operations []models.Operation
    query := `
        SELECT * FROM operations 
        WHERE user_id = $1 
        ORDER BY created_at DESC
        LIMIT $2`
    
    err := r.db.SelectContext(ctx, &operations, query, userID, limit)
    if err != nil {
        return nil, fmt.Errorf("failed to get user operations: %w", err)
    }
    
    return operations, nil
}

func (r *operationRepository) GetByType(ctx context.Context, userID int, operationType string) ([]models.Operation, error) {
    var operations []models.Operation
    query := `
        SELECT * FROM operations 
        WHERE user_id = $1 AND type = $2
        ORDER BY created_at DESC`
    
    err := r.db.SelectContext(ctx, &operations, query, userID, operationType)
    if err != nil {
        return nil, fmt.Errorf("failed to get operations by type: %w", err)
    }
    
    return operations, nil
}

