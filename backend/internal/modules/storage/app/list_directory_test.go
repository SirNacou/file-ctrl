package app

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// normalizeUserPath
// ---------------------------------------------------------------------------

func TestNormalizeUserPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
	}{
		{name: "empty is root", raw: "", want: "."},
		{name: "slash is root", raw: "/", want: "."},
		{name: "dot is root", raw: ".", want: "."},
		{name: "simple", raw: "docs", want: "docs"},
		{name: "leading slash allowed", raw: "/docs/a", want: "docs/a"},
		{name: "dot prefix cleaned", raw: "./docs", want: "docs"},
		{name: "inner double dot stays inside", raw: "a/../b", want: "b"},
		{name: "trailing slash cleaned", raw: "docs/", want: "docs"},

		{name: "parent", raw: "..", wantErr: true},
		{name: "parent file", raw: "../secret", wantErr: true},
		{name: "slash then parent", raw: "/../secret", wantErr: true},
		{name: "escapes after descending", raw: "a/../../secret", wantErr: true},
		{name: "double leading slash", raw: "//etc/passwd", wantErr: true},
		{name: "backslash", raw: `a\b`, wantErr: true},
		{name: "null byte", raw: "a\x00b", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeUserPath(tt.raw)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("normalizeUserPath(%q) = %q, want error", tt.raw, got)
				}

				return
			}

			if err != nil {
				t.Fatalf("normalizeUserPath(%q) unexpected error: %v", tt.raw, err)
			}

			if got != tt.want {
				t.Fatalf("normalizeUserPath(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}

func TestDisplayPath(t *testing.T) {
	t.Parallel()

	if got := displayPath("."); got != "" {
		t.Fatalf("displayPath(\".\") = %q, want empty", got)
	}

	if got := displayPath("docs/a"); got != "docs/a" {
		t.Fatalf("displayPath(\"docs/a\") = %q, want %q", got, "docs/a")
	}
}

// ---------------------------------------------------------------------------
// openDir + readItems against a real temp filesystem
// ---------------------------------------------------------------------------

// fixture layout:
//
//	<base>/outside/secret.txt          (must never be reachable)
//	<base>/storage/src/A.txt
//	<base>/storage/src/b.txt
//	<base>/storage/src/alpha/          (dir)
//	<base>/storage/src/docs/c.txt
//	<base>/storage/src/inside-link  -> docs     (symlink staying inside)
//	<base>/storage/src/escape-link  -> outside  (symlink leaving the source)
//	<base>/storage/other/other.txt     (sibling source)
type fixture struct {
	handler *ListDirectoryHandler
	base    string
}

func newFixture(t *testing.T) *fixture {
	t.Helper()

	base := t.TempDir()
	storage := filepath.Join(base, "storage")
	src := filepath.Join(storage, "src")
	outside := filepath.Join(base, "outside")
	other := filepath.Join(storage, "other")

	for _, d := range []string{
		filepath.Join(src, "alpha"),
		filepath.Join(src, "docs"),
		outside,
		other,
	} {
		if err := os.MkdirAll(d, 0o750); err != nil {
			t.Fatal(err)
		}
	}

	files := map[string]string{
		filepath.Join(src, "A.txt"):          "a",
		filepath.Join(src, "b.txt"):          "bb",
		filepath.Join(src, "docs", "c.txt"):  "ccc",
		filepath.Join(outside, "secret.txt"): "top secret",
		filepath.Join(other, "other.txt"):    "other",
	}
	for p, content := range files {
		if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	// Symlinks may be unavailable (e.g. Windows without privileges).
	if err := os.Symlink(
		"docs",
		filepath.Join(src, "inside-link"),
	); err != nil {
		t.Logf("symlinks unavailable, symlink tests will skip: %v", err)
	}

	_ = os.Symlink(filepath.Join("..", "..", "outside"), filepath.Join(src, "escape-link"))

	return &fixture{
		handler: &ListDirectoryHandler{storageDir: storage},
		base:    base,
	}
}

// list runs the same chain Handle uses after the DB lookup.
func (f *fixture) list(t testing.TB, subpath, userPath string) ([]DirectoryEntry, error) {
	dir, err := f.handler.openDir(t.Context(), subpath, userPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = dir.Close()
		if err != nil {
			t.Logf("failed to close dir: %v", err)
		}
	}()

	return readItems(dir, userPath)
}

func entryNames(items []DirectoryEntry) []string {
	names := make([]string, 0, len(items))
	for _, it := range items {
		if strings.HasSuffix(it.Name, "-link") {
			continue
		}

		names = append(names, it.Name)
	}

	return names
}

func TestListRootSortedDirsFirst(t *testing.T) {
	t.Parallel()
	f := newFixture(t)

	items, err := f.list(t, "src", ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// directories first, then files, each group sorted case-insensitively
	want := []string{"alpha", "docs", "A.txt", "b.txt"}
	if got := entryNames(items); !slices.Equal(got, want) {
		t.Fatalf("entries = %v, want %v", got, want)
	}
}

func TestListSubdirectoryBuildsRelativePaths(t *testing.T) {
	t.Parallel()
	f := newFixture(t)

	items, err := f.list(t, "src", "docs")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("got %d items, want 1", len(items))
	}

	got := items[0]
	if got.Name != "c.txt" || got.Path != "docs/c.txt" || got.IsDir || got.Size != 3 {
		t.Fatalf("unexpected entry: %+v", got)
	}
}

func TestListRootPathsHaveNoDotPrefix(t *testing.T) {
	t.Parallel()
	f := newFixture(t)

	items, err := f.list(t, "src", ".")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, it := range items {
		if it.Path != it.Name {
			t.Fatalf("root entry path = %q, want %q", it.Path, it.Name)
		}
	}
}

func TestOpenDirRejectsEscapes(t *testing.T) {
	t.Parallel()
	f := newFixture(t)

	tests := []struct {
		name     string
		subpath  string
		userPath string
	}{
		{name: "parent of source", subpath: "src", userPath: ".."},
		{name: "reach sibling source", subpath: "src", userPath: "../other"},
		{name: "reach outside file", subpath: "src", userPath: "../../outside/secret.txt"},
		{name: "absolute user path", subpath: "src", userPath: "/etc/passwd"},
		{name: "absolute outside dir", subpath: "src", userPath: filepath.Join(f.base, "outside")},
		{name: "subpath escapes storage root", subpath: "../outside", userPath: "."},
		{name: "subpath absolute", subpath: filepath.Join(f.base, "outside"), userPath: "."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := f.handler.openDir(t.Context(), tt.subpath, tt.userPath)
			if err == nil {
				err = dir.Close()
				if err != nil {
					t.Logf("failed to close dir: %v", err)
				}

				t.Fatalf("openDir(%q, %q) succeeded, want error", tt.subpath, tt.userPath)
			}
		})
	}
}

func TestSymlinkEscapeIsBlocked(t *testing.T) {
	t.Parallel()
	f := newFixture(t)

	if _, err := os.Lstat(filepath.Join(f.handler.storageDir, "src", "escape-link")); err != nil {
		t.Skip("symlinks not available on this system")
	}

	// opening the link directly
	if dir, err := f.handler.openDir(t.Context(), "src", "escape-link"); err == nil {
		err = dir.Close()
		if err != nil {
			t.Logf("failed to close dir: %v", err)
		}

		t.Fatal("opening a symlink that leaves the source must fail")
	}

	// reaching a file through the link
	if dir, err := f.handler.openDir(t.Context(), "src", "escape-link/secret.txt"); err == nil {
		err = dir.Close()
		if err != nil {
			t.Logf("failed to close dir: %v", err)
		}

		t.Fatal("reading a file through an escaping symlink must fail")
	}
}

func TestSymlinkInsideSourceIsAllowed(t *testing.T) {
	t.Parallel()
	f := newFixture(t)

	if _, err := os.Lstat(filepath.Join(f.handler.storageDir, "src", "inside-link")); err != nil {
		t.Skip("symlinks not available on this system")
	}

	items, err := f.list(t, "src", "inside-link")
	if err != nil {
		t.Fatalf("symlink inside the source should work: %v", err)
	}

	if len(items) != 1 || items[0].Name != "c.txt" {
		t.Fatalf("unexpected items: %+v", items)
	}
}

func TestListMissingDirectory(t *testing.T) {
	t.Parallel()
	f := newFixture(t)

	if _, err := f.list(t, "src", "does-not-exist"); err == nil {
		t.Fatal("expected error for missing directory")
	}

	if _, err := f.list(t, "missing-source", "."); err == nil {
		t.Fatal("expected error for missing source directory")
	}
}

func TestListRegularFileFails(t *testing.T) {
	t.Parallel()
	f := newFixture(t)

	// openDir succeeds on a file, but readItems must fail (ReadDir on a file).
	if _, err := f.list(t, "src", "b.txt"); err == nil {
		t.Fatal("listing a regular file must fail")
	}
}

func TestListFollowsSourceReplacedBySymlinkOutside(t *testing.T) {
	t.Parallel()
	f := newFixture(t)

	// Simulates the TOCTOU case: a source directory that was valid at
	// registration is later replaced by a symlink pointing outside storage.
	storage := f.handler.storageDir
	if err := os.Remove(filepath.Join(storage, "other", "other.txt")); err != nil {
		t.Fatal(err)
	}

	if err := os.Remove(filepath.Join(storage, "other")); err != nil {
		t.Fatal(err)
	}

	if err := os.Symlink(
		filepath.Join(f.base, "outside"),
		filepath.Join(storage, "other"),
	); err != nil {
		t.Skip("symlinks not available on this system")
	}

	if dir, err := f.handler.openDir(t.Context(), "other", "."); err == nil {
		err = dir.Close()
		if err != nil {
			t.Logf("failed to close dir: %v", err)
		}

		t.Fatal("a source replaced by an outside symlink must not be openable")
	}
}
