package repository

import (
    "context"
    "database/sql"
    "fmt"
    
    "github.com/KotovBoris/AutoSave/backend/internal/models"
    "github.com/jmoiron/sqlx"
)

type bankRepository struct {
    db *sqlx.DB
}

func NewBankRepository(db *sqlx.DB) BankRepository {
    return &bankRepository{db: db}
}

func (r *bankRepository) GetAll(ctx context.Context) ([]models.Bank, error) {
    var banks []models.Bank
    query := `SELECT id, name, api_base_url, deposit_rate, is_active, created_at FROM banks WHERE is_active = true ORDER BY name`
    
    err := r.db.SelectContext(ctx, &banks, query)
    if err != nil {
        return nil, fmt.Errorf("failed to get banks: %w", err)
    }
    
    return banks, nil
}

func (r *bankRepository) GetByID(ctx context.Context, id string) (*models.Bank, error) {
    var bank models.Bank
    query := `SELECT id, name, api_base_url, deposit_rate, is_active, created_at FROM banks WHERE id = $1`
    
    err := r.db.GetContext(ctx, &bank, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("bank not found")
        }
        return nil, fmt.Errorf("failed to get bank: %w", err)
    }
    
    return &bank, nil
}

func (r *bankRepository) CreateConnection(ctx context.Context, conn *models.BankConnection) error {
    query := `
        INSERT INTO user_banks (
            user_id, bank_id, external_client_id, bank_token,
            token_expires_at, account_consent_id, product_consent_id,
            payment_consent_id, connected, connected_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
        RETURNING id, connected_at`
    
    err := r.db.QueryRowxContext(ctx, query,
        conn.UserID, conn.BankID, conn.ExternalClientID, conn.BankToken,
        conn.TokenExpiresAt, conn.AccountConsentID, conn.ProductConsentID,
        conn.PaymentConsentID, conn.Connected,
    ).Scan(&conn.ID, &conn.ConnectedAt)
    
    if err != nil {
        return fmt.Errorf("failed to create bank connection: %w", err)
    }
    
    return nil
}

func (r *bankRepository) GetUserConnections(ctx context.Context, userID int) ([]models.BankConnection, error) {
    var connections []models.BankConnection
    query := `
        SELECT 
            ub.id, ub.user_id, ub.bank_id, b.name as bank_name,
            ub.external_client_id, ub.bank_token, ub.token_expires_at,
            ub.account_consent_id, ub.product_consent_id, ub.payment_consent_id,
            ub.connected, ub.connected_at, ub.last_sync_at, ub.error
        FROM user_banks ub
        JOIN banks b ON ub.bank_id = b.id
        WHERE ub.user_id = $1 AND ub.connected = true
        ORDER BY ub.connected_at DESC`
    
    err := r.db.SelectContext(ctx, &connections, query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to get user connections: %w", err)
    }
    
    return connections, nil
}

func (r *bankRepository) GetConnection(ctx context.Context, userID int, bankID string) (*models.BankConnection, error) {
    var conn models.BankConnection
    query := `
        SELECT 
            ub.id, ub.user_id, ub.bank_id, b.name as bank_name,
            ub.external_client_id, ub.bank_token, ub.token_expires_at,
            ub.account_consent_id, ub.product_consent_id, ub.payment_consent_id,
            ub.connected, ub.connected_at, ub.last_sync_at, ub.error
        FROM user_banks ub
        JOIN banks b ON ub.bank_id = b.id
        WHERE ub.user_id = $1 AND ub.bank_id = $2`
    
    err := r.db.GetContext(ctx, &conn, query, userID, bankID)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get connection: %w", err)
    }
    
    return &conn, nil
}

func (r *bankRepository) GetConnectionByID(ctx context.Context, id int) (*models.BankConnection, error) {
    var conn models.BankConnection
    query := `
        SELECT 
            id, user_id, bank_id, external_client_id, bank_token,
            token_expires_at, account_consent_id, product_consent_id,
            payment_consent_id, connected, connected_at, last_sync_at, error
        FROM user_banks
        WHERE id = $1`
    
    err := r.db.GetContext(ctx, &conn, query, id)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("connection not found")
        }
        return nil, fmt.Errorf("failed to get connection: %w", err)
    }
    
    return &conn, nil
}

func (r *bankRepository) UpdateConnection(ctx context.Context, conn *models.BankConnection) error {
    query := `
        UPDATE user_banks 
        SET bank_token = $2, token_expires_at = $3, 
            account_consent_id = $4, product_consent_id = $5,
            payment_consent_id = $6, last_sync_at = $7, error = $8
        WHERE id = $1`
    
    _, err := r.db.ExecContext(ctx, query,
        conn.ID, conn.BankToken, conn.TokenExpiresAt,
        conn.AccountConsentID, conn.ProductConsentID,
        conn.PaymentConsentID, conn.LastSyncAt, conn.Error,
    )
    
    if err != nil {
        return fmt.Errorf("failed to update connection: %w", err)
    }
    
    return nil
}

func (r *bankRepository) DeleteConnection(ctx context.Context, userID int, bankID string) error {
    query := `UPDATE user_banks SET connected = false WHERE user_id = $1 AND bank_id = $2`
    
    _, err := r.db.ExecContext(ctx, query, userID, bankID)
    if err != nil {
        return fmt.Errorf("failed to delete connection: %w", err)
    }
    
    return nil
}

