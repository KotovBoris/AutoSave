package bankadapter

import (
	"fmt"
	"math/rand"
	"time"
)

// MockAdapter implements BankAdapter for testing
type MockAdapter struct {
	BankID      string
	BankName    string
	DepositRate float64
	Healthy     bool
}

// NewMockAdapter creates a new mock adapter
func NewMockAdapter(bankID string) *MockAdapter {
	bankNames := map[string]string{
		"vbank": "Virtual Bank",
		"abank": "Awesome Bank",
		"sbank": "Smart Bank",
	}

	rates := map[string]float64{
		"vbank": 8.0,
		"abank": 7.5,
		"sbank": 9.0,
	}

	return &MockAdapter{
		BankID:      bankID,
		BankName:    bankNames[bankID],
		DepositRate: rates[bankID],
		Healthy:     true,
	}
}

func (m *MockAdapter) GetBankToken(clientID, clientSecret string) (*TokenResponse, error) {
	if !m.Healthy {
		return nil, &BankError{Code: ErrCodeBankUnavailable, Message: "Bank is unavailable"}
	}

	return &TokenResponse{
		AccessToken:  fmt.Sprintf("mock_token_%s_%d", m.BankID, time.Now().Unix()),
		TokenType:    "Bearer",
		ExpiresIn:    86400,
		RefreshToken: fmt.Sprintf("mock_refresh_%s_%d", m.BankID, time.Now().Unix()),
		ClientID:     clientID,
		IssuedAt:     time.Now(),
	}, nil
}

func (m *MockAdapter) RefreshToken(refreshToken string) (*TokenResponse, error) {
	return m.GetBankToken("mock_client", "mock_secret")
}

func (m *MockAdapter) CreateAccountConsent(token, clientID, requestingBank string, permissions []string) (*ConsentResponse, error) {
	if !m.Healthy {
		return nil, &BankError{Code: ErrCodeBankUnavailable, Message: "Bank is unavailable"}
	}

	return &ConsentResponse{
		ConsentID:    fmt.Sprintf("consent_%s_%d", m.BankID, time.Now().UnixNano()),
		Status:       "approved",
		Permissions:  permissions,
		ExpiresAt:    time.Now().Add(90 * 24 * time.Hour),
		CreatedAt:    time.Now(),
		AutoApproved: true,
	}, nil
}

func (m *MockAdapter) CreateProductConsent(token, clientID, requestingBank string, permissions []string) (*ConsentResponse, error) {
	return m.CreateAccountConsent(token, clientID, requestingBank, permissions)
}

func (m *MockAdapter) CreatePaymentConsent(token, clientID, requestingBank string, consent PaymentConsentRequest) (*ConsentResponse, error) {
	return m.CreateAccountConsent(token, clientID, requestingBank, []string{"payments"})
}

func (m *MockAdapter) GetConsent(token, consentID string) (*ConsentResponse, error) {
	return &ConsentResponse{
		ConsentID:   consentID,
		Status:      "approved",
		Permissions: []string{"ReadAccountsDetail", "ReadBalances", "ReadTransactionsDetail"},
		ExpiresAt:   time.Now().Add(90 * 24 * time.Hour),
		CreatedAt:   time.Now().Add(-24 * time.Hour),
	}, nil
}

func (m *MockAdapter) DeleteConsent(token, consentID string) error {
	return nil
}

func (m *MockAdapter) GetAccounts(token, clientID, consentID, requestingBank string) ([]Account, error) {
	if !m.Healthy {
		return nil, &BankError{Code: ErrCodeBankUnavailable, Message: "Bank is unavailable"}
	}

	accounts := []Account{
		{
			ID:             fmt.Sprintf("acc_%s_1", m.BankID),
			Identification: "40817810099910001234",
			Currency:       "RUB",
			AccountType:    "Personal",
			Nickname:       "�������� ����",
			SchemeName:     "RU.CBR.PAN",
			Servicer: Servicer{
				SchemeName:     "RU.CBR.BIK",
				Identification: "044525225",
				Name:           m.BankName,
			},
			Balance: Balance{
				Amount:   150000.00,
				Currency: "RUB",
				Type:     "InterimAvailable",
				DateTime: time.Now(),
			},
		},
		{
			ID:             fmt.Sprintf("acc_%s_2", m.BankID),
			Identification: "40817810099910005678",
			Currency:       "RUB",
			AccountType:    "Personal",
			Nickname:       "�������������",
			SchemeName:     "RU.CBR.PAN",
			Servicer: Servicer{
				SchemeName:     "RU.CBR.BIK",
				Identification: "044525225",
				Name:           m.BankName,
			},
			Balance: Balance{
				Amount:   50000.00,
				Currency: "RUB",
				Type:     "InterimAvailable",
				DateTime: time.Now(),
			},
		},
	}

	return accounts, nil
}

