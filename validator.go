package forge

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Validator provides input validation utilities
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateRequired checks if a value is not empty
func (v *Validator) ValidateRequired(value string, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	return nil
}

// ValidateEmail checks if a value is a valid email
func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email is required")
	}
	
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}
	return nil
}

// ValidateLength checks if a string has valid length
func (v *Validator) ValidateLength(value string, min, max int, fieldName string) error {
	length := len(strings.TrimSpace(value))
	if length < min {
		return fmt.Errorf("%s must be at least %d characters", fieldName, min)
	}
	if max > 0 && length > max {
		return fmt.Errorf("%s must be at most %d characters", fieldName, max)
	}
	return nil
}

// ValidateNumeric checks if a string is a valid number
func (v *Validator) ValidateNumeric(value string, fieldName string) error {
	if _, err := strconv.ParseFloat(value, 64); err != nil {
		return fmt.Errorf("%s must be a valid number", fieldName)
	}
	return nil
}

// ValidateInteger checks if a string is a valid integer
func (v *Validator) ValidateInteger(value string, fieldName string) error {
	if _, err := strconv.Atoi(value); err != nil {
		return fmt.Errorf("%s must be a valid integer", fieldName)
	}
	return nil
}

// ValidateRange checks if a number is within a range
func (v *Validator) ValidateRange(value string, min, max float64, fieldName string) error {
	num, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("%s must be a valid number", fieldName)
	}
	
	if num < min || num > max {
		return fmt.Errorf("%s must be between %.2f and %.2f", fieldName, min, max)
	}
	return nil
}

// ValidateAlphanumeric checks if a string contains only alphanumeric characters
func (v *Validator) ValidateAlphanumeric(value string, fieldName string) error {
	alphanumericRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !alphanumericRegex.MatchString(value) {
		return fmt.Errorf("%s must contain only letters and numbers", fieldName)
	}
	return nil
}

// Context validation helpers
func (c *Context) ValidateParam(key string, validator func(string) error) error {
	value, exists := c.Params[key]
	if !exists {
		return fmt.Errorf("parameter %s is missing", key)
	}
	return validator(value)
}

func (c *Context) ValidateQuery(key string, validator func(string) error) error {
	value, exists := c.Query[key]
	if !exists {
		return fmt.Errorf("query parameter %s is missing", key)
	}
	return validator(value)
}

// Validation middleware
func ValidationMiddleware() MiddlewareFunc {
	return func(c *Context) error {
		// Add validator to context
		c.Set("validator", NewValidator())
		return c.Next()
	}
}

// Helper to get validator from context
func (c *Context) GetValidator() *Validator {
	if v := c.Get("validator"); v != nil {
		return v.(*Validator)
	}
	return NewValidator()
}