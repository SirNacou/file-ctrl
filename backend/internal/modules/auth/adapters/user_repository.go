package adapters

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	"github.com/sirnacou/file-ctrl/backend/internal/common/database"
	"github.com/sirnacou/file-ctrl/backend/internal/common/database/gen/models"
	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/domain"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/sqlite"
	"github.com/stephenafamo/bob/dialect/sqlite/sm"
	"uuid"
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
		PasswordHash: omitnull.FromPtr(user.PasswordHash),
		Role:         omit.From(string(user.Role)),
		CreatedAt:    omit.From(database.FormatTime(now)),
		UpdatedAt:    omit.From(database.FormatTime(now)),
		Provider:     omit.From(string(user.Provider)),
		ProviderID:   omitnull.FromPtr(user.ProviderID),
	}).Exec(ctx, r.db)

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

	return userFromRow(row), nil
}

func (r *UserRepository) FindByProvider(
	ctx context.Context,
	provider string,
	providerID string,
) (*domain.User, error) {
	row, err := models.Users.Query(
		sm.Where(models.Users.Columns.Provider.EQ(sqlite.Arg(provider))),
		sm.Where(models.Users.Columns.ProviderID.EQ(sqlite.Arg(providerID))),
	).One(ctx, r.db)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return userFromRow(row), nil
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

	return userFromRow(row), nil
}

func userFromRow(row *models.User) *domain.User {
	return &domain.User{
		ID:           uuid.MustParse(row.ID),
		Username:     row.Username,
		PasswordHash: row.PasswordHash.Ptr(),
		Role:         domain.Role(row.Role),
		Provider:     domain.Provider(row.Provider),
		ProviderID:   row.ProviderID.Ptr(),
	}
}
