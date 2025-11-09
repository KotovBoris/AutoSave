package vbank

import (
	"net/http"

	"github.com/autosave/backend/internal/bankadapter"
	"github.com/rs/zerolog"
)

// Adapter implements BankAdapter for VBank
type Adapter struct {
	*bankadapter.BaseAdapter
	config Config
	logger zerolog.Logger
}

// NewAdapter creates new VBank adapter
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
		logger: cfg.Logger.With().Str("bank", "vbank").Logger(),
	}
}

// GetBankInfo returns static bank information
func (a *Adapter) GetBankInfo() bankadapter.BankInfo {
	return bankadapter.BankInfo{
		ID:          "vbank",
		Name:        "Virtual Bank",
		BaseURL:     a.config.BaseURL,
		DepositRate: 8.0,
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
	return resp.StatusCode == http.StatusOK
}

