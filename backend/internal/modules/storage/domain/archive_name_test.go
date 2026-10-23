package domain

import (
	"errors"
	"testing"
)

func mustSource(t *testing.T, s string) SourcePath {
	t.Helper()
	sp, err := NewSourcePath(s, true)
	if err != nil {
		t.Fatalf("NewSourcePath(%q): %v", s, err)
	}
	return sp
}

func mustEntries(t *testing.T, ps ...string) []EntryPath {
	t.Helper()
	out := make([]EntryPath, 0, len(ps))
	for _, p := range ps {
		e, err := NewEntryPath(p)
		if err != nil {
			t.Fatalf("NewEntryPath(%q): %v", p, err)
		}
		out = append(out, e)
	}
	return out
}

func TestNewArchiveName(t *testing.T) {
	const src = "my-source"

	tests := []struct {
		name    string
		paths   []string
		want    string
		wantErr error
	}{
		// errors
		{name: "nil selection", paths: nil, wantErr: ErrEmptySelection},
		{name: "empty selection", paths: []string{}, wantErr: ErrEmptySelection},

		// single entry
		{
			name:  "single file uses file name",
			paths: []string{"docs/report.pdf"},
			want:  "report.pdf.zip",
		},
		{name: "single file at root", paths: []string{"readme.md"}, want: "readme.md.zip"},
		{name: "single directory uses dir name", paths: []string{"photos/2024"}, want: "2024.zip"},
		{
			name:  "single entry is dot (root) falls back to source",
			paths: []string{"."},
			want:  src + ".zip",
		},
		{
			name:  "single entry is slash (root) falls back to source",
			paths: []string{"/"},
			want:  src + ".zip",
		},
		{
			name:  "single entry is empty (root) falls back to source",
			paths: []string{""},
			want:  src + ".zip",
		},

		// multiple entries, common parent
		{name: "same directory", paths: []string{"a/b/1.txt", "a/b/2.txt"}, want: "b.zip"},
		{
			name:  "three files same directory",
			paths: []string{"a/b/1.txt", "a/b/2.txt", "a/b/3.txt"},
			want:  "b.zip",
		},
		{
			name:  "common parent is shallower than files",
			paths: []string{"a/b/c/1.txt", "a/b/d/2.txt"},
			want:  "b.zip",
		},
		{name: "file and nested file", paths: []string{"a/1.txt", "a/b/2.txt"}, want: "a.zip"},
		{
			name:  "deeply nested shared prefix",
			paths: []string{"x/y/z/p/1", "x/y/z/q/2"},
			want:  "z.zip",
		},
		{
			name:  "common parent narrows across 3 paths",
			paths: []string{"a/b/c/1", "a/b/c/2", "a/b/d/3"},
			want:  "b.zip",
		},
		{
			name:  "leading slash is normalized",
			paths: []string{"/a/x.txt", "/a/y.txt"},
			want:  "a.zip",
		},
		{
			name:  "mixed leading slash and relative",
			paths: []string{"/a/x.txt", "a/y.txt"},
			want:  "a.zip",
		},
		{name: "directories as selection", paths: []string{"a/b/c", "a/b/d"}, want: "b.zip"},

		// multiple entries, no common parent -> source name
		{
			name:  "different top-level dirs",
			paths: []string{"a/1.txt", "b/2.txt"},
			want:  src + ".zip",
		},
		{name: "files all at root", paths: []string{"x.txt", "y.txt"}, want: src + ".zip"},
		{
			name:  "root files with leading slash",
			paths: []string{"/x.txt", "/y.txt"},
			want:  src + ".zip",
		},
		{
			name:  "root file and nested file",
			paths: []string{"x.txt", "a/y.txt"},
			want:  src + ".zip",
		},
		{
			name:  "segment compare, not string prefix",
			paths: []string{"ab/x.txt", "a/y.txt"},
			want:  src + ".zip",
		},
		{
			name:  "no common parent among three",
			paths: []string{"a/b/1", "a/b/2", "c/3"},
			want:  src + ".zip",
		},
		{name: "root entry plus other entry", paths: []string{".", "a/b"}, want: src + ".zip"},

		// duplicates
		{name: "duplicate paths", paths: []string{"a/b/1.txt", "a/b/1.txt"}, want: "b.zip"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewArchiveName(mustSource(t, src), mustEntries(t, tt.paths...))

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				if got.String() != "" {
					t.Errorf("expected zero ArchiveName on error, got %q", got.String())
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tt.want {
				t.Errorf("got %q, want %q", got.String(), tt.want)
			}
		})
	}
}

func TestCommonParent(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		want  string
	}{
		{"single path returns its dir", []string{"a/b/c.txt"}, "a/b"},
		{"same dir", []string{"a/b/1", "a/b/2"}, "a/b"},
		{"diverging", []string{"a/b/1", "a/c/2"}, "a"},
		{"nothing in common", []string{"a/1", "b/2"}, ""},
		{"root files", []string{"x", "y"}, "."},
		{"order independent", []string{"a/b/c/1", "a/b/2"}, "a/b"},
		{"order independent reversed", []string{"a/b/2", "a/b/c/1"}, "a/b"},
		{"segment not string prefix", []string{"ab/x", "a/y"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := commonParent(mustEntries(t, tt.paths...)); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBaseOr(t *testing.T) {
	src := mustSource(t, "my-source")

	tests := []struct {
		in   string
		want string
	}{
		{"a/b/c", "c"},
		{"file.txt", "file.txt"},
		{"/a/b", "b"},
		{".", "my-source"},
		{"/", "my-source"},
		{"", "my-source"},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if got := baseOr(src, tt.in); got != tt.want {
				t.Errorf("baseOr(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
