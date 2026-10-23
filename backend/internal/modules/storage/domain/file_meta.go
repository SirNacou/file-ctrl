package domain

import "time"

type FileMeta struct {
	Name    string
	ModTime time.Time
	Size    int64
	IsDir   bool
}
