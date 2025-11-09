package abank

import (
    "github.com/autosave/backend/internal/bankadapter"
    "github.com/autosave/backend/internal/bankadapter/vbank"
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

// Adapter for ABank - currently using VBank implementation as base
type Adapter struct {
    *vbank.Adapter
}

// NewAdapter creates new ABank adapter
func NewAdapter(cfg Config) *Adapter {
    // For now, ABank uses same implementation as VBank
    // Just with different base URL and credentials
    vbankAdapter := vbank.NewAdapter(vbank.Config{
        ClientID:     cfg.ClientID,
        ClientSecret: cfg.ClientSecret,
        BaseURL:      cfg.BaseURL,
        TeamID:       cfg.TeamID,
        Logger:       cfg.Logger,
    })
    
    return &Adapter{
        Adapter: vbankAdapter,
    }
}

// GetBankInfo returns ABank specific information
func (a *Adapter) GetBankInfo() bankadapter.BankInfo {
    return bankadapter.BankInfo{
        ID:          "abank",
        Name:        "Awesome Bank",
        BaseURL:     a.Adapter.BaseURL,
        DepositRate: 7.5,
    }
}
