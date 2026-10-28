package adapters

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"uuid"

	"github.com/aarondl/opt/omit"
	"github.com/sirnacou/file-ctrl/backend/internal/common/database/gen/models"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/domain"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/dm"
)

type SessionRepository struct {
	db bob.Executor
}

func NewSessionRepository(db bob.Executor) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) Create(
	ctx context.Context,
	session *domain.Session,
) error {
	_, err := models.Sessions.Insert(&models.SessionSetter{
		ID:        omit.From(session.ID),
		UserID:    omit.From(session.UserID.String()),
		ExpiresAt: omit.From(session.ExpiresAt.UTC()),
		CreatedAt: omit.From(time.Now().UTC()),
	}).Exec(ctx, r.db)
	return err
}

func (r *SessionRepository) FindByTokenHash(
	ctx context.Context,
	tokenHash string,
) (*domain.Session, error) {
	sess, err := models.FindSession(ctx, r.db, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, err
	}

	return &domain.Session{
		ID:        sess.ID,
		UserID:    uuid.MustParse(sess.UserID),
		ExpiresAt: sess.ExpiresAt,
		CreatedAt: sess.CreatedAt,
	}, nil
}

func (r *SessionRepository) Delete(ctx context.Context, tokenHash string) error {
	_, err := models.Sessions.Delete(
		dm.Where(models.Sessions.Columns.ID.EQ(sqlite.Arg(tokenHash))),
	).Exec(ctx, r.db)
	return err
}
