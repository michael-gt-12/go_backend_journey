package day2

import "fmt"

// ValidationError is a custom error type. Any type with an Error() string
// method satisfies the built-in error interface.
type ValidationError struct {
	Field   string
	Message string
}

// Error makes *ValidationError satisfy the error interface.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %q: %s", e.Field, e.Message)
}

// ValidateAge returns a *ValidationError when age is outside 0..150,
// and nil when it is valid.
func ValidateAge(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age", Message: "must not be negative"}
	}
	if age > 150 {
		return &ValidationError{Field: "age", Message: "must not be greater than 150"}
	}
	return nil
}
