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

// Adapter for SBank - using mock for MVP
type Adapter struct {
    config Config
    logger *zerolog.Logger
    mock   *bankadapter.MockAdapter
}

// NewAdapter creates new SBank adapter
func NewAdapter(cfg Config) bankadapter.BankAdapter {
    return &Adapter{
        config: cfg,
        logger: cfg.Logger.With().Str("bank", "sbank").Logger(),
        mock:   bankadapter.NewMockAdapter("sbank"),
    }
}

func (a *Adapter) GetBankToken(clientID, clientSecret string) (*bankadapter.TokenResponse, error) {
    return a.mock.GetBankToken(clientID, clientSecret)
}

func (a *Adapter) RefreshToken(refreshToken string) (*bankadapter.TokenResponse, error) {
    return a.mock.RefreshToken(refreshToken)
}

func (a *Adapter) CreateAccountConsent(token, clientID, requestingBank string, permissions []string) (*bankadapter.ConsentResponse, error) {
    return a.mock.CreateAccountConsent(token, clientID, requestingBank, permissions)
}

func (a *Adapter) CreateProductConsent(token, clientID, requestingBank string, permissions []string) (*bankadapter.ConsentResponse, error) {
    return a.mock.CreateProductConsent(token, clientID, requestingBank, permissions)
}

func (a *Adapter) CreatePaymentConsent(token, clientID, requestingBank string, consent bankadapter.PaymentConsentRequest) (*bankadapter.ConsentResponse, error) {
    return a.mock.CreatePaymentConsent(token, clientID, requestingBank, consent)
}

func (a *Adapter) GetConsent(token, consentID string) (*bankadapter.ConsentResponse, error) {
    return a.mock.GetConsent(token, consentID)
}

func (a *Adapter) DeleteConsent(token, consentID string) error {
    return a.mock.DeleteConsent(token, consentID)
}

func (a *Adapter) GetAccounts(token, clientID, consentID, requestingBank string) ([]bankadapter.Account, error) {
    return a.mock.GetAccounts(token, clientID, consentID, requestingBank)
}

func (a *Adapter) GetAccountDetails(token, accountID, consentID, requestingBank string) (*bankadapter.Account, error) {
    return a.mock.GetAccountDetails(token, accountID, consentID, requestingBank)
}

func (a *Adapter) GetAccountBalance(token, accountID, consentID, requestingBank string) (*bankadapter.Balance, error) {
    return a.mock.GetAccountBalance(token, accountID, consentID, requestingBank)
}

func (a *Adapter) CreateAccount(token, clientID string, accountType string, initialBalance float64) (*bankadapter.Account, error) {
    return a.mock.CreateAccount(token, clientID, accountType, initialBalance)
}

func (a *Adapter) CloseAccount(token, clientID, accountID string, closeRequest bankadapter.AccountCloseRequest) error {
    return a.mock.CloseAccount(token, clientID, accountID, closeRequest)
}

func (a *Adapter) GetTransactions(token, accountID, consentID, requestingBank string, from, to time.Time, limit int) ([]bankadapter.Transaction, error) {
    return a.mock.GetTransactions(token, accountID, consentID, requestingBank, from, to, limit)
}

func (a *Adapter) GetProducts(token string, productType string) ([]bankadapter.Product, error) {
    return a.mock.GetProducts(token, productType)
}

func (a *Adapter) GetProductDetails(token, productID string) (*bankadapter.Product, error) {
    return a.mock.GetProductDetails(token, productID)
}

func (a *Adapter) GetAgreements(token, clientID, consentID, requestingBank string) ([]bankadapter.Agreement, error) {
    return a.mock.GetAgreements(token, clientID, consentID, requestingBank)
}

func (a *Adapter) OpenDeposit(token, clientID, consentID, requestingBank string, request bankadapter.DepositRequest) (*bankadapter.Agreement, error) {
    return a.mock.OpenDeposit(token, clientID, consentID, requestingBank, request)
}

func (a *Adapter) CloseDeposit(token, clientID, consentID, requestingBank, agreementID string) (*bankadapter.CloseDepositResponse, error) {
    return a.mock.CloseDeposit(token, clientID, consentID, requestingBank, agreementID)
}

func (a *Adapter) GetAgreementDetails(token, clientID, consentID, requestingBank, agreementID string) (*bankadapter.Agreement, error) {
    return a.mock.GetAgreementDetails(token, clientID, consentID, requestingBank, agreementID)
}

func (a *Adapter) CreatePayment(token, clientID, requestingBank string, payment bankadapter.PaymentRequest) (*bankadapter.PaymentResponse, error) {
    return a.mock.CreatePayment(token, clientID, requestingBank, payment)
}

func (a *Adapter) GetPaymentStatus(token, clientID, paymentID string) (*bankadapter.PaymentResponse, error) {
    return a.mock.GetPaymentStatus(token, clientID, paymentID)
}

func (a *Adapter) GetCards(token, clientID, consentID, requestingBank string) ([]bankadapter.Card, error) {
    return a.mock.GetCards(token, clientID, consentID, requestingBank)
}

func (a *Adapter) CreateCard(token, clientID, consentID, requestingBank string, request bankadapter.CreateCardRequest) (*bankadapter.Card, error) {
    return a.mock.CreateCard(token, clientID, consentID, requestingBank, request)
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
    return a.mock.IsHealthy()
}
