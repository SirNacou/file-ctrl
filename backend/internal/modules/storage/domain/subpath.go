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

// SourceSubpath is a pure domain value object representing a safe, normalized subpath.
type SourceSubpath struct {
	value string
}

func NewSourceSubpath(raw string, allowNested bool) (SourceSubpath, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return SourceSubpath{}, ErrEmptySubpath
	}

	if err := validateChars(raw); err != nil {
		return SourceSubpath{}, err
	}

	segments := strings.Split(raw, "/")
	if err := validateSegments(segments); err != nil {
		return SourceSubpath{}, err
	}

	if !filepath.IsLocal(filepath.FromSlash(raw)) {
		return SourceSubpath{}, ErrPathEscapesSandbox
	}

	if !allowNested && len(segments) > 1 {
		return SourceSubpath{}, ErrNestedSubpathForbidden
	}

	return SourceSubpath{value: raw}, nil
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

func (s SourceSubpath) String() string {
	return s.value
}
