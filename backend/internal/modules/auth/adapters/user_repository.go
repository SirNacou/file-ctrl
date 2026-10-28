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
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
)

type UserRepository struct {
	db bob.Executor
}

func NewUserRepository(db bob.Executor) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	return models.Users.Query().Count(ctx, r.db)
}

func (r *UserRepository) Create(
	ctx context.Context,
	user *domain.User,
) error {
	now := time.Now().UTC()

	_, err := models.Users.Insert(&models.UserSetter{
		ID:           omit.From(user.ID.String()),
		Username:     omit.From(user.Username),
		PasswordHash: omit.From(user.PasswordHash),
		Role:         omit.From(string(user.Role)),
		CreatedAt:    omit.From(now),
		UpdatedAt:    omit.From(now),
	}).One(ctx, r.db)

	return err
}

func (r *UserRepository) Find(
	ctx context.Context,
	ID uuid.UUID,
) (*domain.User, error) {
	row, err := models.FindUser(ctx, r.db, ID.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	id := uuid.MustParse(row.ID)
	return &domain.User{
		ID:           id,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		Role:         domain.Role(row.Role),
	}, nil
}


func (r *UserRepository) FindByUsername(
	ctx context.Context,
	username string,
) (*domain.User, error) {
	row, err := models.Users.Query(
		sm.Where(models.Users.Columns.Username.EQ(sqlite.Arg(username))),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	id := uuid.MustParse(row.ID)
	return &domain.User{
		ID:           id,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		Role:         domain.Role(row.Role),
	}, nil
}
