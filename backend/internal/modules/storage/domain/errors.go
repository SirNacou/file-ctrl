package domain

import (
	"errors"
	"fmt"
)

var (
	ErrStorageSourceNotFound   = errors.New("storage source not found")
	ErrInvalidSourceDirectory  = errors.New("source directory is invalid")
	ErrSourceDirectoryNotFound = errors.New("source directory is not exists")
	ErrAlreadyExists           = errors.New("already exists")
	ErrFileNotFound            = errors.New("file not found")
	ErrNotFound                = errors.New("not found")
	ErrNotAFile                = errors.New("target path is not a file")
	ErrEmptySelection          = errors.New("no paths selected")
)

type ErrInvalidPath struct {
	path string
}

func NewErrInvalidPath(path string) *ErrInvalidPath {
	return &ErrInvalidPath{
		path: path,
	}
}

func (e *ErrInvalidPath) Error() string {
	return fmt.Sprintf("invalid path %s", e.path)
}

type ErrMustNotBeEmpty struct {
	property string
}

func NewErrMustNotBeEmpty(property string) *ErrMustNotBeEmpty {
	return &ErrMustNotBeEmpty{property: property}
}

func (e *ErrMustNotBeEmpty) Error() string {
	return fmt.Sprintf("%s must not be empty", e.property)
}

type ErrFailToReadDir struct {
	path string
}

func NewErrFailToReadDir(path string) *ErrFailToReadDir {
	return &ErrFailToReadDir{
		path: path,
	}
}

func (e *ErrFailToReadDir) Error() string {
	return fmt.Sprintf("failed to read directory: %s", e.path)
}
