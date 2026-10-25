package app

import (
	"context"
	"strings"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type (
	StagedFileCommitter interface {
		CommitStagedFile(
			stagingFilePath string,
			sourcePath domain.SourcePath,
			target domain.EntryPath,
		) error
	}
	UploadFileRequest struct {
		StoragePath    string
		ParentPath     string
		Filename       string
		StagedFilePath string
	}
	UploadFileResponse struct{}
	UploadFileHandler  struct {
		finder    StorageSourcesFinder
		committer StagedFileCommitter
	}
)

func NewUploadFileHandler(
	finder StorageSourcesFinder,
	committer StagedFileCommitter,
) *UploadFileHandler {
	return &UploadFileHandler{
		finder:    finder,
		committer: committer,
	}
}

func (h *UploadFileHandler) Handle(
	ctx context.Context,
	req *UploadFileRequest,
) (*UploadFileResponse, error) {
	stagedFilePath := strings.TrimSpace(req.StagedFilePath)
	if stagedFilePath == "" {
		return nil, domain.NewErrMustNotBeEmpty("staged file path")
	}

	source, err := h.finder.Find(ctx, req.StoragePath)
	if err != nil {
		return nil, domain.ErrStorageSourceNotFound
	}

	parent, err := domain.NewEntryPath(req.ParentPath)
	if err != nil {
		return nil, err
	}

	target, err := parent.ChildPath(req.Filename)
	if err != nil {
		return nil, err
	}

	if err := h.committer.CommitStagedFile(req.StagedFilePath, source.Path, target); err != nil {
	}

	return &UploadFileResponse{}, nil
}
