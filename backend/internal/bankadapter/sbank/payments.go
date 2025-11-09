package sbank

import (
    "encoding/json"
    "fmt"
    "net/url"
    "time"
    
    "github.com/KotovBoris/AutoSave/backend/internal/bankadapter"
)

// CreatePayment creates a new payment
func (a *Adapter) CreatePayment(token, clientID, requestingBank string, payment bankadapter.PaymentRequest) (*bankadapter.PaymentResponse, error) {
    params := url.Values{}
    if clientID != "" {
        params.Set("client_id", clientID)
    }
    
    headers := a.GetAuthHeaders(token)
    if requestingBank != "" {
        headers["X-Requesting-Bank"] = requestingBank
    }
    
    body := map[string]interface{}{
        "data": map[string]interface{}{
            "initiation": map[string]interface{}{
                "instructedAmount": map[string]interface{}{
                    "amount":   fmt.Sprintf("%.2f", payment.Amount),
                    "currency": payment.Currency,
                },
                "debtorAccount": map[string]interface{}{
                    "schemeName":     "RU.CBR.PAN",
                    "identification": payment.DebtorAccountID,
                },
                "creditorAccount": map[string]interface{}{
                    "schemeName":     "RU.CBR.PAN",
                    "identification": payment.CreditorAccountID,
                    "bank_code":      payment.CreditorBankCode,
                },
                "comment": payment.Description,
            },
        },
    }
    
    fullURL := a.BaseURL + endpointPayments
    if len(params) > 0 {
        fullURL += "?" + params.Encode()
    }
    
    resp, err := a.DoRequest("POST", fullURL, headers, body)
    if err != nil {
        return nil, fmt.Errorf("failed to create payment: %w", err)
    }
    defer resp.Body.Close()
    
    var response struct {
        Data struct {
            PaymentID    string `json:"paymentId"`
            Status       string `json:"status"`
            CreationDateTime string `json:"creationDateTime"`
            Amount       string `json:"amount"`
            Currency     string `json:"currency"`
        } `json:"data"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode payment response: %w", err)
    }
    
    amount := 0.0
    fmt.Sscanf(response.Data.Amount, "%f", &amount)
    
    createdAt, _ := time.Parse(time.RFC3339, response.Data.CreationDateTime)
    
    return &bankadapter.PaymentResponse{
        PaymentID: response.Data.PaymentID,
        Status:    response.Data.Status,
        Amount:    amount,
        Currency:  response.Data.Currency,
        CreatedAt: createdAt,
    }, nil
}

// GetPaymentStatus retrieves payment status
func (a *Adapter) GetPaymentStatus(token, clientID, paymentID string) (*bankadapter.PaymentResponse, error) {
    path := fmt.Sprintf(endpointPayment, paymentID)
    params := url.Values{}
    if clientID != "" {
        params.Set("client_id", clientID)
    }
    
    headers := a.GetAuthHeaders(token)
    
    fullURL := a.BaseURL + path
    if len(params) > 0 {
        fullURL += "?" + params.Encode()
    }
    
    resp, err := a.DoRequest("GET", fullURL, headers, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get payment status: %w", err)
    }
    defer resp.Body.Close()
    
    var response struct {
        Data struct {
            PaymentID    string `json:"paymentId"`
            Status       string `json:"status"`
            CreationDateTime string `json:"creationDateTime"`
            StatusUpdateDateTime string `json:"statusUpdateDateTime"`
            Amount       string `json:"amount"`
            Currency     string `json:"currency"`
        } `json:"data"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode payment status: %w", err)
    }
    
    amount := 0.0
    fmt.Sscanf(response.Data.Amount, "%f", &amount)
    
    createdAt, _ := time.Parse(time.RFC3339, response.Data.CreationDateTime)
    completedAt, _ := time.Parse(time.RFC3339, response.Data.StatusUpdateDateTime)
    
    paymentResp := &bankadapter.PaymentResponse{
        PaymentID: response.Data.PaymentID,
        Status:    response.Data.Status,
        Amount:    amount,
        Currency:  response.Data.Currency,
        CreatedAt: createdAt,
    }
    
    if response.Data.Status == "completed" {
        paymentResp.CompletedAt = completedAt
    }
    
    return paymentResp, nil
}

