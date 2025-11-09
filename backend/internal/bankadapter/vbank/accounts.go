package vbank

import (
    "encoding/json"
    "fmt"
    "net/url"
    "time"
    
    "github.com/KotovBoris/AutoSave/backend/internal/bankadapter"
)

// GetAccounts retrieves all accounts for a client
func (a *Adapter) GetAccounts(token, clientID, consentID, requestingBank string) ([]bankadapter.Account, error) {
    params := url.Values{}
    if clientID != "" {
        params.Set("client_id", clientID)
    }
    
    headers := a.GetConsentHeaders(token, consentID, requestingBank)
    
    fullURL := a.BuildURL(endpointAccounts, map[string]string{"client_id": clientID})
    
    resp, err := a.DoRequest("GET", fullURL, headers, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get accounts: %w", err)
    }
    defer resp.Body.Close()
    
    var response struct {
        Data struct {
            Account []struct {
                AccountID      string   `json:"accountId"`
                Currency       string   `json:"currency"`
                AccountType    string   `json:"accountType"`
                AccountSubType string   `json:"accountSubType"`
                Nickname       string   `json:"nickname"`
                Account        struct {
                    SchemeName     string `json:"schemeName"`
                    Identification string `json:"identification"`
                    Name           string `json:"name"`
                } `json:"account"`
                Servicer struct {
                    SchemeName     string `json:"schemeName"`
                    Identification string `json:"identification"`
                    Name           string `json:"name"`
                } `json:"servicer"`
            } `json:"account"`
        } `json:"data"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode accounts: %w", err)
    }
    
    accounts := make([]bankadapter.Account, 0, len(response.Data.Account))
    for _, acc := range response.Data.Account {
        // Get balance for each account
        balance, _ := a.GetAccountBalance(token, acc.AccountID, consentID, requestingBank)
        if balance == nil {
            balance = &bankadapter.Balance{
                Amount:   0,
                Currency: acc.Currency,
                Type:     "InterimAvailable",
                DateTime: time.Now(),
            }
        }
        
        accounts = append(accounts, bankadapter.Account{
            ID:             acc.AccountID,
            Identification: acc.Account.Identification,
            Currency:       acc.Currency,
            AccountType:    acc.AccountType,
            AccountSubType: acc.AccountSubType,
            Nickname:       acc.Nickname,
            SchemeName:     acc.Account.SchemeName,
            Servicer: bankadapter.Servicer{
                SchemeName:     acc.Servicer.SchemeName,
                Identification: acc.Servicer.Identification,
                Name:           acc.Servicer.Name,
            },
            Balance: *balance,
        })
    }
    
    return accounts, nil
}

// GetAccountDetails retrieves details of a specific account
func (a *Adapter) GetAccountDetails(token, accountID, consentID, requestingBank string) (*bankadapter.Account, error) {
    path := fmt.Sprintf(endpointAccountDetails, accountID)
    headers := a.GetConsentHeaders(token, consentID, requestingBank)
    
    resp, err := a.DoRequest("GET", path, headers, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get account details: %w", err)
    }
    defer resp.Body.Close()
    
    var response struct {
        Data struct {
            Account struct {
                AccountID      string `json:"accountId"`
                Currency       string `json:"currency"`
                AccountType    string `json:"accountType"`
                AccountSubType string `json:"accountSubType"`
                Nickname       string `json:"nickname"`
                Account        struct {
                    SchemeName     string `json:"schemeName"`
                    Identification string `json:"identification"`
                } `json:"account"`
                Servicer struct {
                    SchemeName     string `json:"schemeName"`
                    Identification string `json:"identification"`
                    Name           string `json:"name"`
                } `json:"servicer"`
            } `json:"account"`
        } `json:"data"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode account details: %w", err)
    }
    
    balance, _ := a.GetAccountBalance(token, accountID, consentID, requestingBank)
    if balance == nil {
        balance = &bankadapter.Balance{
            Amount:   0,
            Currency: response.Data.Account.Currency,
            Type:     "InterimAvailable",
            DateTime: time.Now(),
        }
    }
    
    return &bankadapter.Account{
        ID:             response.Data.Account.AccountID,
        Identification: response.Data.Account.Account.Identification,
        Currency:       response.Data.Account.Currency,
        AccountType:    response.Data.Account.AccountType,
        AccountSubType: response.Data.Account.AccountSubType,
        Nickname:       response.Data.Account.Nickname,
        SchemeName:     response.Data.Account.Account.SchemeName,
        Servicer: bankadapter.Servicer{
            SchemeName:     response.Data.Account.Servicer.SchemeName,
            Identification: response.Data.Account.Servicer.Identification,
            Name:           response.Data.Account.Servicer.Name,
        },
        Balance: *balance,
    }, nil
}

