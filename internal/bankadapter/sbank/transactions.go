package sbank

import (
    "encoding/json"
    "fmt"
    "net/url"
    "strconv"
    "time"
    
    "github.com/autosave/backend/internal/bankadapter"
)

// GetTransactions retrieves account transactions
func (a *Adapter) GetTransactions(token, accountID, consentID, requestingBank string, from, to time.Time, limit int) ([]bankadapter.Transaction, error) {
    path := fmt.Sprintf(endpointAccountTransactions, accountID)
    
    params := url.Values{}
    params.Set("from_booking_date_time", from.Format(time.RFC3339))
    params.Set("to_booking_date_time", to.Format(time.RFC3339))
    if limit > 0 {
        params.Set("limit", strconv.Itoa(limit))
    }
    
    headers := a.GetConsentHeaders(token, consentID, requestingBank)
    fullURL := a.BaseURL + path + "?" + params.Encode()
    
    resp, err := a.DoRequest("GET", fullURL, headers, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get transactions: %w", err)
    }
    defer resp.Body.Close()
    
    var response struct {
        Data struct {
            Transaction []struct {
                TransactionID        string `json:"transactionId"`
                AccountID            string `json:"accountId"`
                Amount               struct {
                    Amount   string `json:"amount"`
                    Currency string `json:"currency"`
                } `json:"amount"`
                CreditDebitIndicator string `json:"creditDebitIndicator"`
                Status               string `json:"status"`
                BookingDateTime      string `json:"bookingDateTime"`
                ValueDateTime        string `json:"valueDateTime"`
                TransactionInfo      struct {
                    Description          string `json:"description"`
                    TransactionReference string `json:"transactionReference"`
                } `json:"transactionInformation"`
                CounterpartyName    string `json:"counterpartyName"`
                CounterpartyAccount string `json:"counterpartyAccount"`
                Category            string `json:"category"`
            } `json:"transaction"`
        } `json:"data"`
        Links struct {
            Self  string `json:"self"`
            First string `json:"first"`
            Next  string `json:"next"`
            Last  string `json:"last"`
        } `json:"links"`
        Meta struct {
            TotalPages int `json:"totalPages"`
            FirstAvailableDateTime string `json:"firstAvailableDateTime"`
            LastAvailableDateTime string `json:"lastAvailableDateTime"`
        } `json:"meta"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode transactions: %w", err)
    }
    
    transactions := make([]bankadapter.Transaction, 0, len(response.Data.Transaction))
    for _, tx := range response.Data.Transaction {
        amount := 0.0
        fmt.Sscanf(tx.Amount.Amount, "%f", &amount)
        
        // Make amount negative for debits
        if tx.CreditDebitIndicator == "Debit" {
            amount = -amount
        }
        
        bookingDateTime, _ := time.Parse(time.RFC3339, tx.BookingDateTime)
        valueDateTime, _ := time.Parse(time.RFC3339, tx.ValueDateTime)
        
        transactions = append(transactions, bankadapter.Transaction{
            TransactionID:        tx.TransactionID,
            AccountID:            tx.AccountID,
            Amount:               amount,
            Currency:             tx.Amount.Currency,
            CreditDebitIndicator: tx.CreditDebitIndicator,
            Status:               tx.Status,
            BookingDateTime:      bookingDateTime,
            ValueDateTime:        valueDateTime,
            TransactionInfo: bankadapter.TransInfo{
                Description:          tx.TransactionInfo.Description,
                TransactionReference: tx.TransactionInfo.TransactionReference,
            },
            CounterpartyName:    tx.CounterpartyName,
            CounterpartyAccount: tx.CounterpartyAccount,
            Category:            tx.Category,
        })
    }
    
    return transactions, nil
}

