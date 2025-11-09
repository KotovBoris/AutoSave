package services

import (
	"context"
	"fmt"

	"github.com/autosave/backend/internal/models"
	"github.com/autosave/backend/internal/repository"
	"github.com/autosave/backend/pkg/jwt"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.UserRepository
	jwtUtil  *jwt.JWTUtil
	logger   *zerolog.Logger
}

func NewAuthService(userRepo repository.UserRepository, jwtUtil *jwt.JWTUtil, logger *zerolog.Logger) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtUtil:  jwtUtil,
		logger:   logger,
	}
}

// Register creates new user
func (s *AuthService) Register(ctx context.Context, req models.UserRegistration) (*models.User, string, error) {
	s.logger.Info().Str("email", req.Email).Msg("Registering new user")

	// Check if user already exists
	existing, _ := s.userRepo.GetByEmail(ctx, req.Email)
	if existing != nil {
		return nil, "", fmt.Errorf("user with email %s already exists", req.Email)
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to hash password")
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:            req.Email,
		PasswordHash:     string(hashedPassword),
		AutopilotEnabled: false,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error().Err(err).Msg("Failed to create user")
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.jwtUtil.GenerateToken(user.ID, user.Email)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to generate token")
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Info().Int("userId", user.ID).Msg("User registered successfully")

	return user, token, nil
}

// Login authenticates user
func (s *AuthService) Login(ctx context.Context, req models.UserLogin) (*models.User, string, error) {
	s.logger.Info().Str("email", req.Email).Msg("User login attempt")

	// Find user
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.Warn().Str("email", req.Email).Msg("User not found")
		return nil, "", fmt.Errorf("invalid email or password")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		s.logger.Warn().Str("email", req.Email).Msg("Invalid password")
		return nil, "", fmt.Errorf("invalid email or password")
	}

	// Generate JWT token
	token, err := s.jwtUtil.GenerateToken(user.ID, user.Email)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to generate token")
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Info().Int("userId", user.ID).Msg("User logged in successfully")

	return user, token, nil
}

// GetUser returns user by ID
func (s *AuthService) GetUser(ctx context.Context, userID int) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return user, nil
}

// UpdateAutopilot enables/disables autopilot
func (s *AuthService) UpdateAutopilot(ctx context.Context, userID int, enabled bool) error {
	s.logger.Info().Int("userId", userID).Bool("enabled", enabled).Msg("Updating autopilot setting")

	if err := s.userRepo.UpdateAutopilot(ctx, userID, enabled); err != nil {
		s.logger.Error().Err(err).Msg("Failed to update autopilot")
		return fmt.Errorf("failed to update autopilot: %w", err)
	}

	return nil
}
