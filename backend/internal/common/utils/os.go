package utils

import (
	"context"
	"fmt"
	"io"
	"log/slog"
)

func CloseLog(ctx context.Context, c io.Closer, what string) {
	if err := c.Close(); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("failed to close %s", what), slog.Any("error", err))
	}
}
