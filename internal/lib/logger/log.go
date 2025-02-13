package logger

import (
	"context"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envStage = "stage"
)

func InitLogger(env string) {
	var handler slog.Handler
	switch env {
	case envLocal:
		handler = slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envStage:
		handler = slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	handler = NewHandlerMiddleware(handler)
	slog.SetDefault(slog.New(handler))
}

type HandlerMiddleware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddleware {
	return &HandlerMiddleware{next: next}
}

func (h *HandlerMiddleware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *HandlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(key).(logCtx); ok {
		if c.UserID != 0 {
			rec.Add("userId", c.UserID)
		}
		if c.ToUser != 0 {
			rec.Add("to_user", c.ToUser)
		}
		if c.CoinBalance > 0 {
			rec.Add("coin_balance", c.CoinBalance)
		}
		if c.SendAmount > 0 {
			rec.Add("send_amount", c.SendAmount)
		}
		if c.Item != "" {
			rec.Add("item", c.Item)
		}
	}
	return h.next.Handle(ctx, rec)
}

func (h *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithAttrs(attrs)}
}

func (h *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return &HandlerMiddleware{next: h.next.WithGroup(name)}
}
