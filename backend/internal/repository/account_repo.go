package repository

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    
    "github.com/KotovBoris/AutoSave/backend/internal/models"
    "github.com/jmoiron/sqlx"
)

type accountRepository struct {
    db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) AccountRepository {
    return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *models.Account) error {
    query := `
        INSERT INTO accounts (
            user_id, user_bank_id, bank_id, external_id, identification,
            scheme_name, account_type, nickname, balance, currency,
            servicer_name, is_active
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id, created_at, updated_at`
    
    err := r.db.QueryRowxContext(ctx, query,
        account.UserID, account.UserBankID, account.BankID, account.ExternalID,
        account.Identification, account.SchemeName, account.AccountType,
        account.Nickname, account.Balance, account.Currency,
        account.ServicerName, account.IsActive,
    ).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create account: %w", err)
    }
    
    return nil
}

func (r *accountRepository) CreateBatch(ctx context.Context, accounts []models.Account) error {
    if len(accounts) == 0 {
        return nil
    }
    
    valueStrings := make([]string, 0, len(accounts))
    valueArgs := make([]interface{}, 0, len(accounts)*12)
    
    for i, acc := range accounts {
        valueStrings = append(valueStrings, fmt.Sprintf(
            "($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
            i*12+1, i*12+2, i*12+3, i*12+4, i*12+5, i*12+6,
            i*12+7, i*12+8, i*12+9, i*12+10, i*12+11, i*12+12,
        ))
        
        valueArgs = append(valueArgs,
            acc.UserID, acc.UserBankID, acc.BankID, acc.ExternalID,
            acc.Identification, acc.SchemeName, acc.AccountType,
            acc.Nickname, acc.Balance, acc.Currency,
            acc.ServicerName, acc.IsActive,
        )
    }
    
    query := fmt.Sprintf(`
        INSERT INTO accounts (
            user_id, user_bank_id, bank_id, external_id, identification,
            scheme_name, account_type, nickname, balance, currency,
            servicer_name, is_active
        ) VALUES %s
        ON CONFLICT (user_bank_id, external_id) 
        DO UPDATE SET 
            balance = EXCLUDED.balance,
            updated_at = NOW()`,
        strings.Join(valueStrings, ","),
    )
    
    _, err := r.db.ExecContext(ctx, query, valueArgs...)
    if err != nil {
        return fmt.Errorf("failed to create batch accounts: %w", err)
    }
    
    return nil
}

func (r *accountRepository) GetByID(ctx context.Context, id int) (*models.Account, error) {
    var account models.Account
    query := `
        SELECT 
            a.*, b.name as bank_name
        FROM accounts a
        JOIN banks b ON a.bank_id = b.id
        WHERE a.id = $1`
    
    err := r.db.GetContext(ctx, &account, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("account not found")
        }
        return nil, fmt.Errorf("failed to get account: %w", err)
    }
    
    return &account, nil
}

func (r *accountRepository) GetByExternalID(ctx context.Context, userBankID int, externalID string) (*models.Account, error) {
    var account models.Account
    query := `
        SELECT 
            a.*, b.name as bank_name
        FROM accounts a
        JOIN banks b ON a.bank_id = b.id
        WHERE a.user_bank_id = $1 AND a.external_id = $2`
    
    err := r.db.GetContext(ctx, &account, query, userBankID, externalID)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get account: %w", err)
    }
    
    return &account, nil
}

func (r *accountRepository) GetUserAccounts(ctx context.Context, userID int) ([]models.Account, error) {
    var accounts []models.Account
    query := `
        SELECT 
            a.*, b.name as bank_name
        FROM accounts a
        JOIN banks b ON a.bank_id = b.id
        WHERE a.user_id = $1 AND a.is_active = true
        ORDER BY a.bank_id, a.balance DESC`
    
    err := r.db.SelectContext(ctx, &accounts, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user accounts: %w", err)
    }
    
    return accounts, nil
}

func (r *accountRepository) GetBankAccounts(ctx context.Context, userID int, bankID string) ([]models.Account, error) {
    var accounts []models.Account
    query := `
        SELECT 
            a.*, b.name as bank_name
        FROM accounts a
        JOIN banks b ON a.bank_id = b.id
        WHERE a.user_id = $1 AND a.bank_id = $2 AND a.is_active = true
        ORDER BY a.balance DESC`
    
    err := r.db.SelectContext(ctx, &accounts, query, userID, bankID)
    if err != nil {
        return nil, fmt.Errorf("failed to get bank accounts: %w", err)
    }
    
    return accounts, nil
}

func (r *accountRepository) Update(ctx context.Context, account *models.Account) error {
    query := `
        UPDATE accounts 
        SET nickname = $2, balance = $3, updated_at = NOW()
        WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query, account.ID, account.Nickname, account.Balance)
    if err != nil {
        return fmt.Errorf("failed to update account: %w", err)
    }
    
    return nil
}

func (r *accountRepository) UpdateBalance(ctx context.Context, id int, balance float64) error {
    query := `UPDATE accounts SET balance = $2, updated_at = NOW() WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query, id, balance)
    if err != nil {
        return fmt.Errorf("failed to update balance: %w", err)
    }
    
    return nil
}

func (r *accountRepository) Delete(ctx context.Context, id int) error {
    query := `UPDATE accounts SET is_active = false, updated_at = NOW() WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete account: %w", err)
    }
    
    return nil
}

