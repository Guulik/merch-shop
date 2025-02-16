package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"merch/internal/app"
	"merch/internal/configure"
)

func main() {
	ctx := context.Background()
	cfg := configure.MustLoadConfig()

	a := app.New(ctx, cfg)

	go func() {
		a.MustRun(cfg.Port)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	if err := a.Stop(ctx); err != nil {
		fmt.Println(fmt.Errorf("failed to gracefully stop app: err=%s", err.Error()))
	}

	slog.Info("Gracefully stopped")
}
