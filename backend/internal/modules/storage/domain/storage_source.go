package domain

import (
	"strings"
	"time"
)

type StorageSource struct {
	Subpath   StorageSubpath
	CreatedAt time.Time
	Name      string
}

func NewStorageSource(subpath StorageSubpath, name string) (*StorageSource, error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return nil, ErrNameMustNotBeEmpty
	}

	return &StorageSource{
		Subpath:   subpath,
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
