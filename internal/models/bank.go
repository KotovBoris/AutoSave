package models

import "time"

type Bank struct {
    ID          string    `db:"id" json:"id"`
    Name        string    `db:"name" json:"name"`
    APIBaseURL  string    `db:"api_base_url" json:"apiBaseUrl"`
    DepositRate float64   `db:"deposit_rate" json:"depositRate"`
    IsActive    bool      `db:"is_active" json:"isActive"`
    CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

type BankConnection struct {
    ID                 int        `db:"id" json:"id"`
    UserID             int        `db:"user_id" json:"userId"`
    BankID             string     `db:"bank_id" json:"bankId"`
    BankName           string     `db:"bank_name" json:"bankName,omitempty"`
    ExternalClientID   string     `db:"external_client_id" json:"externalClientId"`
    BankToken          string     `db:"bank_token" json:"-"`
    TokenExpiresAt     *time.Time `db:"token_expires_at" json:"tokenExpiresAt,omitempty"`
    AccountConsentID   *string    `db:"account_consent_id" json:"accountConsentId,omitempty"`
    ProductConsentID   *string    `db:"product_consent_id" json:"productConsentId,omitempty"`
    PaymentConsentID   *string    `db:"payment_consent_id" json:"paymentConsentId,omitempty"`
    Connected          bool       `db:"connected" json:"connected"`
    ConnectedAt        time.Time  `db:"connected_at" json:"connectedAt"`
    LastSyncAt         *time.Time `db:"last_sync_at" json:"lastSyncAt,omitempty"`
    Error              *string    `db:"error" json:"error,omitempty"`
}

type ConnectBankRequest struct {
    BankID string `json:"bankId" validate:"required,oneof=vbank abank sbank"`
}

type SyncBankResponse struct {
    SyncedBanks  []string           `json:"syncedBanks"`
    FailedBanks  []BankSyncError    `json:"failedBanks"`
    LastSyncAt   time.Time          `json:"lastSyncAt"`
}

type BankSyncError struct {
    BankID string `json:"bankId"`
    Error  string `json:"error"`
}
