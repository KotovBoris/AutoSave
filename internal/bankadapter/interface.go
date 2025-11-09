package bankadapter

import (
    "time"
)

// BankAdapter is the common interface for all bank integrations
type BankAdapter interface {
    // Authentication
    GetBankToken() (*TokenResponse, error)
    RefreshToken(refreshToken string) (*TokenResponse, error)
    
    // Consents
    CreateAccountConsent(token, clientID, requestingBank string, permissions []string) (*ConsentResponse, error)
    CreateProductConsent(token, clientID, requestingBank string, permissions []string) (*ConsentResponse, error)
    CreatePaymentConsent(token, clientID, requestingBank string, consent PaymentConsentRequest) (*ConsentResponse, error)
    GetConsent(token, consentID string) (*ConsentResponse, error)
    DeleteConsent(token, consentID string) error
    
    // Accounts
    GetAccounts(token, clientID, consentID, requestingBank string) ([]Account, error)
    GetAccountDetails(token, accountID, consentID, requestingBank string) (*Account, error)
    GetAccountBalance(token, accountID, consentID, requestingBank string) (*Balance, error)
    CreateAccount(token, clientID string, accountType string, initialBalance float64) (*Account, error)
    CloseAccount(token, clientID, accountID string, closeRequest AccountCloseRequest) error
    
    // Transactions
    GetTransactions(token, accountID, consentID, requestingBank string, from, to time.Time, limit int) ([]Transaction, error)
    
    // Products
    GetProducts(token string, productType string) ([]Product, error)
    GetProductDetails(token, productID string) (*Product, error)
    
    // Product Agreements (Deposits, Loans, Cards)
    GetAgreements(token, clientID, consentID, requestingBank string) ([]Agreement, error)
    OpenDeposit(token, clientID, consentID, requestingBank string, request DepositRequest) (*Agreement, error)
    CloseDeposit(token, clientID, consentID, requestingBank, agreementID string) (*CloseDepositResponse, error)
    GetAgreementDetails(token, clientID, consentID, requestingBank, agreementID string) (*Agreement, error)
    
    // Payments
    CreatePayment(token, clientID, requestingBank string, payment PaymentRequest) (*PaymentResponse, error)
    GetPaymentStatus(token, clientID, paymentID string) (*PaymentResponse, error)
    
    // Cards
    GetCards(token, clientID, consentID, requestingBank string) ([]Card, error)
    CreateCard(token, clientID, consentID, requestingBank string, request CreateCardRequest) (*Card, error)
    
    // Utility
    GetBankInfo() BankInfo
    IsHealthy() bool
}

// BankInfo contains static information about the bank
type BankInfo struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    BaseURL     string  `json:"baseUrl"`
    DepositRate float64 `json:"depositRate"`
}

// TokenResponse from bank OAuth
type TokenResponse struct {
    AccessToken  string    `json:"access_token"`
    TokenType    string    `json:"token_type"`
    ExpiresIn    int       `json:"expires_in"`
    RefreshToken string    `json:"refresh_token,omitempty"`
    ClientID     string    `json:"client_id,omitempty"`
    IssuedAt     time.Time `json:"issued_at"`
}

// ConsentResponse from consent APIs
type ConsentResponse struct {
    ConsentID     string    `json:"consent_id"`
    Status        string    `json:"status"`
    Permissions   []string  `json:"permissions,omitempty"`
    ExpiresAt     time.Time `json:"expires_at,omitempty"`
    CreatedAt     time.Time `json:"created_at"`
    AutoApproved  bool      `json:"auto_approved,omitempty"`
}

// Account from bank API
type Account struct {
    ID             string    `json:"account_id"`
    Identification string    `json:"identification"`
    Currency       string    `json:"currency"`
    AccountType    string    `json:"account_type"`
    AccountSubType string    `json:"account_sub_type,omitempty"`
    Nickname       string    `json:"nickname,omitempty"`
    SchemeName     string    `json:"scheme_name"`
    Servicer       Servicer  `json:"servicer"`
    Balance        Balance   `json:"balance"`
}

// Servicer information
type Servicer struct {
    SchemeName     string `json:"scheme_name"`
    Identification string `json:"identification"`
    Name           string `json:"name"`
}

// Balance of an account
type Balance struct {
    Amount   float64   `json:"amount"`
    Currency string    `json:"currency"`
    Type     string    `json:"type"`
    DateTime time.Time `json:"date_time"`
}

// Transaction from bank API
type Transaction struct {
    TransactionID        string    `json:"transaction_id"`
    AccountID            string    `json:"account_id"`
    Amount               float64   `json:"amount"`
    Currency             string    `json:"currency"`
    CreditDebitIndicator string    `json:"credit_debit_indicator"`
    Status               string    `json:"status"`
    BookingDateTime      time.Time `json:"booking_date_time"`
    ValueDateTime        time.Time `json:"value_date_time"`
    TransactionInfo      TransInfo `json:"transaction_information"`
    CounterpartyName     string    `json:"counterparty_name,omitempty"`
    CounterpartyAccount  string    `json:"counterparty_account,omitempty"`
    Category             string    `json:"category,omitempty"`
}

// TransInfo contains transaction details
type TransInfo struct {
    Description          string `json:"description"`
    TransactionReference string `json:"transaction_reference,omitempty"`
}

