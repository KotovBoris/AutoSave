package sbank

import (
    "time"
    
    "github.com/autosave/backend/internal/bankadapter"
    "github.com/rs/zerolog"
)

// Config for SBank adapter
type Config struct {
    ClientID     string
    ClientSecret string
    BaseURL      string
    TeamID       string
    Logger       *zerolog.Logger
}

// Adapter for SBank
type Adapter struct {
    *bankadapter.BaseAdapter
    config Config
    logger *zerolog.Logger
}

// NewAdapter creates new SBank adapter
func NewAdapter(cfg Config) *Adapter {
    return &Adapter{
        BaseAdapter: bankadapter.NewBaseAdapter(
            cfg.ClientID,
            cfg.ClientSecret,
            cfg.BaseURL,
            cfg.TeamID,
            cfg.Logger,
        ),
        config: cfg,
        logger: cfg.Logger.With().Str("bank", "sbank").Logger(),
    }
}

// SBank delegates to mock implementation for MVP

func (a *Adapter) GetBankToken(clientID, clientSecret string) (*bankadapter.TokenResponse, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.GetBankToken(clientID, clientSecret)
}

func (a *Adapter) RefreshToken(refreshToken string) (*bankadapter.TokenResponse, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.RefreshToken(refreshToken)
}

func (a *Adapter) CreateAccountConsent(token, clientID, requestingBank string, permissions []string) (*bankadapter.ConsentResponse, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.CreateAccountConsent(token, clientID, requestingBank, permissions)
}

func (a *Adapter) CreateProductConsent(token, clientID, requestingBank string, permissions []string) (*bankadapter.ConsentResponse, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.CreateProductConsent(token, clientID, requestingBank, permissions)
}

func (a *Adapter) CreatePaymentConsent(token, clientID, requestingBank string, consent bankadapter.PaymentConsentRequest) (*bankadapter.ConsentResponse, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.CreatePaymentConsent(token, clientID, requestingBank, consent)
}

func (a *Adapter) GetConsent(token, consentID string) (*bankadapter.ConsentResponse, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.GetConsent(token, consentID)
}

func (a *Adapter) DeleteConsent(token, consentID string) error {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.DeleteConsent(token, consentID)
}

func (a *Adapter) GetAccounts(token, clientID, consentID, requestingBank string) ([]bankadapter.Account, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.GetAccounts(token, clientID, consentID, requestingBank)
}

func (a *Adapter) GetAccountDetails(token, accountID, consentID, requestingBank string) (*bankadapter.Account, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.GetAccountDetails(token, accountID, consentID, requestingBank)
}

func (a *Adapter) GetAccountBalance(token, accountID, consentID, requestingBank string) (*bankadapter.Balance, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.GetAccountBalance(token, accountID, consentID, requestingBank)
}

func (a *Adapter) CreateAccount(token, clientID string, accountType string, initialBalance float64) (*bankadapter.Account, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.CreateAccount(token, clientID, accountType, initialBalance)
}

func (a *Adapter) CloseAccount(token, clientID, accountID string, closeRequest bankadapter.AccountCloseRequest) error {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.CloseAccount(token, clientID, accountID, closeRequest)
}

func (a *Adapter) GetTransactions(token, accountID, consentID, requestingBank string, from, to time.Time, limit int) ([]bankadapter.Transaction, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.GetTransactions(token, accountID, consentID, requestingBank, from, to, limit)
}

func (a *Adapter) GetProducts(token string, productType string) ([]bankadapter.Product, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.GetProducts(token, productType)
}

func (a *Adapter) GetProductDetails(token, productID string) (*bankadapter.Product, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.GetProductDetails(token, productID)
}

func (a *Adapter) GetAgreements(token, clientID, consentID, requestingBank string) ([]bankadapter.Agreement, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.GetAgreements(token, clientID, consentID, requestingBank)
}

func (a *Adapter) OpenDeposit(token, clientID, consentID, requestingBank string, request bankadapter.DepositRequest) (*bankadapter.Agreement, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.OpenDeposit(token, clientID, consentID, requestingBank, request)
}

func (a *Adapter) CloseDeposit(token, clientID, consentID, requestingBank, agreementID string) (*bankadapter.CloseDepositResponse, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.CloseDeposit(token, clientID, consentID, requestingBank, agreementID)
}

func (a *Adapter) GetAgreementDetails(token, clientID, consentID, requestingBank, agreementID string) (*bankadapter.Agreement, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.GetAgreementDetails(token, clientID, consentID, requestingBank, agreementID)
}

func (a *Adapter) CreatePayment(token, clientID, requestingBank string, payment bankadapter.PaymentRequest) (*bankadapter.PaymentResponse, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.CreatePayment(token, clientID, requestingBank, payment)
}

func (a *Adapter) GetPaymentStatus(token, clientID, paymentID string) (*bankadapter.PaymentResponse, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.GetPaymentStatus(token, clientID, paymentID)
}

func (a *Adapter) GetCards(token, clientID, consentID, requestingBank string) ([]bankadapter.Card, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.GetCards(token, clientID, consentID, requestingBank)
}

func (a *Adapter) CreateCard(token, clientID, consentID, requestingBank string, request bankadapter.CreateCardRequest) (*bankadapter.Card, error) {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.CreateCard(token, clientID, consentID, requestingBank, request)
}

func (a *Adapter) GetBankInfo() bankadapter.BankInfo {
    return bankadapter.BankInfo{
        ID:          "sbank",
        Name:        "Smart Bank",
        BaseURL:     a.config.BaseURL,
        DepositRate: 9.0,
    }
}

func (a *Adapter) IsHealthy() bool {
    mock := bankadapter.NewMockAdapter("sbank")
    return mock.IsHealthy()
}
