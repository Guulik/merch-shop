package main

import (
	"context"
	"fmt"
	"log/slog"
	"merch/internal/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()

	a := app.New(ctx)

	go func() {
		a.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	if err := a.Stop(ctx); err != nil {
		fmt.Println(fmt.Errorf("failed to gracefully stop app: err=%s", err.Error()))
	}

	slog.Info("Gracefully stopped")
}