// Product from catalog
type Product struct {
    ProductID          string            `json:"product_id"`
    ProductType        string            `json:"product_type"`
    ProductName        string            `json:"product_name"`
    Description        string            `json:"description"`
    InterestRate       float64           `json:"interest_rate,omitempty"`
    MinAmount          float64           `json:"min_amount,omitempty"`
    MaxAmount          float64           `json:"max_amount,omitempty"`
    TermMonths         []int             `json:"term_months,omitempty"`
    Currency           string            `json:"currency"`
    Features           []string          `json:"features,omitempty"`
    Requirements       []string          `json:"requirements,omitempty"`
    AdditionalInfo     map[string]string `json:"additional_info,omitempty"`
}

// Agreement (deposit, loan, card contract)
type Agreement struct {
    AgreementID   string    `json:"agreement_id"`
    ProductID     string    `json:"product_id"`
    ProductType   string    `json:"product_type"`
    ProductName   string    `json:"product_name"`
    Amount        float64   `json:"amount"`
    Currency      string    `json:"currency"`
    InterestRate  float64   `json:"interest_rate,omitempty"`
    TermMonths    int       `json:"term_months,omitempty"`
    Status        string    `json:"status"`
    OpenedDate    time.Time `json:"opened_date"`
    MaturityDate  time.Time `json:"maturity_date,omitempty"`
    ClosedDate    time.Time `json:"closed_date,omitempty"`
    NextPayment   time.Time `json:"next_payment,omitempty"`
    CurrentDebt   float64   `json:"current_debt,omitempty"`
    AccruedInterest float64 `json:"accrued_interest,omitempty"`
}

// Card information
type Card struct {
    CardID         string    `json:"card_id"`
    CardNumber     string    `json:"card_number"`
    CardType       string    `json:"card_type"`
    CardBrand      string    `json:"card_brand"`
    CardStatus     string    `json:"card_status"`
    AccountID      string    `json:"account_id"`
    ExpiryDate     string    `json:"expiry_date"`
    DailyLimit     float64   `json:"daily_limit"`
    MonthlyLimit   float64   `json:"monthly_limit"`
    IssuedDate     time.Time `json:"issued_date"`
}

// Request types

// PaymentConsentRequest for creating payment consent
type PaymentConsentRequest struct {
    ConsentType        string   `json:"consent_type"`
    Amount             float64  `json:"amount,omitempty"`
    Currency           string   `json:"currency"`
    DebtorAccount      string   `json:"debtor_account"`
    CreditorAccount    string   `json:"creditor_account,omitempty"`
    CreditorName       string   `json:"creditor_name,omitempty"`
    Reference          string   `json:"reference,omitempty"`
    MaxUses            int      `json:"max_uses,omitempty"`
    MaxAmountPerPayment float64 `json:"max_amount_per_payment,omitempty"`
    ValidUntil         time.Time `json:"valid_until,omitempty"`
}

// AccountCloseRequest for closing account
type AccountCloseRequest struct {
    Action              string `json:"action"`
    DestinationAccountID string `json:"destination_account_id,omitempty"`
}

// DepositRequest for opening deposit
type DepositRequest struct {
    ProductID        string  `json:"product_id"`
    Amount           float64 `json:"amount"`
    TermMonths       int     `json:"term_months"`
    SourceAccountID  string  `json:"source_account_id"`
    AutoRenewal      bool    `json:"auto_renewal"`
}

// CloseDepositResponse when closing deposit
type CloseDepositResponse struct {
    AgreementID      string  `json:"agreement_id"`
    ClosedAt         time.Time `json:"closed_at"`
    ReturnedAmount   float64 `json:"returned_amount"`
    AccruedInterest  float64 `json:"accrued_interest"`
    PenaltyAmount    float64 `json:"penalty_amount,omitempty"`
}

// PaymentRequest for making payment
type PaymentRequest struct {
    DebtorAccountID   string  `json:"debtor_account_id"`
    CreditorAccountID string  `json:"creditor_account_id"`
    CreditorBankCode  string  `json:"creditor_bank_code,omitempty"`
    Amount            float64 `json:"amount"`
    Currency          string  `json:"currency"`
    Reference         string  `json:"reference"`
    Description       string  `json:"description,omitempty"`
}

// PaymentResponse for payment status
type PaymentResponse struct {
    PaymentID   string    `json:"payment_id"`
    Status      string    `json:"status"`
    Amount      float64   `json:"amount"`
    Currency    string    `json:"currency"`
    CreatedAt   time.Time `json:"created_at"`
    CompletedAt time.Time `json:"completed_at,omitempty"`
    Error       string    `json:"error,omitempty"`
}

// CreateCardRequest for creating new card
type CreateCardRequest struct {
    AccountNumber string `json:"account_number"`
    CardName      string `json:"card_name,omitempty"`
    CardType      string `json:"card_type"`
    DailyLimit    float64 `json:"daily_limit,omitempty"`
    MonthlyLimit  float64 `json:"monthly_limit,omitempty"`
}

// Error types

// BankError represents an error from bank API
type BankError struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
    Status  int    `json:"status"`
}

func (e *BankError) Error() string {
    return e.Message
}

// Common error codes
const (
    ErrCodeUnauthorized      = "UNAUTHORIZED"
    ErrCodeForbidden         = "FORBIDDEN"
    ErrCodeNotFound          = "NOT_FOUND"
    ErrCodeInvalidRequest    = "INVALID_REQUEST"
    ErrCodeInsufficientFunds = "INSUFFICIENT_FUNDS"
    ErrCodeConsentRequired   = "CONSENT_REQUIRED"
    ErrCodeConsentExpired    = "CONSENT_EXPIRED"
    ErrCodeBankUnavailable   = "BANK_UNAVAILABLE"
    ErrCodeRateLimited       = "RATE_LIMITED"
)

