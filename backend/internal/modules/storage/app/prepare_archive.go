package app

import (
	"context"
	"io"
	"io/fs"
	"log/slog"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/storage/domain"
)

type PrepareArchiveRequest struct {
	StoragePath string
	FilePaths   []string
}
type PrepareArchiveResponse struct {
	Archive *domain.Archive
}

type PrepareArchiveHandler struct {
	fsOpener FSOpener
	finder   StorageSourcesFinder
}

type FSOpener interface {
	OpenFS(sourcePath domain.SourcePath) (fs.FS, io.Closer, error)
}

func NewPrepareArchiveHandler(
	fsOpener FSOpener,
	finder StorageSourcesFinder,
) *PrepareArchiveHandler {
	return &PrepareArchiveHandler{
		fsOpener: fsOpener,
		finder:   finder,
	}
}

func (h *PrepareArchiveHandler) Handle(
	ctx context.Context,
	req *PrepareArchiveRequest,
) (*PrepareArchiveResponse, error) {
	slog.DebugContext(ctx, "archive", slog.Any("paths", req.FilePaths))

	source, err := h.finder.Find(ctx, req.StoragePath)
	if err != nil {
		return nil, err
	}

	entryPaths := make([]domain.EntryPath, 0, len(req.FilePaths))
	for _, path := range req.FilePaths {
		entryPath, err := domain.NewEntryPath(path)
		if err != nil {
			slog.WarnContext(
				ctx,
				"failed to create entry path",
				slog.String("path", path),
				slog.Any("error", err),
			)
			return nil, err
		}
		entryPaths = append(entryPaths, entryPath)
	}

	fsys, closer, err := h.fsOpener.OpenFS(source.Path)
	if err != nil {
		return nil, err
	}

	archiveName, err := domain.NewArchiveName(source.Path, entryPaths)
	if err != nil {
		return nil, err
	}
	slog.DebugContext(ctx, "archive name", slog.String("name", archiveName.String()))

	a, err := domain.NewArchive(fsys, closer, archiveName, entryPaths)
	if err != nil {
		_ = closer.Close()
		return nil, err
	}

	return &PrepareArchiveResponse{
		Archive: a,
	}, nil
}
