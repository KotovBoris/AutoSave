package sbank

import (
    "github.com/autosave/backend/internal/bankadapter"
    "github.com/autosave/backend/internal/bankadapter/vbank"
    "github.com/rs/zerolog"
)

// Config for SBank adapter
type Config struct {
    ClientID     string
    ClientSecret string
    BaseURL      string
    TeamID       string
    Logger       *zerolog.Logger
}

// Adapter for SBank - currently using VBank implementation as base
type Adapter struct {
    *vbank.Adapter
}

// NewAdapter creates new SBank adapter
func NewAdapter(cfg Config) *Adapter {
    // For now, SBank uses same implementation as VBank
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

// GetBankInfo returns SBank specific information
func (a *Adapter) GetBankInfo() bankadapter.BankInfo {
    return bankadapter.BankInfo{
        ID:          "sbank",
        Name:        "Smart Bank",
        BaseURL:     a.Adapter.BaseURL,
        DepositRate: 9.0,
    }
}
