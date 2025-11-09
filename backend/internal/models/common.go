package models

// Pagination for list endpoints
type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

func NewPagination(page, limit, total int) Pagination {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	totalPages := (total + limit - 1) / limit

	return Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}

// ErrorResponse for API errors
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Common error codes
const (
	ErrCodeValidation        = "VALIDATION_ERROR"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeForbidden         = "FORBIDDEN"
	ErrCodeNotFound          = "NOT_FOUND"
	ErrCodeConflict          = "CONFLICT"
	ErrCodeInternal          = "INTERNAL_ERROR"
	ErrCodeBankUnavailable   = "BANK_UNAVAILABLE"
	ErrCodeInsufficientFunds = "INSUFFICIENT_FUNDS"
	ErrCodeConsentRequired   = "CONSENT_REQUIRED"
)

// Auth response types
type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// Settings types
type UpdateProfileRequest struct {
	AvgSalary   *float64 `json:"avgSalary,omitempty"`
	SalaryDates []int    `json:"salaryDates,omitempty"`
}

type UpdateAutopilotRequest struct {
	Enabled bool `json:"enabled"`
}

type AutopilotResponse struct {
	AutopilotEnabled bool `json:"autopilotEnabled"`
}

// Success response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

