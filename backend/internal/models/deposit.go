package models

import "time"

type Deposit struct {
    ID               int        `db:"id" json:"id"`
    GoalID           int        `db:"goal_id" json:"goalId"`
    UserID           int        `db:"user_id" json:"userId"`
    BankID           string     `db:"bank_id" json:"bankId"`
    ProductID        *string    `db:"product_id" json:"productId,omitempty"`
    AgreementID      *string    `db:"agreement_id" json:"agreementId,omitempty"`
    Amount           float64    `db:"amount" json:"amount"`
    Rate             float64    `db:"rate" json:"rate"`
    TermMonths       int        `db:"term_months" json:"termMonths"`
    Status           string     `db:"status" json:"status"`
    OpenedAt         *time.Time `db:"opened_at" json:"openedAt,omitempty"`
    MaturesAt        *time.Time `db:"matures_at" json:"maturesAt,omitempty"`
    ClosedAt         *time.Time `db:"closed_at" json:"closedAt,omitempty"`
    AccruedInterest  float64    `db:"accrued_interest" json:"accruedInterest"`
    Error            *string    `db:"error" json:"error,omitempty"`
    CreatedAt        time.Time  `db:"created_at" json:"createdAt"`
    UpdatedAt        time.Time  `db:"updated_at" json:"updatedAt"`
}

type CreateDepositRequest struct {
    GoalID        int     `json:"goalId" validate:"required"`
    Amount        float64 `json:"amount" validate:"required,min=1000"`
    SourceAccountID string `json:"sourceAccountId" validate:"required"`
}

type DepositStatus string

const (
    DepositStatusPending DepositStatus = "pending"
    DepositStatusActive  DepositStatus = "active"
    DepositStatusMatured DepositStatus = "matured"
    DepositStatusClosed  DepositStatus = "closed"
    DepositStatusFailed  DepositStatus = "failed"
)

