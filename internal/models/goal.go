package models

import (
    "time"
)

type Goal struct {
    ID               int        `db:"id" json:"id"`
    UserID           int        `db:"user_id" json:"userId"`
    Name             string     `db:"name" json:"name"`
    TargetAmount     float64    `db:"target_amount" json:"targetAmount"`
    CurrentAmount    float64    `db:"current_amount" json:"currentAmount"`
    MonthlyAmount    float64    `db:"monthly_amount" json:"monthlyAmount"`
    BankID           string     `db:"bank_id" json:"bankId"`
    DepositRate      float64    `db:"deposit_rate" json:"depositRate"`
    Position         int        `db:"position" json:"position"`
    Status           string     `db:"status" json:"status"`
    NextDepositDate  *time.Time `db:"next_deposit_date" json:"nextDepositDate,omitempty"`
    CreatedAt        time.Time  `db:"created_at" json:"createdAt"`
    CompletedAt      *time.Time `db:"completed_at" json:"completedAt,omitempty"`
    UpdatedAt        time.Time  `db:"updated_at" json:"updatedAt"`
    
    // Joined fields
    BankName         string     `db:"bank_name" json:"bankName,omitempty"`
    Deposits         []Deposit  `json:"deposits,omitempty"`
}

type CreateGoalRequest struct {
    Name          string  `json:"name" validate:"required,min=1,max=100"`
    TargetAmount  float64 `json:"targetAmount" validate:"required,min=1000"`
    MonthlyAmount float64 `json:"monthlyAmount" validate:"required,min=1000"`
    BankID        string  `json:"bankId" validate:"required,oneof=vbank abank sbank"`
}

type UpdateGoalRequest struct {
    Name          *string  `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
    MonthlyAmount *float64 `json:"monthlyAmount,omitempty" validate:"omitempty,min=1000"`
}

type ReorderGoalsRequest struct {
    GoalIDs []int `json:"goalIds" validate:"required,min=1"`
}

type GoalPlan struct {
    MonthsToComplete    int       `json:"monthsToComplete"`
    EstimatedInterest   float64   `json:"estimatedInterest"`
    EstimatedTotal      float64   `json:"estimatedTotal"`
    EstimatedCompletion time.Time `json:"estimatedCompletion"`
}

type GoalResponse struct {
    ID                   int        `json:"id"`
    Name                 string     `json:"name"`
    TargetAmount         float64    `json:"targetAmount"`
    CurrentAmount        float64    `json:"currentAmount"`
    MonthlyAmount        float64    `json:"monthlyAmount"`
    BankID               string     `json:"bankId"`
    BankName             string     `json:"bankName"`
    DepositRate          float64    `json:"depositRate"`
    Position             int        `json:"position"`
    Status               string     `json:"status"`
    NextDepositDate      *time.Time `json:"nextDepositDate,omitempty"`
    CreatedAt            time.Time  `json:"createdAt"`
    CompletedAt          *time.Time `json:"completedAt,omitempty"`
    Deposits             []Deposit  `json:"deposits"`
    EstimatedCompletion  *time.Time `json:"estimatedCompletion,omitempty"`
    EstimatedInterest    float64    `json:"estimatedInterest"`
    ProgressPercentage   float64    `json:"progressPercentage"`
}

type CloseGoalResponse struct {
    Message         string              `json:"message"`
    ClosedDeposits  []ClosedDepositInfo `json:"closedDeposits"`
    TotalReturned   float64             `json:"totalReturned"`
    TotalLostInterest float64           `json:"totalLostInterest"`
}

type ClosedDepositInfo struct {
    DepositID       int     `json:"depositId"`
    Amount          float64 `json:"amount"`
    AccruedInterest float64 `json:"accruedInterest"`
    LostInterest    float64 `json:"lostInterest"`
}
