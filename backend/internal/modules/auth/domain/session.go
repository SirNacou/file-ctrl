package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	SessionTTL  = 7 * 24 * time.Hour
	TokenLength = 32
)

type SessionToken struct {
	Raw       string
	Hash      string
	ExpiresAt time.Time
}

func NewSessionToken() (SessionToken, error) {
	b := make([]byte, TokenLength)
	if _, err := rand.Read(b); err != nil {
		return SessionToken{}, fmt.Errorf("failed to read random bytes: %w", err)
	}

	raw := hex.EncodeToString(b)
	return SessionToken{
		Raw:       raw,
		Hash:      hashSessionToken(raw),
		ExpiresAt: time.Now().UTC().Add(SessionTTL),
	}, nil
}

func hashSessionToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}
