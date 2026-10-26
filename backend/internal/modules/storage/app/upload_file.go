package app

import (
	"context"
	"strings"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type (
	StagedFileCommitter interface {
		CommitStagedFile(
			ctx context.Context,
			stagingFilePath string,
			sourcePath domain.SourcePath,
			target domain.EntryPath,
		) error
	}
	UploadFileRequest struct {
		StoragePath    string
		StorageSubpath     string
		DirectoryPath   string
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

	sourceSubpath, err := domain.NewEntryPath(req.StorageSubpath)
	if err != nil {
		return nil, err
	}

	targetDirectory, err := sourceSubpath.Join(req.DirectoryPath)
	if err != nil {
		return nil, err
	}

	target, err := targetDirectory.Join(req.Filename)
	if err != nil {
		return nil, err
	}

	if err := h.committer.CommitStagedFile(
		ctx,
		req.StagedFilePath,
		source.Path,
		target,
	); err != nil {
	}

	return &UploadFileResponse{}, nil
}
