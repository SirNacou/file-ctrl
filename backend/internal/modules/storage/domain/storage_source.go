package domain

import (
	"strings"
	"time"
	"uuid"
)

type StorageSource struct {
	CreatedAt time.Time
	Subpath   SourceSubpath
	Name      string
	ID        uuid.UUID
}

func NewStorageSource(name string, subPath SourceSubpath) (*StorageSource, error) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return nil, ErrNameMustNotBeEmpty
	}

	return &StorageSource{
		ID:        uuid.NewV7(),
		Name:      name,
		Subpath:   subPath,
		CreatedAt: time.Now().UTC(),
	}, nil
}
