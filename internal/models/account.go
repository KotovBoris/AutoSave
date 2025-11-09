package models

import (
	"time"
)

type Account struct {
	ID             int       `db:"id" json:"id"`
	UserID         int       `db:"user_id" json:"userId"`
	UserBankID     int       `db:"user_bank_id" json:"userBankId"`
	BankID         string    `db:"bank_id" json:"bankId"`
	ExternalID     string    `db:"external_id" json:"externalId"`
	Identification string    `db:"identification" json:"identification"`
	SchemeName     *string   `db:"scheme_name" json:"schemeName,omitempty"`
	AccountType    *string   `db:"account_type" json:"accountType,omitempty"`
	Nickname       *string   `db:"nickname" json:"nickname,omitempty"`
	Balance        float64   `db:"balance" json:"balance"`
	Currency       string    `db:"currency" json:"currency"`
	ServicerName   *string   `db:"servicer_name" json:"servicerName,omitempty"`
	IsActive       bool      `db:"is_active" json:"isActive"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`

	// Joined fields
	BankName string `db:"bank_name" json:"bankName,omitempty"`
}

type AccountResponse struct {
	ID            int       `json:"id"`
	UserBankID    int       `json:"userBankId"`
	BankID        string    `json:"bankId"`
	BankName      string    `json:"bankName"`
	AccountNumber string    `json:"accountNumber"`
	AccountName   *string   `json:"accountName,omitempty"`
	AccountType   *string   `json:"accountType,omitempty"`
	Balance       float64   `json:"balance"`
	Currency      string    `json:"currency"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (a *Account) ToResponse() AccountResponse {
	// Show last 4 digits of account number
	accountNumber := a.Identification
	if len(accountNumber) > 4 {
		accountNumber = "****" + accountNumber[len(accountNumber)-4:]
	}

	return AccountResponse{
		ID:            a.ID,
		UserBankID:    a.UserBankID,
		BankID:        a.BankID,
		BankName:      a.BankName,
		AccountNumber: accountNumber,
		AccountName:   a.Nickname,
		AccountType:   a.AccountType,
		Balance:       a.Balance,
		Currency:      a.Currency,
		UpdatedAt:     a.UpdatedAt,
	}
}

type Balance struct {
	AccountID int       `json:"accountId"`
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"`
	UpdatedAt time.Time `json:"updatedAt"`
}
