package models

import (
	"errors"
	"fmt"
)

// Common errors for the orders service
var (
	// ErrNotFound is returned when a requested resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrInvalidStatus is returned when an invalid status transition is attempted
	ErrInvalidStatus = errors.New("invalid status transition")

	// ErrValidation is returned when input validation fails
	ErrValidation = errors.New("validation error")
)

// ValidationError represents a validation error with field-specific messages
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed for field '%s': %s", e.Field, e.Message)
}

// NotFoundError is returned when a resource is not found
type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s not found", e.Resource)
}

// StatusTransitionError is returned when an invalid status transition is attempted
type StatusTransitionError struct {
	CurrentStatus string
	NewStatus     string
}

func (e *StatusTransitionError) Error() string {
	return fmt.Sprintf("invalid status transition from %s to %s", e.CurrentStatus, e.NewStatus)
}

// ConflictError is returned when a resource conflict occurs
type ConflictError struct {
	Resource string
	Message  string
}

func (e *ConflictError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("conflict with %s: %s", e.Resource, e.Message)
	}
	return fmt.Sprintf("conflict with %s", e.Resource)
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	var ve *ValidationError
	return errors.As(err, &ve)
}

// IsNotFoundError checks if an error is a not found error
func IsNotFoundError(err error) bool {
	var nfe *NotFoundError
	return errors.Is(err, ErrNotFound) || errors.As(err, &nfe)
}

// IsStatusTransitionError checks if an error is a status transition error
func IsStatusTransitionError(err error) bool {
	var ste *StatusTransitionError
	return errors.Is(err, ErrInvalidStatus) || errors.As(err, &ste)
}

// IsConflictError checks if an error is a conflict error
func IsConflictError(err error) bool {
	var ce *ConflictError
	return errors.As(err, &ce)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// NewNotFoundError creates a new not found error
func NewNotFoundError(resource, id string) error {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

// NewStatusTransitionError creates a new status transition error
func NewStatusTransitionError(currentStatus, newStatus string) error {
	return &StatusTransitionError{
		CurrentStatus: currentStatus,
		NewStatus:     newStatus,
	}
}

// NewConflictError creates a new conflict error
func NewConflictError(resource, message string) error {
	return &ConflictError{
		Resource: resource,
		Message:  message,
	}
}