func (m *MockAdapter) GetAccountDetails(token, accountID, consentID, requestingBank string) (*Account, error) {
	accounts, err := m.GetAccounts(token, "", consentID, requestingBank)
	if err != nil {
		return nil, err
	}

	for _, acc := range accounts {
		if acc.ID == accountID {
			return &acc, nil
		}
	}

	return nil, &BankError{Code: ErrCodeNotFound, Message: "Account not found"}
}

func (m *MockAdapter) GetAccountBalance(token, accountID, consentID, requestingBank string) (*Balance, error) {
	account, err := m.GetAccountDetails(token, accountID, consentID, requestingBank)
	if err != nil {
		return nil, err
	}
	return &account.Balance, nil
}

func (m *MockAdapter) CreateAccount(token, clientID string, accountType string, initialBalance float64) (*Account, error) {
	return &Account{
		ID:             fmt.Sprintf("acc_%s_%d", m.BankID, time.Now().UnixNano()),
		Identification: fmt.Sprintf("408178100999100%05d", rand.Intn(99999)),
		Currency:       "RUB",
		AccountType:    accountType,
		SchemeName:     "RU.CBR.PAN",
		Balance: Balance{
			Amount:   initialBalance,
			Currency: "RUB",
			Type:     "InterimAvailable",
			DateTime: time.Now(),
		},
	}, nil
}

func (m *MockAdapter) CloseAccount(token, clientID, accountID string, closeRequest AccountCloseRequest) error {
	return nil
}

func (m *MockAdapter) GetTransactions(token, accountID, consentID, requestingBank string, from, to time.Time, limit int) ([]Transaction, error) {
	if !m.Healthy {
		return nil, &BankError{Code: ErrCodeBankUnavailable, Message: "Bank is unavailable"}
	}

	transactions := []Transaction{}

	// Generate salary transactions (monthly on 15th)
	for d := from; d.Before(to); d = d.AddDate(0, 1, 0) {
		salaryDate := time.Date(d.Year(), d.Month(), 15, 10, 0, 0, 0, time.UTC)
		if salaryDate.After(from) && salaryDate.Before(to) {
			transactions = append(transactions, Transaction{
				TransactionID:        fmt.Sprintf("tx_%d", salaryDate.Unix()),
				AccountID:            accountID,
				Amount:               85000.00,
				Currency:             "RUB",
				CreditDebitIndicator: "Credit",
				Status:               "Booked",
				BookingDateTime:      salaryDate,
				ValueDateTime:        salaryDate,
				TransactionInfo: TransInfo{
					Description: "��������",
				},
				CounterpartyName: "��� ������������",
				Category:         "salary",
			})
		}
	}

	// Generate random expense transactions
	categories := []struct {
		name string
		desc string
		min  float64
		max  float64
	}{
		{"groceries", "�����������", 1000, 5000},
		{"transport", "���������", 200, 1000},
		{"restaurants", "��������", 500, 3000},
		{"utilities", "���", 3000, 8000},
	}

	current := from
	for i := 0; i < 30 && current.Before(to); i++ {
		cat := categories[rand.Intn(len(categories))]
		amount := cat.min + rand.Float64()*(cat.max-cat.min)

		transactions = append(transactions, Transaction{
			TransactionID:        fmt.Sprintf("tx_%d_%d", current.Unix(), i),
			AccountID:            accountID,
			Amount:               -amount,
			Currency:             "RUB",
			CreditDebitIndicator: "Debit",
			Status:               "Booked",
			BookingDateTime:      current,
			ValueDateTime:        current,
			TransactionInfo: TransInfo{
				Description: cat.desc,
			},
			Category: cat.name,
		})

		current = current.AddDate(0, 0, rand.Intn(3)+1)
	}

	if limit > 0 && len(transactions) > limit {
		transactions = transactions[:limit]
	}

	return transactions, nil
}

func (m *MockAdapter) GetProducts(token string, productType string) ([]Product, error) {
	products := []Product{
		{
			ProductID:    fmt.Sprintf("prod-%s-deposit-001", m.BankID),
			ProductType:  "deposit",
			ProductName:  "����� ��������",
			Description:  "������������ ������� �����",
			InterestRate: m.DepositRate,
			MinAmount:    10000,
			MaxAmount:    10000000,
			TermMonths:   []int{3, 6, 12, 24},
			Currency:     "RUB",
		},
		{
			ProductID:    fmt.Sprintf("prod-%s-deposit-002", m.BankID),
			ProductType:  "deposit",
			ProductName:  "����� �������������",
			Description:  "����� � ������������ ����������",
			InterestRate: m.DepositRate - 0.5,
			MinAmount:    1000,
			MaxAmount:    5000000,
			TermMonths:   []int{6, 12},
			Currency:     "RUB",
		},
	}

	if productType != "" {
		filtered := []Product{}
		for _, p := range products {
			if p.ProductType == productType {
				filtered = append(filtered, p)
			}
		}
		return filtered, nil
	}

	return products, nil
}

