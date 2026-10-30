package domain

import (
	"strings"
	"uuid"

	"github.com/sirnacou/file-ctrl/backend/internal/common"
)

var (
	ErrUserNotFound       = common.NewNotFoundError("user not found")
	ErrInvalidCredentials = common.NewUnauthorizedError("invalid username or password")
	ErrSetupAlreadyDone   = common.NewForbiddenError(
		"initial admin setup has already been completed",
	)
	ErrSessionNotFound = common.NewNotFoundError("session not found or expired")
	ErrInvalidUsername = common.NewValidationError(
		"username",
		"must be at least 3 characters",
	)
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
	Role         Role
}

func NewUser(username, passwordHash string, role Role) (*User, error) {
	username = strings.TrimSpace(username)
	if len(username) < 3 {
		return nil, common.JoinValidationErrors("failed to create user", ErrInvalidUsername)
	}

	return &User{
		ID:           uuid.NewV7(),
		Username:     username,
		PasswordHash: passwordHash,
		Role:         RoleAdmin,
	}, nil
}
