package abank

import (
    "github.com/autosave/backend/internal/bankadapter"
    "github.com/rs/zerolog"
)

// Config for ABank adapter
type Config struct {
    ClientID     string
    ClientSecret string
    BaseURL      string
    TeamID       string
    Logger       *zerolog.Logger
}

// Adapter implements BankAdapter for ABank
type Adapter struct {
    *bankadapter.BaseAdapter
    config Config
    logger zerolog.Logger
}

// NewAdapter creates new ABank adapter
func NewAdapter(cfg Config) bankadapter.BankAdapter {
    return &Adapter{
        BaseAdapter: bankadapter.NewBaseAdapter(
            cfg.ClientID,
            cfg.ClientSecret,
            cfg.BaseURL,
            cfg.TeamID,
            cfg.Logger,
        ),
        config: cfg,
        logger: cfg.Logger.With().Str("bank", "abank").Logger(),
    }
}

// GetBankInfo returns static bank information
func (a *Adapter) GetBankInfo() bankadapter.BankInfo {
    return bankadapter.BankInfo{
        ID:          "abank",
        Name:        "Awesome Bank",
        BaseURL:     a.config.BaseURL,
        DepositRate: 7.5,
    }
}

// IsHealthy checks if bank API is available
func (a *Adapter) IsHealthy() bool {
    resp, err := a.DoRequest("GET", "/health", nil, nil)
    if err != nil {
        a.logger.Error().Err(err).Msg("Health check failed")
        return false
    }
    defer resp.Body.Close()
    return resp.StatusCode == 200
}


