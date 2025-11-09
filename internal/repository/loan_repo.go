package repository

import (
    "context"
    "database/sql"
    "fmt"
    "time"
    
    "github.com/autosave/backend/internal/models"
    "github.com/jmoiron/sqlx"
)

type loanRepository struct {
    db *sqlx.DB
}

func NewLoanRepository(db *sqlx.DB) LoanRepository {
    return &loanRepository{db: db}
}

func (r *loanRepository) Create(ctx context.Context, loan *models.Loan) error {
    query := `
        INSERT INTO loans (
            user_id, name, original_debt, current_debt, rate,
            monthly_payment, autopay_enabled, autopay_bank_id,
            autopay_day, status, next_payment_date
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id, created_at, updated_at`
    
    loan.OriginalDebt = loan.CurrentDebt
    loan.Status = "active"
    
    err := r.db.QueryRowxContext(ctx, query,
        loan.UserID, loan.Name, loan.OriginalDebt, loan.CurrentDebt,
        loan.Rate, loan.MonthlyPayment, loan.AutopayEnabled,
        loan.AutopayBankID, loan.AutopayDay, loan.Status, loan.NextPaymentDate,
    ).Scan(&loan.ID, &loan.CreatedAt, &loan.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create loan: %w", err)
    }
    
    return nil
}

func (r *loanRepository) GetByID(ctx context.Context, id int) (*models.Loan, error) {
    var loan models.Loan
    query := `
        SELECT l.*, b.name as autopay_bank_name
        FROM loans l
        LEFT JOIN banks b ON l.autopay_bank_id = b.id
        WHERE l.id = $1`
    
    err := r.db.GetContext(ctx, &loan, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("loan not found")
        }
        return nil, fmt.Errorf("failed to get loan: %w", err)
    }
    
    return &loan, nil
}

func (r *loanRepository) GetUserLoans(ctx context.Context, userID int) ([]models.Loan, error) {
    var loans []models.Loan
    query := `
        SELECT l.*, b.name as autopay_bank_name
        FROM loans l
        LEFT JOIN banks b ON l.autopay_bank_id = b.id
        WHERE l.user_id = $1
        ORDER BY l.created_at DESC`
    
    err := r.db.SelectContext(ctx, &loans, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user loans: %w", err)
    }
    
    return loans, nil
}

func (r *loanRepository) GetActiveLoans(ctx context.Context, userID int) ([]models.Loan, error) {
    var loans []models.Loan
    query := `
        SELECT l.*, b.name as autopay_bank_name
        FROM loans l
        LEFT JOIN banks b ON l.autopay_bank_id = b.id
        WHERE l.user_id = $1 AND l.status = 'active'
        ORDER BY l.next_payment_date`
    
    err := r.db.SelectContext(ctx, &loans, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get active loans: %w", err)
    }
    
    return loans, nil
}

func (r *loanRepository) Update(ctx context.Context, loan *models.Loan) error {
    query := `
        UPDATE loans 
        SET name = $2, monthly_payment = $3, autopay_enabled = $4,
            autopay_bank_id = $5, autopay_day = $6, next_payment_date = $7
        WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query,
        loan.ID, loan.Name, loan.MonthlyPayment, loan.AutopayEnabled,
        loan.AutopayBankID, loan.AutopayDay, loan.NextPaymentDate,
    )
    
    if err != nil {
        return fmt.Errorf("failed to update loan: %w", err)
    }
    
    return nil
}

func (r *loanRepository) UpdateDebt(ctx context.Context, id int, currentDebt float64) error {
    query := `UPDATE loans SET current_debt = $2 WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query, id, currentDebt)
    if err != nil {
        return fmt.Errorf("failed to update debt: %w", err)
    }
    
    return nil
}

func (r *loanRepository) UpdateStatus(ctx context.Context, id int, status string) error {
    var query string
    if status == "paid_off" {
        query = `UPDATE loans SET status = $2, paid_off_at = NOW() WHERE id = $1`
    } else {
        query = `UPDATE loans SET status = $2 WHERE id = $1`
    }
    
    _, err := r.db.ExecContext(ctx, query, id, status)
    if err != nil {
        return fmt.Errorf("failed to update status: %w", err)
    }
    
    return nil
}

func (r *loanRepository) Delete(ctx context.Context, id int) error {
    query := `UPDATE loans SET status = 'cancelled' WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete loan: %w", err)
    }
    
    return nil
}

func (r *loanRepository) CreatePayment(ctx context.Context, payment *models.LoanPayment) error {
    query := `
        INSERT INTO loan_payments (
            loan_id, user_id, amount, is_autopay, bank_payment_id,
            status, scheduled_date, completed_at, error
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        RETURNING id, created_at`
    
    err := r.db.QueryRowxContext(ctx, query,
        payment.LoanID, payment.UserID, payment.Amount, payment.IsAutopay,
        payment.BankPaymentID, payment.Status, payment.ScheduledDate,
        payment.CompletedAt, payment.Error,
    ).Scan(&payment.ID, &payment.CreatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create payment: %w", err)
    }
    
    return nil
}

func (r *loanRepository) GetPayments(ctx context.Context, loanID int) ([]models.LoanPayment, error) {
    var payments []models.LoanPayment
    query := `
        SELECT * FROM loan_payments 
        WHERE loan_id = $1 
        ORDER BY scheduled_date DESC`
    
    err := r.db.SelectContext(ctx, &payments, query, loanID)
    if err != nil {
        return nil, fmt.Errorf("failed to get payments: %w", err)
    }
    
    return payments, nil
}

func (r *loanRepository) GetScheduledPayments(ctx context.Context, date time.Time) ([]models.LoanPayment, error) {
    var payments []models.LoanPayment
    query := `
        SELECT * FROM loan_payments 
        WHERE status = 'scheduled' AND scheduled_date <= $1`
    
    err := r.db.SelectContext(ctx, &payments, query, date)
    if err != nil {
        return nil, fmt.Errorf("failed to get scheduled payments: %w", err)
    }
    
    return payments, nil
}

func (r *loanRepository) UpdatePayment(ctx context.Context, payment *models.LoanPayment) error {
    query := `
        UPDATE loan_payments 
        SET bank_payment_id = $2, status = $3, completed_at = $4, error = $5
        WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query,
        payment.ID, payment.BankPaymentID, payment.Status,
        payment.CompletedAt, payment.Error,
    )
    
    if err != nil {
        return fmt.Errorf("failed to update payment: %w", err)
    }
    
    return nil
}
