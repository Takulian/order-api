package observability

import (
	"context"
	"log/slog"
	"runtime/debug"
)

func SafeGo(ctx context.Context, logger *slog.Logger, name string, fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.ErrorContext(ctx, "goroutine panic recovered",
					"goroutine:", name,
					"panic:", r,
					"stack", string(debug.Stack()),
				)
			}
		}()
		fn()
	}()
}
