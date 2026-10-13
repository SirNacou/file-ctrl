package domain_test

import (
	"errors"
	"testing"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

func TestNewSourceSubpath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		wantErr     error
		name        string
		raw         string
		want        string
		allowNested bool
	}{
		// valid
		{name: "single segment", raw: "docs", want: "docs"},
		{name: "trims whitespace", raw: "  docs  ", want: "docs"},
		{name: "dash and underscore", raw: "my-files_2024", want: "my-files_2024"},
		{name: "dots in name are fine", raw: "..backup", want: "..backup"},
		{name: "nested allowed", raw: "a/b/c", allowNested: true, want: "a/b/c"},

		// empty
		{name: "empty", raw: "", wantErr: domain.ErrEmptySubpath},
		{name: "only spaces", raw: "   ", wantErr: domain.ErrEmptySubpath},

		// traversal and non-canonical shapes
		{name: "double dot", raw: "..", wantErr: domain.ErrPathEscapesSandbox},
		{name: "leading double dot", raw: "../etc", wantErr: domain.ErrPathEscapesSandbox},
		{
			name:        "middle double dot",
			raw:         "a/../b",
			allowNested: true,
			wantErr:     domain.ErrPathEscapesSandbox,
		},
		{name: "single dot", raw: ".", wantErr: domain.ErrPathEscapesSandbox},
		{
			name:        "dot segment",
			raw:         "a/./b",
			allowNested: true,
			wantErr:     domain.ErrPathEscapesSandbox,
		},
		{name: "absolute", raw: "/etc", wantErr: domain.ErrPathEscapesSandbox},
		{name: "trailing slash", raw: "docs/", wantErr: domain.ErrPathEscapesSandbox},
		{
			name:        "double slash",
			raw:         "a//b",
			allowNested: true,
			wantErr:     domain.ErrPathEscapesSandbox,
		},

		// illegal characters
		{name: "backslash", raw: `a\b`, wantErr: domain.ErrInvalidSubpathChars},
		{name: "backslash traversal", raw: `..\..\etc`, wantErr: domain.ErrInvalidSubpathChars},
		{name: "null byte", raw: "a\x00b", wantErr: domain.ErrInvalidSubpathChars},
		{name: "newline", raw: "a\nb", wantErr: domain.ErrInvalidSubpathChars},
		{name: "tab", raw: "a\tb", wantErr: domain.ErrInvalidSubpathChars},
		{name: "delete char", raw: "a\x7fb", wantErr: domain.ErrInvalidSubpathChars},

		// nesting rule
		{
			name:        "nested forbidden",
			raw:         "a/b",
			allowNested: false,
			wantErr:     domain.ErrNestedSubpathForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := domain.NewSourceSubpath(tt.raw, tt.allowNested)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("NewSourceSubpath(%q) error = %v, want %v", tt.raw, err, tt.wantErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("NewSourceSubpath(%q) unexpected error: %v", tt.raw, err)
			}

			if got.String() != tt.want {
				t.Fatalf("NewSourceSubpath(%q) = %q, want %q", tt.raw, got.String(), tt.want)
			}
		})
	}
}