// GetAccountBalance retrieves account balance
func (a *Adapter) GetAccountBalance(token, accountID, consentID, requestingBank string) (*bankadapter.Balance, error) {
    path := fmt.Sprintf(endpointAccountBalances, accountID)
    headers := a.GetConsentHeaders(token, consentID, requestingBank)
    
    resp, err := a.DoRequest("GET", path, headers, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get account balance: %w", err)
    }
    defer resp.Body.Close()
    
    var response struct {
        Data struct {
            Balance []struct {
                Amount struct {
                    Amount   string `json:"amount"`
                    Currency string `json:"currency"`
                } `json:"amount"`
                CreditDebitIndicator string `json:"creditDebitIndicator"`
                Type                 string `json:"type"`
                DateTime             string `json:"dateTime"`
            } `json:"balance"`
        } `json:"data"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode balance: %w", err)
    }
    
    if len(response.Data.Balance) == 0 {
        return nil, fmt.Errorf("no balance found")
    }
    
    bal := response.Data.Balance[0]
    amount := 0.0
    fmt.Sscanf(bal.Amount.Amount, "%f", &amount)
    
    dateTime, _ := time.Parse(time.RFC3339, bal.DateTime)
    
    return &bankadapter.Balance{
        Amount:   amount,
        Currency: bal.Amount.Currency,
        Type:     bal.Type,
        DateTime: dateTime,
    }, nil
}

// CreateAccount creates a new account
func (a *Adapter) CreateAccount(token, clientID string, accountType string, initialBalance float64) (*bankadapter.Account, error) {
    params := url.Values{}
    if clientID != "" {
        params.Set("client_id", clientID)
    }
    
    request := map[string]interface{}{
        "account_type":    accountType,
        "initial_balance": initialBalance,
    }
    
    headers := a.GetAuthHeaders(token)
    fullURL := a.BuildURL(endpointAccounts, map[string]string{"client_id": clientID})
    
    resp, err := a.DoRequest("POST", fullURL, headers, request)
    if err != nil {
        return nil, fmt.Errorf("failed to create account: %w", err)
    }
    defer resp.Body.Close()
    
    var response struct {
        AccountID      string `json:"account_id"`
        Identification string `json:"identification"`
        AccountType    string `json:"account_type"`
        Currency       string `json:"currency"`
        Balance        float64 `json:"balance"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode account creation response: %w", err)
    }
    
    return &bankadapter.Account{
        ID:             response.AccountID,
        Identification: response.Identification,
        AccountType:    response.AccountType,
        Currency:       response.Currency,
        Balance: bankadapter.Balance{
            Amount:   response.Balance,
            Currency: response.Currency,
            Type:     "InterimAvailable",
            DateTime: time.Now(),
        },
    }, nil
}

// CloseAccount closes an account
func (a *Adapter) CloseAccount(token, clientID, accountID string, closeRequest bankadapter.AccountCloseRequest) error {
    path := fmt.Sprintf(endpointAccountClose, accountID)
    params := url.Values{}
    if clientID != "" {
        params.Set("client_id", clientID)
    }
    
    headers := a.GetAuthHeaders(token)
    fullURL := a.BuildURL(path, map[string]string{"client_id": clientID})
    
    resp, err := a.DoRequest("PUT", fullURL, headers, closeRequest)
    if err != nil {
        return fmt.Errorf("failed to close account: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 200 && resp.StatusCode != 204 {
        return fmt.Errorf("failed to close account: status %d", resp.StatusCode)
    }
    
    return nil
}

