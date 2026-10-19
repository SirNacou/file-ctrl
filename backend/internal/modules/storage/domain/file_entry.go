package domain

import "time"

type FileEntry struct {
	Name    string
	Size    int64
	ModTime time.Time
	IsDir   bool
}
