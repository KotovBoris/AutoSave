package bankadapter

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "time"
    
    "github.com/rs/zerolog"
)

// BaseAdapter provides common HTTP functionality
type BaseAdapter struct {
    ClientID     string
    ClientSecret string
    BaseURL      string
    TeamID       string
    HTTPClient   *http.Client
    Logger       *zerolog.Logger
}

// NewBaseAdapter creates a new base adapter
func NewBaseAdapter(clientID, clientSecret, baseURL, teamID string, logger *zerolog.Logger) *BaseAdapter {
    return &BaseAdapter{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        BaseURL:      baseURL,
        TeamID:       teamID,
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        Logger: logger,
    }
}

// DoRequest performs an HTTP request with common error handling
func (b *BaseAdapter) DoRequest(method, path string, headers map[string]string, body interface{}) (*http.Response, error) {
    fullURL := b.BaseURL + path
    
    var bodyReader io.Reader
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal request body: %w", err)
        }
        bodyReader = bytes.NewReader(jsonBody)
    }
    
    req, err := http.NewRequest(method, fullURL, bodyReader)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    // Set default headers
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Accept", "application/json")
    
    // Set custom headers
    for key, value := range headers {
        req.Header.Set(key, value)
    }
    
    // Log request
    b.Logger.Debug().
        Str("method", method).
        Str("url", fullURL).
        Interface("headers", headers).
        Msg("Making bank API request")
    
    resp, err := b.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }
    
    // Log response
    b.Logger.Debug().
        Int("status", resp.StatusCode).
        Str("url", fullURL).
        Msg("Bank API response")
    
    return resp, nil
}

// ParseResponse reads and unmarshals response body
func (b *BaseAdapter) ParseResponse(resp *http.Response, target interface{}) error {
    defer resp.Body.Close()
    
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read response body: %w", err)
    }
    
    // Check for error status codes
    if resp.StatusCode >= 400 {
        var bankErr BankError
        if err := json.Unmarshal(body, &bankErr); err != nil {
            // If can't parse as BankError, return generic error
            return &BankError{
                Code:    "API_ERROR",
                Message: fmt.Sprintf("API returned status %d: %s", resp.StatusCode, string(body)),
                Status:  resp.StatusCode,
            }
        }
        bankErr.Status = resp.StatusCode
        return &bankErr
    }
    
    if target != nil && len(body) > 0 {
        if err := json.Unmarshal(body, target); err != nil {
            return fmt.Errorf("failed to unmarshal response: %w", err)
        }
    }
    
    return nil
}

// BuildURL builds URL with query parameters
func (b *BaseAdapter) BuildURL(path string, params map[string]string) string {
    baseURL := b.BaseURL + path
    
    if len(params) == 0 {
        return baseURL
    }
    
    u, err := url.Parse(baseURL)
    if err != nil {
        return baseURL
    }
    
    q := u.Query()
    for key, value := range params {
        q.Set(key, value)
    }
    u.RawQuery = q.Encode()
    
    return u.String()
}

// GetAuthHeaders returns common auth headers
func (b *BaseAdapter) GetAuthHeaders(token string) map[string]string {
    return map[string]string{
        "Authorization": "Bearer " + token,
    }
}

// GetConsentHeaders returns headers for consent-based requests
func (b *BaseAdapter) GetConsentHeaders(token, consentID, requestingBank string) map[string]string {
    headers := b.GetAuthHeaders(token)
    if consentID != "" {
        headers["X-Consent-Id"] = consentID
    }
    if requestingBank != "" {
        headers["X-Requesting-Bank"] = requestingBank
    }
    return headers
}

// FormatClientID formats client ID for the bank
func (b *BaseAdapter) FormatClientID(userID int) string {
    return fmt.Sprintf("%s-%d", b.TeamID, userID)
}

// ParseClientID extracts user ID from formatted client ID
func (b *BaseAdapter) ParseClientID(clientID string) (int, error) {
    var userID int
    _, err := fmt.Sscanf(clientID, b.TeamID+"-%d", &userID)
    if err != nil {
        return 0, fmt.Errorf("invalid client ID format: %s", clientID)
    }
    return userID, nil
}

// CalculateInterest calculates interest for deposit
func CalculateInterest(principal float64, rate float64, days int) float64 {
    // Simple interest calculation: P * R * T / 365
    return principal * (rate / 100) * float64(days) / 365
}

// CalculateMaturityDate calculates maturity date for deposit
func CalculateMaturityDate(openDate time.Time, termMonths int) time.Time {
    return openDate.AddDate(0, termMonths, 0)
}

// CalculatePenalty calculates early withdrawal penalty
func CalculatePenalty(accruedInterest float64, daysHeld, totalDays int) float64 {
    if daysHeld < 30 {
        return accruedInterest // Lose all interest if withdrawn within 30 days
    }
    
    // Proportional penalty based on time held
    penaltyRate := 1.0 - (float64(daysHeld) / float64(totalDays))
    return accruedInterest * penaltyRate * 0.5 // 50% of proportional interest
}

