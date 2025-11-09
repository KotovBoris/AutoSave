package models

import (
    "time"
    "encoding/json"
    "database/sql/driver"
)

type Operation struct {
    ID               int        `db:"id" json:"id"`
    UserID           int        `db:"user_id" json:"userId"`
    Type             string     `db:"type" json:"type"`
    Amount           *float64   `db:"amount" json:"amount,omitempty"`
    RelatedGoalID    *int       `db:"related_goal_id" json:"relatedGoalId,omitempty"`
    RelatedLoanID    *int       `db:"related_loan_id" json:"relatedLoanId,omitempty"`
    RelatedDepositID *int       `db:"related_deposit_id" json:"relatedDepositId,omitempty"`
    Status           string     `db:"status" json:"status"`
    Error            *string    `db:"error" json:"error,omitempty"`
    Metadata         JSONB      `db:"metadata" json:"metadata,omitempty"`
    CreatedAt        time.Time  `db:"created_at" json:"createdAt"`
    
    // Joined fields
    Goal    *Goal    `json:"goal,omitempty"`
    Loan    *Loan    `json:"loan,omitempty"`
    Deposit *Deposit `json:"deposit,omitempty"`
}

type OperationType string

const (
    OperationDepositOpened   OperationType = "deposit_opened"
    OperationDepositClosed   OperationType = "deposit_closed"
    OperationLoanPayment     OperationType = "loan_payment"
    OperationEmergencyWithdraw OperationType = "emergency_withdraw"
    OperationGoalCreated     OperationType = "goal_created"
    OperationGoalCompleted   OperationType = "goal_completed"
)

type OperationStatus string

const (
    OperationStatusSuccess OperationStatus = "success"
    OperationStatusFailed  OperationStatus = "failed"
    OperationStatusPending OperationStatus = "pending"
)

type EmergencyWithdrawRequest struct {
    Amount float64 `json:"amount" validate:"required,min=1"`
}

type EmergencyWithdrawPlan struct {
    RequestedAmount   float64                `json:"requestedAmount"`
    DepositsToClose   []DepositToClose      `json:"depositsToClose"`
    TotalAmount       float64                `json:"totalAmount"`
    TotalAccruedInterest float64             `json:"totalAccruedInterest"`
    TotalLostInterest float64                `json:"totalLostInterest"`
    TotalReturned     float64                `json:"totalReturned"`
    AffectedGoals     []AffectedGoal         `json:"affectedGoals"`
}

type DepositToClose struct {
    DepositID       int     `json:"depositId"`
    GoalID          int     `json:"goalId"`
    GoalName        string  `json:"goalName"`
    Amount          float64 `json:"amount"`
    AccruedInterest float64 `json:"accruedInterest"`
    LostInterest    float64 `json:"lostInterest"`
    BankID          string  `json:"bankId"`
}

type AffectedGoal struct {
    GoalID        int     `json:"goalId"`
    GoalName      string  `json:"goalName"`
    CurrentAmount float64 `json:"currentAmount"`
    AfterWithdraw float64 `json:"afterWithdraw"`
    WillBePaused  bool    `json:"willBePaused"`
}

type EmergencyWithdrawConfirm struct {
    DepositIDs []int `json:"depositIds" validate:"required,min=1"`
}

type EmergencyWithdrawResult struct {
    Success         bool                   `json:"success"`
    ClosedDeposits  []int                  `json:"closedDeposits"`
    TotalReturned   float64                `json:"totalReturned"`
    TotalLostInterest float64              `json:"totalLostInterest"`
    Operations      []Operation            `json:"operations"`
}

// JSONB type for PostgreSQL jsonb columns
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
    if j == nil {
        return nil, nil
    }
    return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
    if value == nil {
        *j = make(JSONB)
        return nil
    }
    
    bytes, ok := value.([]byte)
    if !ok {
        return fmt.Errorf("cannot scan %T into JSONB", value)
    }
    
    return json.Unmarshal(bytes, j)
}
