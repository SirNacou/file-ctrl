package domain

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	ErrEmptyPath              = errors.New("storage path cannot be empty")
	ErrInvalidSubpathChars    = errors.New("subpath contains illegal characters or null bytes")
	ErrNestedSubpathForbidden = errors.New("subpath must be a single top-level directory")
	ErrPathEscapesSandbox     = errors.New("subpath attempts to escape sandbox")
)

// StoragePath is a pure domain value object representing a safe, normalized subpath.
type StoragePath struct {
	value string
}

func NewSourcePathWithDefault(raw string) (StoragePath, error) {
	return NewSourcePath(raw, false)
}

func NewSourcePath(raw string, allowNested bool) (StoragePath, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return StoragePath{}, ErrEmptyPath
	}

	if err := validateChars(raw); err != nil {
		return StoragePath{}, err
	}

	segments := strings.Split(raw, "/")
	if err := validateSegments(segments); err != nil {
		return StoragePath{}, err
	}

	if !filepath.IsLocal(filepath.FromSlash(raw)) {
		return StoragePath{}, ErrPathEscapesSandbox
	}

	if !allowNested && len(segments) > 1 {
		return StoragePath{}, ErrNestedSubpathForbidden
	}

	return StoragePath{value: raw}, nil
}

// validateChars rejects control characters and backslashes.
func validateChars(s string) error {
	for _, r := range s {
		if r < 0x20 || r == 0x7f || r == '\\' {
			return ErrInvalidSubpathChars
		}
	}

	return nil
}

// validateSegments rejects empty, "." and ".." segments, which also
// covers leading/trailing/double slashes.
func validateSegments(segments []string) error {
	for _, s := range segments {
		if s == "" || s == "." || s == ".." {
			return ErrPathEscapesSandbox
		}
	}

	return nil
}

func (s StoragePath) String() string {
	return s.value
}
