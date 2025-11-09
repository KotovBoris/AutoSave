package abank

import (
    "encoding/json"
    "fmt"
    "net/url"
    "time"
    
    "github.com/autosave/backend/internal/bankadapter"
)

// GetProducts retrieves available products
func (a *Adapter) GetProducts(token string, productType string) ([]bankadapter.Product, error) {
    params := url.Values{}
    if productType != "" {
        params.Set("product_type", productType)
    }
    
    headers := a.GetAuthHeaders(token)
    fullURL := a.BaseURL + endpointProducts
    if len(params) > 0 {
        fullURL += "?" + params.Encode()
    }
    
    resp, err := a.DoRequest("GET", fullURL, headers, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get products: %w", err)
    }
    defer resp.Body.Close()
    
    var response struct {
        Products []struct {
            ProductID    string   `json:"productId"`
            ProductType  string   `json:"productType"`
            ProductName  string   `json:"productName"`
            Description  string   `json:"description"`
            InterestRate float64  `json:"interestRate"`
            MinAmount    float64  `json:"minAmount"`
            MaxAmount    float64  `json:"maxAmount"`
            TermMonths   []int    `json:"termMonths"`
            Currency     string   `json:"currency"`
            Features     []string `json:"features"`
        } `json:"products"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode products: %w", err)
    }
    
    products := make([]bankadapter.Product, 0, len(response.Products))
    for _, p := range response.Products {
        products = append(products, bankadapter.Product{
            ProductID:    p.ProductID,
            ProductType:  p.ProductType,
            ProductName:  p.ProductName,
            Description:  p.Description,
            InterestRate: p.InterestRate,
            MinAmount:    p.MinAmount,
            MaxAmount:    p.MaxAmount,
            TermMonths:   p.TermMonths,
            Currency:     p.Currency,
            Features:     p.Features,
        })
    }
    
    return products, nil
}

// GetProductDetails retrieves details of a specific product
func (a *Adapter) GetProductDetails(token, productID string) (*bankadapter.Product, error) {
    path := fmt.Sprintf(endpointProduct, productID)
    headers := a.GetAuthHeaders(token)
    
    resp, err := a.DoRequest("GET", path, headers, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get product details: %w", err)
    }
    defer resp.Body.Close()
    
    var product struct {
        ProductID    string   `json:"productId"`
        ProductType  string   `json:"productType"`
        ProductName  string   `json:"productName"`
        Description  string   `json:"description"`
        InterestRate float64  `json:"interestRate"`
        MinAmount    float64  `json:"minAmount"`
        MaxAmount    float64  `json:"maxAmount"`
        TermMonths   []int    `json:"termMonths"`
        Currency     string   `json:"currency"`
        Features     []string `json:"features"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
        return nil, fmt.Errorf("failed to decode product details: %w", err)
    }
    
    return &bankadapter.Product{
        ProductID:    product.ProductID,
        ProductType:  product.ProductType,
        ProductName:  product.ProductName,
        Description:  product.Description,
        InterestRate: product.InterestRate,
        MinAmount:    product.MinAmount,
        MaxAmount:    product.MaxAmount,
        TermMonths:   product.TermMonths,
        Currency:     product.Currency,
        Features:     product.Features,
    }, nil
}

// GetAgreements retrieves all product agreements
func (a *Adapter) GetAgreements(token, clientID, consentID, requestingBank string) ([]bankadapter.Agreement, error) {
    params := url.Values{}
    if clientID != "" {
        params.Set("client_id", clientID)
    }
    
    headers := a.GetConsentHeaders(token, consentID, requestingBank)
    headers["X-Product-Agreement-Consent-Id"] = consentID
    
    fullURL := a.BaseURL + endpointProductAgreements
    if len(params) > 0 {
        fullURL += "?" + params.Encode()
    }
    
    resp, err := a.DoRequest("GET", fullURL, headers, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get agreements: %w", err)
    }
    defer resp.Body.Close()
    
    var response struct {
        Agreements []struct {
            AgreementID     string  `json:"agreementId"`
            ProductID       string  `json:"productId"`
            ProductType     string  `json:"productType"`
            ProductName     string  `json:"productName"`
            Amount          float64 `json:"amount"`
            Currency        string  `json:"currency"`
            InterestRate    float64 `json:"interestRate"`
            TermMonths      int     `json:"termMonths"`
            Status          string  `json:"status"`
            OpenedDate      string  `json:"openedDate"`
            MaturityDate    string  `json:"maturityDate"`
            AccruedInterest float64 `json:"accruedInterest"`
        } `json:"agreements"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode agreements: %w", err)
    }
    
    agreements := make([]bankadapter.Agreement, 0, len(response.Agreements))
    for _, agr := range response.Agreements {
        openedDate, _ := time.Parse(time.RFC3339, agr.OpenedDate)
        maturityDate, _ := time.Parse(time.RFC3339, agr.MaturityDate)
        
        agreements = append(agreements, bankadapter.Agreement{
            AgreementID:     agr.AgreementID,
            ProductID:       agr.ProductID,
            ProductType:     agr.ProductType,
            ProductName:     agr.ProductName,
            Amount:          agr.Amount,
            Currency:        agr.Currency,
            InterestRate:    agr.InterestRate,
            TermMonths:      agr.TermMonths,
            Status:          agr.Status,
            OpenedDate:      openedDate,
            MaturityDate:    maturityDate,
            AccruedInterest: agr.AccruedInterest,
        })
    }
    
    return agreements, nil
}

// OpenDeposit opens a new deposit
func (a *Adapter) OpenDeposit(token, clientID, consentID, requestingBank string, request bankadapter.DepositRequest) (*bankadapter.Agreement, error) {
    params := url.Values{}
    if clientID != "" {
        params.Set("client_id", clientID)
    }
    
    headers := a.GetConsentHeaders(token, consentID, requestingBank)
    headers["X-Product-Agreement-Consent-Id"] = consentID
    
    body := map[string]interface{}{
        "product_id":        request.ProductID,
        "amount":            request.Amount,
        "term_months":       request.TermMonths,
        "source_account_id": request.SourceAccountID,
    }
    
    fullURL := a.BaseURL + endpointProductAgreements
    if len(params) > 0 {
        fullURL += "?" + params.Encode()
    }
    
    resp, err := a.DoRequest("POST", fullURL, headers, body)
    if err != nil {
        return nil, fmt.Errorf("failed to open deposit: %w", err)
    }
    defer resp.Body.Close()
    
    var response struct {
        AgreementID     string  `json:"agreementId"`
        ProductID       string  `json:"productId"`
        ProductType     string  `json:"productType"`
        ProductName     string  `json:"productName"`
        Amount          float64 `json:"amount"`
        Currency        string  `json:"currency"`
        InterestRate    float64 `json:"interestRate"`
        TermMonths      int     `json:"termMonths"`
        Status          string  `json:"status"`
        OpenedDate      string  `json:"openedDate"`
        MaturityDate    string  `json:"maturityDate"`
        AccountID       string  `json:"accountId"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode deposit response: %w", err)
    }
    
    openedDate, _ := time.Parse(time.RFC3339, response.OpenedDate)
    maturityDate, _ := time.Parse(time.RFC3339, response.MaturityDate)
    
    return &bankadapter.Agreement{
        AgreementID:  response.AgreementID,
        ProductID:    response.ProductID,
        ProductType:  response.ProductType,
        ProductName:  response.ProductName,
        Amount:       response.Amount,
        Currency:     response.Currency,
        InterestRate: response.InterestRate,
        TermMonths:   response.TermMonths,
        Status:       response.Status,
        OpenedDate:   openedDate,
        MaturityDate: maturityDate,
    }, nil
}

