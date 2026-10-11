package app

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListFilesHandler(t *testing.T) {
	root := t.TempDir()

	_ = os.Mkdir(filepath.Join(root, "docs"), 0755)
	_ = os.WriteFile(filepath.Join(root, "file.txt"), []byte("content"), 0644)
	_ = os.Mkdir(filepath.Join(root, "empty_dir"), 0755)

	handler := NewListFilesHandler(root)

	tests := []struct {
		name          string
		inputPath     string
		expectedCount int
		expectErr     bool
	}{
		{
			name:          "root directory listing",
			inputPath:     "/",
			expectedCount: 3, // docs, file.txt, empty_dir
			expectErr:     false,
		},
		{
			name:          "empty subdirectory",
			inputPath:     "/empty_dir",
			expectedCount: 0,
			expectErr:     false,
		},
		{
			name:          "non-existent directory",
			inputPath:     "/does-not-exist",
			expectedCount: 0,
			expectErr:     true,
		},
		{
			name:          "path traversal attempt",
			inputPath:     "/../../etc",
			expectedCount: 0,
			expectErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := handler.Handle(t.Context(), &ListFilesRequest{
				Path: tt.inputPath,
			})

			if tt.expectErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.expectedCount != len(resp.Items) {
				t.Errorf("expected %d items, got %d", tt.expectedCount, len(resp.Items))
			}
		})
	}
}