// GetCards retrieves all cards
func (a *Adapter) GetCards(token, clientID, consentID, requestingBank string) ([]bankadapter.Card, error) {
    params := url.Values{}
    if clientID != "" {
        params.Set("client_id", clientID)
    }
    
    headers := a.GetConsentHeaders(token, consentID, requestingBank)
    
    fullURL := a.BaseURL + endpointCards
    if len(params) > 0 {
        fullURL += "?" + params.Encode()
    }
    
    resp, err := a.DoRequest("GET", fullURL, headers, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get cards: %w", err)
    }
    defer resp.Body.Close()
    
    var response struct {
        Cards []struct {
            CardID       string  `json:"cardId"`
            CardNumber   string  `json:"cardNumber"`
            CardType     string  `json:"cardType"`
            CardBrand    string  `json:"cardBrand"`
            CardStatus   string  `json:"cardStatus"`
            AccountID    string  `json:"accountId"`
            ExpiryDate   string  `json:"expiryDate"`
            DailyLimit   float64 `json:"dailyLimit"`
            MonthlyLimit float64 `json:"monthlyLimit"`
            IssuedDate   string  `json:"issuedDate"`
        } `json:"cards"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode cards: %w", err)
    }
    
    cards := make([]bankadapter.Card, 0, len(response.Cards))
    for _, c := range response.Cards {
        issuedDate, _ := time.Parse(time.RFC3339, c.IssuedDate)
        
        cards = append(cards, bankadapter.Card{
            CardID:       c.CardID,
            CardNumber:   c.CardNumber,
            CardType:     c.CardType,
            CardBrand:    c.CardBrand,
            CardStatus:   c.CardStatus,
            AccountID:    c.AccountID,
            ExpiryDate:   c.ExpiryDate,
            DailyLimit:   c.DailyLimit,
            MonthlyLimit: c.MonthlyLimit,
            IssuedDate:   issuedDate,
        })
    }
    
    return cards, nil
}

// CreateCard creates a new card
func (a *Adapter) CreateCard(token, clientID, consentID, requestingBank string, request bankadapter.CreateCardRequest) (*bankadapter.Card, error) {
    params := url.Values{}
    if clientID != "" {
        params.Set("client_id", clientID)
    }
    
    headers := a.GetConsentHeaders(token, consentID, requestingBank)
    
    fullURL := a.BaseURL + endpointCards
    if len(params) > 0 {
        fullURL += "?" + params.Encode()
    }
    
    resp, err := a.DoRequest("POST", fullURL, headers, request)
    if err != nil {
        return nil, fmt.Errorf("failed to create card: %w", err)
    }
    defer resp.Body.Close()
    
    var card struct {
        CardID       string  `json:"cardId"`
        CardNumber   string  `json:"cardNumber"`
        CardType     string  `json:"cardType"`
        CardBrand    string  `json:"cardBrand"`
        CardStatus   string  `json:"cardStatus"`
        AccountID    string  `json:"accountId"`
        ExpiryDate   string  `json:"expiryDate"`
        DailyLimit   float64 `json:"dailyLimit"`
        MonthlyLimit float64 `json:"monthlyLimit"`
        IssuedDate   string  `json:"issuedDate"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&card); err != nil {
        return nil, fmt.Errorf("failed to decode card: %w", err)
    }
    
    issuedDate, _ := time.Parse(time.RFC3339, card.IssuedDate)
    
    return &bankadapter.Card{
        CardID:       card.CardID,
        CardNumber:   card.CardNumber,
        CardType:     card.CardType,
        CardBrand:    card.CardBrand,
        CardStatus:   card.CardStatus,
        AccountID:    card.AccountID,
        ExpiryDate:   card.ExpiryDate,
        DailyLimit:   card.DailyLimit,
        MonthlyLimit: card.MonthlyLimit,
        IssuedDate:   issuedDate,
    }, nil
}


