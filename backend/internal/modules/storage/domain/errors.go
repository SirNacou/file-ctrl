package domain

import "fmt"

type ErrFailToReadDir struct {
	Path string
}

func NewErrFailToReadDir(path string) *ErrFailToReadDir {
	return &ErrFailToReadDir{
		Path: path,
	}
}

func (e *ErrFailToReadDir) Error() string {
	return fmt.Sprintf("Failed to read directory: %s", e.Path)
}
