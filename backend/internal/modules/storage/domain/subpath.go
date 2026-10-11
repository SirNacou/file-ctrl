package domain

import (
	"errors"
	"path"
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
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return SourceSubpath{}, ErrEmptySubpath
	}

	// 1. Invariant: No control characters or null bytes
	if strings.ContainsRune(trimmed, 0) || strings.ContainsAny(trimmed, "\r\n") {
		return SourceSubpath{}, ErrInvalidSubpathChars
	}

	// 2. Invariant: Normalized virtual path
	cleaned := path.Clean("/" + trimmed)
	sanitized := strings.Trim(cleaned, "/")

	if sanitized == "" || sanitized == "." {
		return SourceSubpath{}, ErrEmptySubpath
	}

	// 3. Invariant: Traversal prevention
	if strings.HasPrefix(sanitized, "..") {
		return SourceSubpath{}, ErrPathEscapesSandbox
	}

	// 4. Invariant: Single-level directory rule
	if !allowNested && strings.Contains(sanitized, "/") {
		return SourceSubpath{}, ErrNestedSubpathForbidden
	}

	return SourceSubpath{value: sanitized}, nil
}

func (s SourceSubpath) String() string {
	return s.value
}
