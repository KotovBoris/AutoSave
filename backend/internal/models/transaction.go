package models

import "time"

type Transaction struct {
    ID                   int        `db:"id" json:"id"`
    AccountID            int        `db:"account_id" json:"accountId"`
    ExternalID           string     `db:"external_id" json:"externalId"`
    BookingDateTime      time.Time  `db:"booking_date_time" json:"bookingDateTime"`
    ValueDateTime        *time.Time `db:"value_date_time" json:"valueDateTime,omitempty"`
    Amount               float64    `db:"amount" json:"amount"`
    Currency             string     `db:"currency" json:"currency"`
    Description          *string    `db:"description" json:"description,omitempty"`
    CreditDebitIndicator *string    `db:"credit_debit_indicator" json:"creditDebitIndicator,omitempty"`
    CounterpartyName     *string    `db:"counterparty_name" json:"counterpartyName,omitempty"`
    CounterpartyAccount  *string    `db:"counterparty_account" json:"counterpartyAccount,omitempty"`
    Category             *string    `db:"category" json:"category,omitempty"`
    IsSalary             bool       `db:"is_salary" json:"isSalary"`
    CreatedAt            time.Time  `db:"created_at" json:"createdAt"`
}

type TransactionFilter struct {
    AccountID int
    FromDate  *time.Time
    ToDate    *time.Time
    IsSalary  *bool
    Limit     int
    Offset    int
}

type TransactionsResponse struct {
    AccountID     int           `json:"accountId"`
    Transactions  []Transaction `json:"transactions"`
    Total         int           `json:"total,omitempty"`
}

type SalaryDetection struct {
    TransactionID    int     `json:"transactionId"`
    Date            string  `json:"date"`
    Amount          float64 `json:"amount"`
    Counterparty    string  `json:"counterparty"`
    AccountID       int     `json:"accountId"`
    Confidence      string  `json:"confidence"`
    AutoSelected    bool    `json:"autoSelected"`
}

type ConfirmSalariesRequest struct {
    SalaryTransactionIDs []int `json:"salaryTransactionIds" validate:"required,min=1"`
}

type SalaryAnalysis struct {
    AvgSalary       float64      `json:"avgSalary"`
    AvgExpenses     float64      `json:"avgExpenses"`
    SavingsCapacity float64      `json:"savingsCapacity"`
    SalaryDates     []int        `json:"salaryDates"`
    Analysis        AnalysisData `json:"analysis"`
}

type AnalysisData struct {
    TotalIncome    float64 `json:"totalIncome"`
    TotalExpenses  float64 `json:"totalExpenses"`
    PeriodMonths   int     `json:"periodMonths"`
}

