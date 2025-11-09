package sbank

import (
    "encoding/json"
    "fmt"
    "time"
    
    "github.com/autosave/backend/internal/bankadapter"
)

// CreateAccountConsent creates consent for account access
// NOTE: SBank requires MANUAL consent approval (unlike VBank/ABank which auto-approve)
// This means consent will be created with status "AwaitingAuthorisation"
// and client must manually approve it in SBank UI before it can be used
func (a *Adapter) CreateAccountConsent(token, clientID, requestingBank string, permissions []string) (*bankadapter.ConsentResponse, error) {
    request := map[string]interface{}{
        "client_id":            clientID,
        "permissions":          permissions,
        "reason":               "Account aggregation for AutoSave app",
        "requesting_bank":      requestingBank,
        "requesting_bank_name": "AutoSave Platform",
    }
    
    headers := a.GetAuthHeaders(token)
    headers["X-Requesting-Bank"] = requestingBank
    
    resp, err := a.DoRequest("POST", endpointAccountConsents, headers, request)
    if err != nil {
        return nil, fmt.Errorf("failed to create account consent: %w", err)
    }
    defer resp.Body.Close()
    
    var consentResp struct {
        Status       string   `json:"status"`
        ConsentID    string   `json:"consent_id"`
        Permissions  []string `json:"permissions"`
        ExpiresAt    string   `json:"expires_at,omitempty"`
        AutoApproved bool     `json:"auto_approved"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&consentResp); err != nil {
        return nil, fmt.Errorf("failed to decode consent response: %w", err)
    }
    
    var expiresAt time.Time
    if consentResp.ExpiresAt != "" {
        expiresAt, _ = time.Parse(time.RFC3339, consentResp.ExpiresAt)
    } else {
        expiresAt = time.Now().Add(90 * 24 * time.Hour)
    }
    
    return &bankadapter.ConsentResponse{
        ConsentID:    consentResp.ConsentID,
        Status:       consentResp.Status,
        Permissions:  consentResp.Permissions,
        ExpiresAt:    expiresAt,
        CreatedAt:    time.Now(),
        AutoApproved: consentResp.AutoApproved,
    }, nil
}

// CreateProductConsent creates consent for product management
func (a *Adapter) CreateProductConsent(token, clientID, requestingBank string, permissions []string) (*bankadapter.ConsentResponse, error) {
    request := map[string]interface{}{
        "requesting_bank":        requestingBank,
        "client_id":              clientID,
        "read_product_agreements": contains(permissions, "read_product_agreements"),
        "open_product_agreements": contains(permissions, "open_product_agreements"),
        "close_product_agreements": contains(permissions, "close_product_agreements"),
        "allowed_product_types":   []string{"deposit", "card"},
        "max_amount":              1000000.00,
        "valid_until":             time.Now().Add(365 * 24 * time.Hour).Format(time.RFC3339),
        "reason":                  "Product management for AutoSave app",
    }
    
    headers := a.GetAuthHeaders(token)
    
    resp, err := a.DoRequest("POST", endpointProductAgreementConsents, headers, request)
    if err != nil {
        return nil, fmt.Errorf("failed to create product consent: %w", err)
    }
    defer resp.Body.Close()
    
    var consentResp struct {
        ConsentID    string `json:"consent_id"`
        Status       string `json:"status"`
        AutoApproved bool   `json:"auto_approved"`
        ValidUntil   string `json:"valid_until"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&consentResp); err != nil {
        return nil, fmt.Errorf("failed to decode product consent response: %w", err)
    }
    
    validUntil, _ := time.Parse(time.RFC3339, consentResp.ValidUntil)
    
    return &bankadapter.ConsentResponse{
        ConsentID:    consentResp.ConsentID,
        Status:       consentResp.Status,
        Permissions:  permissions,
        ExpiresAt:    validUntil,
        CreatedAt:    time.Now(),
        AutoApproved: consentResp.AutoApproved,
    }, nil
}

// CreatePaymentConsent creates consent for payments
func (a *Adapter) CreatePaymentConsent(token, clientID, requestingBank string, consent bankadapter.PaymentConsentRequest) (*bankadapter.ConsentResponse, error) {
    request := map[string]interface{}{
        "requesting_bank":   requestingBank,
        "client_id":         clientID,
        "consent_type":      consent.ConsentType,
        "amount":            consent.Amount,
        "currency":          consent.Currency,
        "debtor_account":    consent.DebtorAccount,
        "creditor_account":  consent.CreditorAccount,
        "creditor_name":     consent.CreditorName,
        "reference":         consent.Reference,
        "max_uses":          consent.MaxUses,
        "max_amount_per_payment": consent.MaxAmountPerPayment,
        "valid_until":       consent.ValidUntil.Format(time.RFC3339),
    }
    
    headers := a.GetAuthHeaders(token)
    headers["X-Requesting-Bank"] = requestingBank
    
    resp, err := a.DoRequest("POST", endpointPaymentConsents, headers, request)
    if err != nil {
        return nil, fmt.Errorf("failed to create payment consent: %w", err)
    }
    defer resp.Body.Close()
    
    var consentResp struct {
        ConsentID string `json:"consent_id"`
        Status    string `json:"status"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&consentResp); err != nil {
        return nil, fmt.Errorf("failed to decode payment consent response: %w", err)
    }
    
    return &bankadapter.ConsentResponse{
        ConsentID: consentResp.ConsentID,
        Status:    consentResp.Status,
        ExpiresAt: consent.ValidUntil,
        CreatedAt: time.Now(),
    }, nil
}

// GetConsent retrieves consent details
func (a *Adapter) GetConsent(token, consentID string) (*bankadapter.ConsentResponse, error) {
    path := fmt.Sprintf(endpointAccountConsent, consentID)
    headers := a.GetAuthHeaders(token)
    
    resp, err := a.DoRequest("GET", path, headers, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get consent: %w", err)
    }
    defer resp.Body.Close()
    
    var response struct {
        Data struct {
            ConsentID    string    `json:"consentId"`
            Status       string    `json:"status"`
            Permissions  []string  `json:"permissions"`
            ExpirationDateTime string `json:"expirationDateTime"`
            CreationDateTime string `json:"creationDateTime"`
        } `json:"data"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode consent: %w", err)
    }
    
    expiresAt, _ := time.Parse(time.RFC3339, response.Data.ExpirationDateTime)
    createdAt, _ := time.Parse(time.RFC3339, response.Data.CreationDateTime)
    
    return &bankadapter.ConsentResponse{
        ConsentID:   response.Data.ConsentID,
        Status:      response.Data.Status,
        Permissions: response.Data.Permissions,
        ExpiresAt:   expiresAt,
        CreatedAt:   createdAt,
    }, nil
}

// DeleteConsent revokes consent
func (a *Adapter) DeleteConsent(token, consentID string) error {
    path := fmt.Sprintf(endpointAccountConsent, consentID)
    headers := a.GetAuthHeaders(token)
    
    resp, err := a.DoRequest("DELETE", path, headers, nil)
    if err != nil {
        return fmt.Errorf("failed to delete consent: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 204 && resp.StatusCode != 200 {
        return fmt.Errorf("failed to delete consent: status %d", resp.StatusCode)
    }
    
    return nil
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}


