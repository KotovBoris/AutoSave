package banks

import (
    "fmt"
    
    "github.com/autosave/backend/internal/bankadapter"
    "github.com/autosave/backend/internal/bankadapter/vbank"
    "github.com/autosave/backend/internal/bankadapter/abank"
    "github.com/autosave/backend/internal/bankadapter/sbank"
    "github.com/autosave/backend/internal/config"
    "github.com/rs/zerolog"
)

// Factory creates bank adapters
type Factory struct {
    config *config.Config
    logger *zerolog.Logger
}

// NewFactory creates a new bank adapter factory
func NewFactory(cfg *config.Config, logger *zerolog.Logger) *Factory {
    return &Factory{
        config: cfg,
        logger: logger,
    }
}

// CreateAdapter creates a bank adapter for the specified bank
func (f *Factory) CreateAdapter(bankID string) (bankadapter.BankAdapter, error) {
    clientID, clientSecret, apiURL := f.config.GetBankConfig(bankID)
    
    if clientID == "" || clientSecret == "" || apiURL == "" {
        return nil, fmt.Errorf("bank configuration not found for %s", bankID)
    }
    
    switch bankID {
    case "vbank":
        return vbank.NewAdapter(vbank.Config{
            ClientID:     clientID,
            ClientSecret: clientSecret,
            BaseURL:      apiURL,
            TeamID:       f.config.TeamID,
            Logger:       f.logger,
        }), nil
        
    case "abank":
        return abank.NewAdapter(abank.Config{
            ClientID:     clientID,
            ClientSecret: clientSecret,
            BaseURL:      apiURL,
            TeamID:       f.config.TeamID,
            Logger:       f.logger,
        }), nil
        
    case "sbank":
        return sbank.NewAdapter(sbank.Config{
            ClientID:     clientID,
            ClientSecret: clientSecret,
            BaseURL:      apiURL,
            TeamID:       f.config.TeamID,
            Logger:       f.logger,
        }), nil
        
    default:
        return nil, fmt.Errorf("unknown bank: %s", bankID)
    }
}

// GetSupportedBanks returns list of supported bank IDs
func (f *Factory) GetSupportedBanks() []string {
    return []string{"vbank", "abank", "sbank"}
}

// ValidateBankID checks if bank ID is supported
func (f *Factory) ValidateBankID(bankID string) bool {
    switch bankID {
    case "vbank", "abank", "sbank":
        return true
    default:
        return false
    }
}
