package domain

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	ErrEmptySubpath           = errors.New("storage subpath cannot be empty")
	ErrInvalidSubpathChars    = errors.New("subpath contains illegal characters or null bytes")
	ErrNestedSubpathForbidden = errors.New("subpath must be a single top-level directory")
	ErrPathEscapesSandbox     = errors.New("subpath attempts to escape sandbox")
)

// StorageSubpath is a pure domain value object representing a safe, normalized subpath.
type StorageSubpath struct {
	value string
}

func NewSourceSubpath(raw string, allowNested bool) (StorageSubpath, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return StorageSubpath{}, ErrEmptySubpath
	}

	if err := validateChars(raw); err != nil {
		return StorageSubpath{}, err
	}

	segments := strings.Split(raw, "/")
	if err := validateSegments(segments); err != nil {
		return StorageSubpath{}, err
	}

	if !filepath.IsLocal(filepath.FromSlash(raw)) {
		return StorageSubpath{}, ErrPathEscapesSandbox
	}

	if !allowNested && len(segments) > 1 {
		return StorageSubpath{}, ErrNestedSubpathForbidden
	}

	return StorageSubpath{value: raw}, nil
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

func (s StorageSubpath) String() string {
	return s.value
}
