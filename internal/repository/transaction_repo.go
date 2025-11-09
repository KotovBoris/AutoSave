package repository

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    "time"
    
    "github.com/autosave/backend/internal/models"
    "github.com/jmoiron/sqlx"
    "github.com/lib/pq"
)

type transactionRepository struct {
    db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
    return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *models.Transaction) error {
    query := `
        INSERT INTO transactions (
            account_id, external_id, booking_date_time, value_date_time,
            amount, currency, description, credit_debit_indicator,
            counterparty_name, counterparty_account, category, is_salary
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id, created_at`
    
    err := r.db.QueryRowxContext(ctx, query,
        tx.AccountID, tx.ExternalID, tx.BookingDateTime, tx.ValueDateTime,
        tx.Amount, tx.Currency, tx.Description, tx.CreditDebitIndicator,
        tx.CounterpartyName, tx.CounterpartyAccount, tx.Category, tx.IsSalary,
    ).Scan(&tx.ID, &tx.CreatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create transaction: %w", err)
    }
    
    return nil
}

func (r *transactionRepository) CreateBatch(ctx context.Context, transactions []models.Transaction) error {
    if len(transactions) == 0 {
        return nil
    }
    
    valueStrings := make([]string, 0, len(transactions))
    valueArgs := make([]interface{}, 0, len(transactions)*12)
    
    for i, tx := range transactions {
        valueStrings = append(valueStrings, fmt.Sprintf(
            "($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
            i*12+1, i*12+2, i*12+3, i*12+4, i*12+5, i*12+6,
            i*12+7, i*12+8, i*12+9, i*12+10, i*12+11, i*12+12,
        ))
        
        valueArgs = append(valueArgs,
            tx.AccountID, tx.ExternalID, tx.BookingDateTime, tx.ValueDateTime,
            tx.Amount, tx.Currency, tx.Description, tx.CreditDebitIndicator,
            tx.CounterpartyName, tx.CounterpartyAccount, tx.Category, tx.IsSalary,
        )
    }
    
    query := fmt.Sprintf(`
        INSERT INTO transactions (
            account_id, external_id, booking_date_time, value_date_time,
            amount, currency, description, credit_debit_indicator,
            counterparty_name, counterparty_account, category, is_salary
        ) VALUES %s
        ON CONFLICT (account_id, external_id) DO NOTHING`,
        strings.Join(valueStrings, ","),
    )
    
    _, err := r.db.ExecContext(ctx, query, valueArgs...)
    if err != nil {
        return fmt.Errorf("failed to create batch transactions: %w", err)
    }
    
    return nil
}

func (r *transactionRepository) GetByID(ctx context.Context, id int) (*models.Transaction, error) {
    var tx models.Transaction
    query := `SELECT * FROM transactions WHERE id = $1`
    
    err := r.db.GetContext(ctx, &tx, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("transaction not found")
        }
        return nil, fmt.Errorf("failed to get transaction: %w", err)
    }
    
    return &tx, nil
}

func (r *transactionRepository) GetByExternalID(ctx context.Context, accountID int, externalID string) (*models.Transaction, error) {
    var tx models.Transaction
    query := `SELECT * FROM transactions WHERE account_id = $1 AND external_id = $2`
    
    err := r.db.GetContext(ctx, &tx, query, accountID, externalID)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get transaction: %w", err)
    }
    
    return &tx, nil
}

func (r *transactionRepository) GetAccountTransactions(ctx context.Context, filter models.TransactionFilter) ([]models.Transaction, error) {
    var transactions []models.Transaction
    var args []interface{}
    
    query := `SELECT * FROM transactions WHERE account_id = $1`
    args = append(args, filter.AccountID)
    
    argCount := 1
    
    if filter.FromDate != nil {
        argCount++
        query += fmt.Sprintf(" AND booking_date_time >= $%d", argCount)
        args = append(args, *filter.FromDate)
    }
    
    if filter.ToDate != nil {
        argCount++
        query += fmt.Sprintf(" AND booking_date_time <= $%d", argCount)
        args = append(args, *filter.ToDate)
    }
    
    if filter.IsSalary != nil {
        argCount++
        query += fmt.Sprintf(" AND is_salary = $%d", argCount)
        args = append(args, *filter.IsSalary)
    }
    
    query += " ORDER BY booking_date_time DESC"
    
    if filter.Limit > 0 {
        argCount++
        query += fmt.Sprintf(" LIMIT $%d", argCount)
        args = append(args, filter.Limit)
    }
    
    if filter.Offset > 0 {
        argCount++
        query += fmt.Sprintf(" OFFSET $%d", argCount)
        args = append(args, filter.Offset)
    }
    
    err := r.db.SelectContext(ctx, &transactions, query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to get account transactions: %w", err)
    }
    
    return transactions, nil
}

func (r *transactionRepository) GetUserTransactions(ctx context.Context, userID int, fromDate, toDate time.Time) ([]models.Transaction, error) {
    var transactions []models.Transaction
    query := `
        SELECT t.* 
        FROM transactions t
        JOIN accounts a ON t.account_id = a.id
        WHERE a.user_id = $1 
          AND t.booking_date_time >= $2 
          AND t.booking_date_time <= $3
        ORDER BY t.booking_date_time DESC`
    
    err := r.db.SelectContext(ctx, &transactions, query, userID, fromDate, toDate)
    if err != nil {
        return nil, fmt.Errorf("failed to get user transactions: %w", err)
    }
    
    return transactions, nil
}

func (r *transactionRepository) GetSalaryTransactions(ctx context.Context, userID int) ([]models.Transaction, error) {
    var transactions []models.Transaction
    query := `
        SELECT t.* 
        FROM transactions t
        JOIN accounts a ON t.account_id = a.id
        WHERE a.user_id = $1 
          AND t.amount > 0
          AND t.booking_date_time >= NOW() - INTERVAL '3 months'
        ORDER BY t.amount DESC, t.booking_date_time DESC`
    
    err := r.db.SelectContext(ctx, &transactions, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get salary transactions: %w", err)
    }
    
    return transactions, nil
}

func (r *transactionRepository) MarkAsSalary(ctx context.Context, transactionIDs []int) error {
    query := `UPDATE transactions SET is_salary = true WHERE id = ANY($1)`
    
    _, err := r.db.ExecContext(ctx, query, pq.Array(transactionIDs))
    if err != nil {
        return fmt.Errorf("failed to mark as salary: %w", err)
    }
    
    return nil
}

func (r *transactionRepository) CountAccountTransactions(ctx context.Context, accountID int) (int, error) {
    var count int
    query := `SELECT COUNT(*) FROM transactions WHERE account_id = $1`
    
    err := r.db.GetContext(ctx, &count, query, accountID)
    if err != nil {
        return 0, fmt.Errorf("failed to count transactions: %w", err)
    }
    
    return count, nil
}
