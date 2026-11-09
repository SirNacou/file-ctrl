package domain

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewEntryPath(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{name: "empty is root", raw: "", want: "."},
		{name: "whitespace is root", raw: "  ", want: "."},
		{name: "dot is root", raw: ".", want: "."},
		{name: "slashes are root", raw: "///", want: "."},
		{name: "trims whitespace", raw: "  docs/file.txt  ", want: "docs/file.txt"},
		{name: "leading slash is made relative", raw: "/docs/file.txt", want: "docs/file.txt"},
		{name: "cleans local dot segments", raw: "docs/./file.txt", want: "docs/file.txt"},
		{name: "cleans local parent segments", raw: "docs/../file.txt", want: "file.txt"},
		{name: "rejects parent traversal", raw: "../outside", wantErr: true},
		{name: "rejects traversal above root", raw: "docs/../../outside", wantErr: true},
		{name: "rejects null byte", raw: "docs/\x00file", wantErr: true},
		{name: "rejects backslash", raw: `docs\file.txt`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEntryPath(tt.raw)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("NewEntryPath(%q) error = nil, want error", tt.raw)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewEntryPath(%q) error = %v", tt.raw, err)
			}
			if got.String() != tt.want {
				t.Errorf("NewEntryPath(%q) = %q, want %q", tt.raw, got.String(), tt.want)
			}
		})
	}
}

func TestEntryPathRepresentations(t *testing.T) {
	if !Root.IsRoot() {
		t.Fatal("Root.IsRoot() = false, want true")
	}
	if got := Root.String(); got != "." {
		t.Errorf("Root.String() = %q, want %q", got, ".")
	}
	if got := Root.FS(); got != filepath.FromSlash(".") {
		t.Errorf("Root.FS() = %q, want %q", got, filepath.FromSlash("."))
	}
	if got := Root.Display(); got != "" {
		t.Errorf("Root.Display() = %q, want empty string", got)
	}

	entry, err := NewEntryPath("docs/2026/file.txt")
	if err != nil {
		t.Fatalf("NewEntryPath() error = %v", err)
	}
	if entry.IsRoot() {
		t.Error("nested entry IsRoot() = true, want false")
	}
	if got, want := entry.FS(), filepath.FromSlash("docs/2026/file.txt"); got != want {
		t.Errorf("nested entry FS() = %q, want %q", got, want)
	}
	if got, want := entry.Display(), "docs/2026/file.txt"; got != want {
		t.Errorf("nested entry Display() = %q, want %q", got, want)
	}
}

func TestEntryPathJoin(t *testing.T) {
	tests := []struct {
		name    string
		parent  string
		rel     string
		want    string
		wantErr bool
	}{
		{name: "join from root", parent: ".", rel: "docs/file.txt", want: "docs/file.txt"},
		{name: "join nested path", parent: "docs", rel: "2026/file.txt", want: "docs/2026/file.txt"},
		{name: "reject traversal", parent: "docs", rel: "../outside", wantErr: true},
		{name: "reject dot child", parent: "docs", rel: ".", wantErr: true},
		{name: "reject empty child", parent: "docs", rel: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent, err := NewEntryPath(tt.parent)
			if err != nil {
				t.Fatalf("NewEntryPath(%q): %v", tt.parent, err)
			}
			got, err := parent.Join(tt.rel)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Join(%q) error = nil, want error", tt.rel)
				}
				return
			}
			if err != nil {
				t.Fatalf("Join(%q) error = %v", tt.rel, err)
			}
			if got.String() != tt.want {
				t.Errorf("Join(%q) = %q, want %q", tt.rel, got.String(), tt.want)
			}
		})
	}
}

func TestEntryPathAncestors(t *testing.T) {
	tests := []struct {
		name string
		path string
		want []string
	}{
		{name: "root has no ancestors", path: ".", want: nil},
		{name: "top-level path has no ancestors", path: "file.txt", want: []string{}},
		{name: "nested path returns parents", path: "docs/2026/file.txt", want: []string{"docs", "docs/2026"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry, err := NewEntryPath(tt.path)
			if err != nil {
				t.Fatalf("NewEntryPath(%q): %v", tt.path, err)
			}
			ancestors := entry.Ancestors()
			var got []string
			if ancestors != nil {
				got = make([]string, len(ancestors))
				for i, ancestor := range ancestors {
					got[i] = ancestor.String()
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Ancestors() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
