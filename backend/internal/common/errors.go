package common

import (
	"errors"
	"fmt"
	"sort"
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
		fields := make([]string, 0, len(e.ValidationMap))
		for f := range e.ValidationMap {
			fields = append(fields, f)
		}
		sort.Strings(fields)

		details := make([]string, 0, len(fields))
		for _, f := range fields {
			details = append(
				details,
				fmt.Sprintf("%s: %s", f, strings.Join(e.ValidationMap[f], ", ")),
			)
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

const defaultValidationMessage = "validation failed"

// ValidationErrors accumulates field errors. Zero value is ready to use.
type ValidationErrors struct {
	message string
	fields  map[string][]string
}

func NewValidation() *ValidationErrors {
	return &ValidationErrors{fields: make(map[string][]string)}
}

func (v *ValidationErrors) WithMessage(msg string) *ValidationErrors {
	v.message = msg
	return v
}

func (v *ValidationErrors) Add(field, message string) *ValidationErrors {
	if v.fields == nil {
		v.fields = make(map[string][]string)
	}
	v.fields[field] = append(v.fields[field], message)
	return v
}

func (v *ValidationErrors) Addf(field, format string, args ...any) *ValidationErrors {
	return v.Add(field, fmt.Sprintf(format, args...))
}

// Merge folds in any validation error (e.g. from a nested value object).
// Non-validation and nil errors are ignored, so callers can pass results directly.
func (v *ValidationErrors) Merge(err error) *ValidationErrors {
	var de *DomainError
	if errors.As(err, &de) && de.Type == TypeValidation {
		for field, msgs := range de.ValidationMap {
			for _, m := range msgs {
				v.Add(field, m)
			}
		}
	}
	return v
}

func (v *ValidationErrors) HasErrors() bool { return len(v.fields) > 0 }

func (v *ValidationErrors) Err() error {
	if !v.HasErrors() {
		return nil
	}
	msg := v.message
	if msg == "" {
		msg = defaultValidationMessage
	}
	return &DomainError{
		Type:          TypeValidation,
		Message:       msg,
		ValidationMap: v.fields,
	}
}

// Single-error shortcut, now just a thin wrapper.
func NewValidationError(field, message string) error {
	return NewValidation().Add(field, message).Err()
}
