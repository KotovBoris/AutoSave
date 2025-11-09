package services

import (
	"github.com/KotovBoris/AutoSave/backend/internal/repository"
	"github.com/rs/zerolog"
)

type OperationService struct {
	operationRepo repository.OperationRepository
	logger        *zerolog.Logger
}

func NewOperationService(
	operationRepo repository.OperationRepository,
	logger *zerolog.Logger,
) *OperationService {
	return &OperationService{
		operationRepo: operationRepo,
		logger:        logger,
	}
}

// Placeholder for operations service
// Will be implemented when needed