// CloseDeposit closes an existing deposit
func (a *Adapter) CloseDeposit(token, clientID, consentID, requestingBank, agreementID string) (*bankadapter.CloseDepositResponse, error) {
    path := fmt.Sprintf(endpointProductAgreement, agreementID)
    params := url.Values{}
    if clientID != "" {
        params.Set("client_id", clientID)
    }
    
    headers := a.GetConsentHeaders(token, consentID, requestingBank)
    headers["X-Product-Agreement-Consent-Id"] = consentID
    
    fullURL := a.BaseURL + path
    if len(params) > 0 {
        fullURL += "?" + params.Encode()
    }
    
    resp, err := a.DoRequest("DELETE", fullURL, headers, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to close deposit: %w", err)
    }
    defer resp.Body.Close()
    
    var response struct {
        AgreementID     string  `json:"agreementId"`
        ClosedAt        string  `json:"closedAt"`
        ReturnedAmount  float64 `json:"returnedAmount"`
        AccruedInterest float64 `json:"accruedInterest"`
        PenaltyAmount   float64 `json:"penaltyAmount"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode close deposit response: %w", err)
    }
    
    closedAt, _ := time.Parse(time.RFC3339, response.ClosedAt)
    
    return &bankadapter.CloseDepositResponse{
        AgreementID:     response.AgreementID,
        ClosedAt:        closedAt,
        ReturnedAmount:  response.ReturnedAmount,
        AccruedInterest: response.AccruedInterest,
        PenaltyAmount:   response.PenaltyAmount,
    }, nil
}

// GetAgreementDetails retrieves details of specific agreement
func (a *Adapter) GetAgreementDetails(token, clientID, consentID, requestingBank, agreementID string) (*bankadapter.Agreement, error) {
    path := fmt.Sprintf(endpointProductAgreement, agreementID)
    params := url.Values{}
    if clientID != "" {
        params.Set("client_id", clientID)
    }
    
    headers := a.GetConsentHeaders(token, consentID, requestingBank)
    headers["X-Product-Agreement-Consent-Id"] = consentID
    
    fullURL := a.BaseURL + path
    if len(params) > 0 {
        fullURL += "?" + params.Encode()
    }
    
    resp, err := a.DoRequest("GET", fullURL, headers, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to get agreement details: %w", err)
    }
    defer resp.Body.Close()
    
    var response struct {
        AgreementID     string  `json:"agreementId"`
        ProductID       string  `json:"productId"`
        ProductType     string  `json:"productType"`
        ProductName     string  `json:"productName"`
        Amount          float64 `json:"amount"`
        Currency        string  `json:"currency"`
        InterestRate    float64 `json:"interestRate"`
        TermMonths      int     `json:"termMonths"`
        Status          string  `json:"status"`
        OpenedDate      string  `json:"openedDate"`
        MaturityDate    string  `json:"maturityDate"`
        AccruedInterest float64 `json:"accruedInterest"`
    }
    
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return nil, fmt.Errorf("failed to decode agreement details: %w", err)
    }
    
    openedDate, _ := time.Parse(time.RFC3339, response.OpenedDate)
    maturityDate, _ := time.Parse(time.RFC3339, response.MaturityDate)
    
    return &bankadapter.Agreement{
        AgreementID:     response.AgreementID,
        ProductID:       response.ProductID,
        ProductType:     response.ProductType,
        ProductName:     response.ProductName,
        Amount:          response.Amount,
        Currency:        response.Currency,
        InterestRate:    response.InterestRate,
        TermMonths:      response.TermMonths,
        Status:          response.Status,
        OpenedDate:      openedDate,
        MaturityDate:    maturityDate,
        AccruedInterest: response.AccruedInterest,
    }, nil
}


