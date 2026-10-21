package domain

import "time"

type FileEntry struct {
	Name    string
	ModTime time.Time
	Size    int64
	IsDir   bool
}
