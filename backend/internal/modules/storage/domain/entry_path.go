package domain

import (
	"path"
	"path/filepath"
	"strings"
)

type EntryPath struct{ value string }

// Root is the base directory.
var Root = EntryPath{value: "."}

func NewEntryPath(raw string) (EntryPath, error) {
	raw = strings.TrimLeft(strings.TrimSpace(raw), "/")
	if raw == "" {
		return Root, nil
	}
	if strings.ContainsRune(raw, 0) || strings.Contains(raw, `\`) {
		return EntryPath{}, ErrInvalidPath
	}
	if !filepath.IsLocal(filepath.FromSlash(raw)) {
		return EntryPath{}, ErrInvalidPath
	}
	return EntryPath{value: path.Clean(raw)}, nil
}

func (p EntryPath) IsRoot() bool { return p.value == "." }

func (p EntryPath) FS() string { return filepath.FromSlash(p.value) }

func (p EntryPath) String() string { return p.value }

func (p EntryPath) Display() string {
	if p.IsRoot() {
		return ""
	}
	return p.value
}

// Join appends a relative path (a name or a nested path like "a/b/c.txt").
// The relative part is validated on its own, so it can never escape p.
func (p EntryPath) Join(rel string) (EntryPath, error) {
	child, err := NewEntryPath(rel)
	if err != nil || child.IsRoot() {
		return EntryPath{}, ErrInvalidPath
	}
	return NewEntryPath(path.Join(p.value, child.value))
}

// Ancestors returns every parent directory, outermost first, excluding root
// and the path itself. "a/b/c.txt" -> ["a", "a/b"].
func (p EntryPath) Ancestors() []EntryPath {
	if p.IsRoot() {
		return nil
	}
	parts := strings.Split(p.value, "/")
	out := make([]EntryPath, 0, len(parts)-1)
	for i := 1; i < len(parts); i++ {
		out = append(out, EntryPath{value: strings.Join(parts[:i], "/")})
	}
	return out
}