func (m *MockAdapter) GetProductDetails(token, productID string) (*Product, error) {
	products, err := m.GetProducts(token, "")
	if err != nil {
		return nil, err
	}

	for _, p := range products {
		if p.ProductID == productID {
			return &p, nil
		}
	}

	return nil, &BankError{Code: ErrCodeNotFound, Message: "Product not found"}
}

func (m *MockAdapter) GetAgreements(token, clientID, consentID, requestingBank string) ([]Agreement, error) {
	return []Agreement{}, nil
}

func (m *MockAdapter) OpenDeposit(token, clientID, consentID, requestingBank string, request DepositRequest) (*Agreement, error) {
	if !m.Healthy {
		return nil, &BankError{Code: ErrCodeBankUnavailable, Message: "Bank is unavailable"}
	}

	return &Agreement{
		AgreementID:  fmt.Sprintf("agr_%s_%d", m.BankID, time.Now().UnixNano()),
		ProductID:    request.ProductID,
		ProductType:  "deposit",
		ProductName:  "����� ��������",
		Amount:       request.Amount,
		Currency:     "RUB",
		InterestRate: m.DepositRate,
		TermMonths:   request.TermMonths,
		Status:       "active",
		OpenedDate:   time.Now(),
		MaturityDate: time.Now().AddDate(0, request.TermMonths, 0),
	}, nil
}

func (m *MockAdapter) CloseDeposit(token, clientID, consentID, requestingBank, agreementID string) (*CloseDepositResponse, error) {
	accruedInterest := 1500.0 // Mock interest

	return &CloseDepositResponse{
		AgreementID:     agreementID,
		ClosedAt:        time.Now(),
		ReturnedAmount:  50000 + accruedInterest,
		AccruedInterest: accruedInterest,
		PenaltyAmount:   100,
	}, nil
}

func (m *MockAdapter) GetAgreementDetails(token, clientID, consentID, requestingBank, agreementID string) (*Agreement, error) {
	return &Agreement{
		AgreementID:     agreementID,
		ProductID:       fmt.Sprintf("prod-%s-deposit-001", m.BankID),
		ProductType:     "deposit",
		Amount:          50000,
		Currency:        "RUB",
		InterestRate:    m.DepositRate,
		Status:          "active",
		OpenedDate:      time.Now().AddDate(0, -3, 0),
		MaturityDate:    time.Now().AddDate(0, 9, 0),
		AccruedInterest: 1000,
	}, nil
}

func (m *MockAdapter) CreatePayment(token, clientID, requestingBank string, payment PaymentRequest) (*PaymentResponse, error) {
	return &PaymentResponse{
		PaymentID:   fmt.Sprintf("pay_%s_%d", m.BankID, time.Now().UnixNano()),
		Status:      "completed",
		Amount:      payment.Amount,
		Currency:    payment.Currency,
		CreatedAt:   time.Now(),
		CompletedAt: time.Now(),
	}, nil
}

func (m *MockAdapter) GetPaymentStatus(token, clientID, paymentID string) (*PaymentResponse, error) {
	return &PaymentResponse{
		PaymentID:   paymentID,
		Status:      "completed",
		Amount:      10000,
		Currency:    "RUB",
		CreatedAt:   time.Now().Add(-1 * time.Hour),
		CompletedAt: time.Now().Add(-30 * time.Minute),
	}, nil
}

func (m *MockAdapter) GetCards(token, clientID, consentID, requestingBank string) ([]Card, error) {
	return []Card{}, nil
}

func (m *MockAdapter) CreateCard(token, clientID, consentID, requestingBank string, request CreateCardRequest) (*Card, error) {
	return &Card{
		CardID:       fmt.Sprintf("card_%s_%d", m.BankID, time.Now().UnixNano()),
		CardNumber:   "****1234",
		CardType:     request.CardType,
		CardBrand:    "Visa",
		CardStatus:   "active",
		AccountID:    request.AccountNumber,
		ExpiryDate:   time.Now().AddDate(3, 0, 0).Format("01/06"),
		DailyLimit:   100000,
		MonthlyLimit: 1000000,
		IssuedDate:   time.Now(),
	}, nil
}

func (m *MockAdapter) GetBankInfo() BankInfo {
	return BankInfo{
		ID:          m.BankID,
		Name:        m.BankName,
		BaseURL:     fmt.Sprintf("https://%s.mock.api", m.BankID),
		DepositRate: m.DepositRate,
	}
}

func (m *MockAdapter) IsHealthy() bool {
	return m.Healthy
}
