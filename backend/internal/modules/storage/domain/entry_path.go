package domain

import (
	"path"
	"path/filepath"
	"strings"
)

type EntryPath struct{ value string }

func NewEntryPath(raw string) (EntryPath, error) {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "/")
	if raw == "" {
		return EntryPath{value: "."}, nil
	}
	if strings.ContainsRune(raw, 0) || strings.Contains(raw, `\`) {
		return EntryPath{}, ErrInvalidPath
	}
	if !filepath.IsLocal(filepath.FromSlash(raw)) {
		return EntryPath{}, ErrInvalidPath
	}
	return EntryPath{value: path.Clean(raw)}, nil
}

func (p EntryPath) FS() string {
	return filepath.FromSlash(p.value)
}

func (p EntryPath) Display() string {
	if p.value == "." {
		return ""
	}
	return p.value
}

func (p EntryPath) Child(name string) string {
	return path.Join(p.value, name)
}

func (p EntryPath) ChildPath(name string) (EntryPath, error) {
	if name == "" || name == "." || strings.ContainsAny(name, `/\`) {
		return EntryPath{}, ErrInvalidPath
	}

	return NewEntryPath(p.Child(name))
}
