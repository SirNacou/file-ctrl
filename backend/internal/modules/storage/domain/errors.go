package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNameMustNotBeEmpty      = errors.New("name must not be empty")
	ErrStorageSourceNotFound   = errors.New("storage source not found")
	ErrInvalidSourceDirectory  = errors.New("source directory is invalid")
	ErrSourceDirectoryNotFound = errors.New("source directory is not exists")
	ErrInvalidPath             = errors.New("invalid path")
	ErrAlreadyExists           = errors.New("already exists")
	ErrFileNotFound            = errors.New("file not found")
	ErrNotFound                = errors.New("not found")
	ErrNotAFile                = errors.New("target path is not a file")
	ErrEmptySelection          = errors.New("no paths selected")
)

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
