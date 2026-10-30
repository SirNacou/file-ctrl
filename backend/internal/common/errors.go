package common

import (
	"fmt"
	"strings"
)

type DomainErrorType string

const (
	TypeValidation   DomainErrorType = "VALIDATION_FAILED"
	TypeNotFound     DomainErrorType = "NOT_FOUND"
	TypeUnauthorized DomainErrorType = "UNAUTHORIZED"
	TypeForbidden    DomainErrorType = "FORBIDDEN"
	TypeFailure      DomainErrorType = "INTERNAL_FAILURE"
)

type DomainError struct {
	Type          DomainErrorType     `json:"type"`
	Message       string              `json:"message"`
	ValidationMap map[string][]string `json:"validation_map"`
	Err           error               `json:"-"`
}

func (e *DomainError) Error() string {
	if e.Type == TypeValidation && len(e.ValidationMap) > 0 {
		var details []string
		for field, msgs := range e.ValidationMap {
			details = append(details, fmt.Sprintf("%s: %s", field, strings.Join(msgs, ", ")))
		}
		return fmt.Sprintf("[%s] %s (%s)", e.Type, e.Message, strings.Join(details, "; "))
	}

	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// =========================================================================
// Constructors & Builders
// =========================================================================

func NewNotFoundError(message string) error {
	return &DomainError{
		Type:    TypeNotFound,
		Message: message,
	}
}

func NewUnauthorizedError(message string) error {
	return &DomainError{
		Type:    TypeUnauthorized,
		Message: message,
	}
}

func NewForbiddenError(message string) error {
	return &DomainError{
		Type:    TypeForbidden,
		Message: message,
	}
}

func NewFailureError(message string, err error) error {
	return &DomainError{
		Type:    TypeFailure,
		Message: message,
		Err:     err,
	}
}

func NewValidationError(field, message string) *DomainError {
	return &DomainError{
		Type:    TypeValidation,
		Message: "validation failed",
		ValidationMap: map[string][]string{
			field: {message},
		},
	}
}

func JoinValidationErrors(summaryMessage string, errs ...*DomainError) error {
	var combinedMap map[string][]string
	hasErrors := false

	for _, err := range errs {
		if err == nil || err.Type != TypeValidation {
			continue
		}
		if combinedMap == nil {
			combinedMap = make(map[string][]string)
		}
		hasErrors = true

		// Merge the fields safely
		for field, messages := range err.ValidationMap {
			combinedMap[field] = append(combinedMap[field], messages...)
		}
	}

	if !hasErrors {
		return nil
	}

	return &DomainError{
		Type:          TypeValidation,
		Message:       summaryMessage,
		ValidationMap: combinedMap,
	}
}
