package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
	"uuid"

	"github.com/sirnacou/file-ctrl/backend/internal/common"
)

var (
	ErrInvalidSession = common.NewUnauthorizedError("Session expired or invalid")
)

const (
	SessionTTL  = 7 * 24 * time.Hour
	TokenLength = 32
)

// SessionToken is a Value Object holding the ephemeral raw secret for the client cookie
// and the SHA-256 hash used as the Session ID in persistence.
type SessionToken struct {
	Raw  string
	Hash string
}

func NewSessionToken() (SessionToken, error) {
	b := make([]byte, TokenLength)
	if _, err := rand.Read(b); err != nil {
		return SessionToken{}, fmt.Errorf("failed to read random bytes: %w", err)
	}

	raw := hex.EncodeToString(b)
	return SessionToken{
		Raw:  raw,
		Hash: HashSessionToken(raw),
	}, nil
}

func HashSessionToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}

// Session is the Domain Entity representing an active authenticated session.
type Session struct {
	ID        string // The SHA-256 token hash
	UserID    uuid.UUID
	ExpiresAt time.Time
	CreatedAt time.Time
}

// NewSession creates a new Session entity and its corresponding client SessionToken.
func NewSession(userID uuid.UUID) (*Session, SessionToken, error) {
	token, err := NewSessionToken()
	if err != nil {
		return nil, SessionToken{}, err
	}

	now := time.Now().UTC()
	session := &Session{
		ID:        token.Hash,
		UserID:    userID,
		ExpiresAt: now.Add(SessionTTL),
		CreatedAt: now,
	}

	return session, token, nil
}

// IsExpired enforces the domain rule for session validity.
func (s *Session) IsExpired() bool {
	return time.Now().UTC().After(s.ExpiresAt)
}
