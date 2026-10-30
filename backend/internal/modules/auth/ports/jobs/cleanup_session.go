package jobs

import (
	"context"
	"log/slog"
	"time"

	"github.com/sirnacou/file-ctrl/backend/internal/modules/auth/app"
)

func CleanupSession(ctx context.Context, handler *app.RemoveExpiredSessionsHandler) {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case cur := <-ticker.C:
				ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
				if _, err := handler.Handle(
					ctx,
					&app.RemoveExpiredSessionsRequest{Cutoff: cur},
				); err != nil {
					slog.WarnContext(ctx, "failed to purge expired sessions", slog.Any("err", err))
				}
				cancel()
			}
		}
	}()
}
