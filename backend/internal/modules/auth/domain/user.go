package domain

import (
	"strings"
	"unicode/utf8"

	"github.com/sirnacou/file-ctrl/backend/internal/common"
	"uuid"
)

// Sentinel errors (non-validation, safe to share).
var (
	ErrUserNotFound       = common.NewNotFoundError("user not found")
	ErrInvalidCredentials = common.NewUnauthorizedError("invalid username or password")
	ErrSetupAlreadyDone   = common.NewForbiddenError(
		"initial admin setup has already been completed",
	)
	ErrSessionNotFound = common.NewNotFoundError("session not found or expired")
)

// Validation field names and messages. These are constants, not error vars,
// because validation errors hold a mutable map and must be built per call.
const (
	FieldUsername   = "username"
	FieldProvider   = "provider"
	FieldProviderID = "provider_id"

	MsgMin3Characters = "must be at least 3 characters"
	MsgRequired       = "must not be empty"

	minUsernameLen = 3
)

type Provider string

const (
	ProviderLocal Provider = "local"
	ProviderOIDC  Provider = "oidc"
)

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash *string
	Role         Role
	Provider     Provider
	ProviderID   *string
}

func NewUser(username string, passwordHash string, role Role) (*User, error) {
	username = strings.TrimSpace(username)

	if utf8.RuneCountInString(username) < minUsernameLen {
		return nil, common.NewValidationError(FieldUsername, MsgMin3Characters)
	}

	return &User{
		ID:           uuid.NewV7(),
		Username:     username,
		PasswordHash: &passwordHash,
		Role:         role,
		Provider:     "local",
	}, nil
}

func NewFederatedUser(username, providerID string, provider Provider, role Role) (*User, error) {
	username = strings.TrimSpace(username)
	providerID = strings.TrimSpace(providerID)

	v := common.NewValidation().WithMessage("invalid federated user")

	if utf8.RuneCountInString(username) < minUsernameLen {
		v.Add(FieldUsername, MsgMin3Characters)
	}
	if provider == "" {
		v.Add(FieldProvider, MsgRequired)
	}
	if providerID == "" {
		v.Add(FieldProviderID, MsgRequired)
	}

	if err := v.Err(); err != nil {
		return nil, err
	}

	return &User{
		ID:         uuid.NewV7(),
		Username:   username,
		Role:       role,
		Provider:   provider,
		ProviderID: &providerID,
	}, nil
}
