package app

import (
	"context"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/domain"
)

type SessionCreator interface {
	Create(
		ctx context.Context,
		session *domain.Session,
	) error
}

type SessionFinder interface {
	FindByTokenHash(
		ctx context.Context,
		tokenHash string,
	) (*domain.Session, error)
}
