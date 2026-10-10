package adapters

import (
	"context"
	"fmt"

	"github.com/aarondl/opt/omit"
	"github.com/sirnacou/file-ctrl/backend/internal/common/database/gen/models"
	"github.com/stephenafamo/bob"
)

type StorageSourcesRepository struct {
	db bob.Executor
}

func NewStorageSourcesRepository(db bob.Executor) *StorageSourcesRepository {
	return &StorageSourcesRepository{db: db}
}

func (r *StorageSourcesRepository) WithTx(tx *bob.Tx) *StorageSourcesRepository {
	r.db = tx
	return r
}

func (r *StorageSourcesRepository) Insert(ctx context.Context, model *models.StorageSource) error {
	_, err := models.StorageSources.Insert(&models.StorageSourceSetter{
		ID:        omit.From(model.ID),
		Name:      omit.From(model.Name),
		Subpath:   omit.From(model.Subpath),
		CreatedAt: omit.From(model.CreatedAt),
	}).Exec(ctx, r.db)

	if err != nil {
		return fmt.Errorf("failed to insert new Storage Source: %v", err)
	}

	return nil
}

func (r *StorageSourcesRepository) Update(ctx context.Context, model *models.StorageSource) error {
	_, err := models.StorageSources.
		Update(models.UpdateWhere.StorageSources.ID.EQ(model.ID),
			models.StorageSourceSetter{
				Name:    omit.From(model.Name),
				Subpath: omit.From(model.Subpath),
			}.UpdateMod()).
		Exec(ctx, r.db)

	if err != nil {
		return fmt.Errorf("failed to update Storage Source \"%s\": %v", model.ID, err)
	}

	return nil
}
