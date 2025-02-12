package main

import (
	"context"
	"fmt"
	"merch/internal/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	a := app.New()

	ctx := context.Background()

	go func() {
		a.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	if err := a.Stop(ctx); err != nil {
		fmt.Println(fmt.Errorf("failed to gracefully stop app: err=%s", err.Error()))
	}

	fmt.Println("Gracefully stopped")
}
