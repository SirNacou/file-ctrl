package domain

import (
	"strings"
	"time"
)

type StorageSource struct {
	Path      StoragePath
	CreatedAt time.Time
	Name      string
}

func NewStorageSource(path StoragePath, name string) (*StorageSource, error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return nil, ErrNameMustNotBeEmpty
	}

	return &StorageSource{
		Path:      path,
		Name:      name,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (s *StorageSource) Rename(name string) error {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return ErrNameMustNotBeEmpty
	}

	s.Name = name

	return nil
}
