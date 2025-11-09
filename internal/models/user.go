package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type User struct {
	ID               int       `db:"id" json:"id"`
	Email            string    `db:"email" json:"email"`
	PasswordHash     string    `db:"password_hash" json:"-"`
	AvgSalary        *float64  `db:"avg_salary" json:"avgSalary"`
	AvgExpenses      *float64  `db:"avg_expenses" json:"avgExpenses"`
	SavingsCapacity  *float64  `db:"savings_capacity" json:"savingsCapacity"`
	SalaryDates      IntArray  `db:"salary_dates" json:"salaryDates"`
	AutopilotEnabled bool      `db:"autopilot_enabled" json:"autopilotEnabled"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time `db:"updated_at" json:"updatedAt"`
}

// IntArray for PostgreSQL integer[] type
type IntArray []int

func (a IntArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "{}", nil
	}
	return fmt.Sprintf("{%v}", intArrayToString(a)), nil
}

func (a *IntArray) Scan(src interface{}) error {
	if src == nil {
		*a = IntArray{}
		return nil
	}

	switch s := src.(type) {
	case string:
		return scanIntArray(s, a)
	case []byte:
		return scanIntArray(string(s), a)
	}
	return fmt.Errorf("cannot scan %T into IntArray", src)
}

func intArrayToString(arr []int) string {
	if len(arr) == 0 {
		return ""
	}
	result := ""
	for i, v := range arr {
		if i > 0 {
			result += ","
		}
		result += fmt.Sprintf("%d", v)
	}
	return result
}

func scanIntArray(src string, a *IntArray) error {
	if src == "{}" || src == "" {
		*a = IntArray{}
		return nil
	}

	// Remove curly braces
	src = src[1 : len(src)-1]
	if src == "" {
		*a = IntArray{}
		return nil
	}

	// Parse integers
	var result []int
	for _, s := range splitComma(src) {
		var i int
		if _, err := fmt.Sscanf(s, "%d", &i); err != nil {
			return err
		}
		result = append(result, i)
	}
	*a = result
	return nil
}

func splitComma(s string) []string {
	if s == "" {
		return []string{}
	}
	var result []string
	start := 0
	for i, r := range s {
		if r == ',' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}

type UserRegistration struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID               int       `json:"id"`
	Email            string    `json:"email"`
	AvgSalary        *float64  `json:"avgSalary"`
	AvgExpenses      *float64  `json:"avgExpenses"`
	SavingsCapacity  *float64  `json:"savingsCapacity"`
	SalaryDates      []int     `json:"salaryDates"`
	AutopilotEnabled bool      `json:"autopilotEnabled"`
	CreatedAt        time.Time `json:"createdAt"`
}

func (u *User) ToResponse() UserResponse {
	salaryDates := []int{}
	if u.SalaryDates != nil {
		salaryDates = u.SalaryDates
	}

	return UserResponse{
		ID:               u.ID,
		Email:            u.Email,
		AvgSalary:        u.AvgSalary,
		AvgExpenses:      u.AvgExpenses,
		SavingsCapacity:  u.SavingsCapacity,
		SalaryDates:      salaryDates,
		AutopilotEnabled: u.AutopilotEnabled,
		CreatedAt:        u.CreatedAt,
	}
}
