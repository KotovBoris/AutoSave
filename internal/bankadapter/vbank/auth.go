package vbank

import (
    "encoding/json"
    "fmt"
    "net/url"
    "time"
    
    "github.com/autosave/backend/internal/bankadapter"
)

// GetBankToken obtains bank access token
func (a *Adapter) GetBankToken() (*bankadapter.TokenResponse, error) {
    params := url.Values{}
    params.Set("client_id", a.config.ClientID)
    params.Set("client_secret", a.config.ClientSecret)
    
    fullURL := a.BaseURL + endpointAuthToken + "?" + params.Encode()
    
    resp, err := a.HTTPClient.Post(fullURL, "application/x-www-form-urlencoded", nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get bank token: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 200 {
        return nil, &bankadapter.BankError{
            Code:    bankadapter.ErrCodeUnauthorized,
            Message: "Failed to authenticate with bank",
            Status:  resp.StatusCode,
        }
    }
    
    var tokenResp struct {
        AccessToken  string `json:"access_token"`
        TokenType    string `json:"token_type"`
        ExpiresIn    int    `json:"expires_in"`
        RefreshToken string `json:"refresh_token,omitempty"`
        ClientID     string `json:"client_id,omitempty"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
        return nil, fmt.Errorf("failed to decode token response: %w", err)
    }
    
    return &bankadapter.TokenResponse{
        AccessToken:  tokenResp.AccessToken,
        TokenType:    tokenResp.TokenType,
        ExpiresIn:    tokenResp.ExpiresIn,
        RefreshToken: tokenResp.RefreshToken,
        ClientID:     tokenResp.ClientID,
        IssuedAt:     time.Now(),
    }, nil
}

// RefreshToken refreshes access token
func (a *Adapter) RefreshToken(refreshToken string) (*bankadapter.TokenResponse, error) {
    // VBank uses same endpoint for refresh
    // In production, this would use refresh_token grant type
    return a.GetBankToken(a.ClientID, a.ClientSecret)
}

