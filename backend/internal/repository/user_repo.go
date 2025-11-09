package repository

import (
    "context"
    "database/sql"
    "fmt"
    
    "github.com/KotovBoris/AutoSave/backend/internal/models"
    "github.com/jmoiron/sqlx"
    "github.com/lib/pq"
)

type userRepository struct {
    db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
    return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
    query := `
        INSERT INTO users (email, password_hash, created_at, updated_at)
        VALUES ($1, $2, NOW(), NOW())
        RETURNING id, created_at, updated_at`
    
    err := r.db.QueryRowxContext(ctx, query, 
        user.Email, 
        user.PasswordHash,
    ).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
    
    if err != nil {
        if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
            return fmt.Errorf("email already exists")
        }
        return fmt.Errorf("failed to create user: %w", err)
    }
    
    return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
    var user models.User
    query := `
        SELECT id, email, password_hash, avg_salary, avg_expenses, 
               savings_capacity, salary_dates, autopilot_enabled,
               created_at, updated_at
        FROM users 
        WHERE id = $1`
    
    err := r.db.GetContext(ctx, &user, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
    var user models.User
    query := `
        SELECT id, email, password_hash, avg_salary, avg_expenses, 
               savings_capacity, salary_dates, autopilot_enabled,
               created_at, updated_at
        FROM users 
        WHERE email = $1`
    
    err := r.db.GetContext(ctx, &user, query, email)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("user not found")
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
    query := `
        UPDATE users 
        SET email = $2, avg_salary = $3, avg_expenses = $4,
            savings_capacity = $5, salary_dates = $6, 
            autopilot_enabled = $7, updated_at = NOW()
        WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query,
        user.ID, user.Email, user.AvgSalary, user.AvgExpenses,
        user.SavingsCapacity, pq.Array(user.SalaryDates),
        user.AutopilotEnabled,
    )
    
    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }
    
    return nil
}

func (r *userRepository) UpdateFinancialProfile(ctx context.Context, userID int, avgSalary, avgExpenses, savingsCapacity float64, salaryDates []int) error {
    query := `
        UPDATE users 
        SET avg_salary = $2, avg_expenses = $3, 
            savings_capacity = $4, salary_dates = $5,
            updated_at = NOW()
        WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query,
        userID, avgSalary, avgExpenses, savingsCapacity, pq.Array(salaryDates),
    )
    
    if err != nil {
        return fmt.Errorf("failed to update financial profile: %w", err)
    }
    
    return nil
}

func (r *userRepository) UpdateAutopilot(ctx context.Context, userID int, enabled bool) error {
    query := `UPDATE users SET autopilot_enabled = $2, updated_at = NOW() WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query, userID, enabled)
    if err != nil {
        return fmt.Errorf("failed to update autopilot: %w", err)
    }
    
    return nil
}

