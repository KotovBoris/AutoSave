package models

import "time"

type Loan struct {
    ID              int        `db:"id" json:"id"`
    UserID          int        `db:"user_id" json:"userId"`
    Name            string     `db:"name" json:"name"`
    OriginalDebt    float64    `db:"original_debt" json:"originalDebt"`
    CurrentDebt     float64    `db:"current_debt" json:"currentDebt"`
    Rate            float64    `db:"rate" json:"rate"`
    MonthlyPayment  float64    `db:"monthly_payment" json:"monthlyPayment"`
    AutopayEnabled  bool       `db:"autopay_enabled" json:"autopayEnabled"`
    AutopayBankID   *string    `db:"autopay_bank_id" json:"autopayBankId,omitempty"`
    AutopayDay      *int       `db:"autopay_day" json:"autopayDay,omitempty"`
    Status          string     `db:"status" json:"status"`
    NextPaymentDate *time.Time `db:"next_payment_date" json:"nextPaymentDate,omitempty"`
    CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
    PaidOffAt       *time.Time `db:"paid_off_at" json:"paidOffAt,omitempty"`
    UpdatedAt       time.Time  `db:"updated_at" json:"updatedAt"`
    
    // Joined fields
    AutopayBankName *string        `db:"autopay_bank_name" json:"autopayBankName,omitempty"`
    Payments        []LoanPayment  `json:"payments,omitempty"`
}

type CreateLoanRequest struct {
    Name           string   `json:"name" validate:"required,min=1,max=100"`
    CurrentDebt    float64  `json:"currentDebt" validate:"required,min=1000"`
    Rate           float64  `json:"rate" validate:"required,min=0.1,max=100"`
    MonthlyPayment float64  `json:"monthlyPayment" validate:"required,min=100"`
    AutopayEnabled bool     `json:"autopayEnabled"`
    AutopayBankID  *string  `json:"autopayBankId,omitempty" validate:"omitempty,oneof=vbank abank sbank"`
    AutopayDay     *int     `json:"autopayDay,omitempty" validate:"omitempty,min=1,max=31"`
}

type UpdateLoanRequest struct {
    Name           *string  `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
    MonthlyPayment *float64 `json:"monthlyPayment,omitempty" validate:"omitempty,min=100"`
    AutopayEnabled *bool    `json:"autopayEnabled,omitempty"`
    AutopayBankID  *string  `json:"autopayBankId,omitempty" validate:"omitempty,oneof=vbank abank sbank"`
    AutopayDay     *int     `json:"autopayDay,omitempty" validate:"omitempty,min=1,max=31"`
}

type LoanSchedule struct {
    LoanID   int              `json:"loanId"`
    Schedule []ScheduleEntry  `json:"schedule"`
    Summary  ScheduleSummary  `json:"summary"`
}

type ScheduleEntry struct {
    Month         int       `json:"month"`
    Date          time.Time `json:"date"`
    Payment       float64   `json:"payment"`
    Principal     float64   `json:"principal"`
    Interest      float64   `json:"interest"`
    RemainingDebt float64   `json:"remainingDebt"`
}

type ScheduleSummary struct {
    TotalPayments    float64 `json:"totalPayments"`
    TotalInterest    float64 `json:"totalInterest"`
    MonthsRemaining  int     `json:"monthsRemaining"`
}

type LoanPayment struct {
    ID             int        `db:"id" json:"id"`
    LoanID         int        `db:"loan_id" json:"loanId"`
    UserID         int        `db:"user_id" json:"userId"`
    Amount         float64    `db:"amount" json:"amount"`
    IsAutopay      bool       `db:"is_autopay" json:"isAutopay"`
    BankPaymentID  *string    `db:"bank_payment_id" json:"bankPaymentId,omitempty"`
    Status         string     `db:"status" json:"status"`
    ScheduledDate  time.Time  `db:"scheduled_date" json:"scheduledDate"`
    CompletedAt    *time.Time `db:"completed_at" json:"completedAt,omitempty"`
    Error          *string    `db:"error" json:"error,omitempty"`
    CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
}

