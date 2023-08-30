package main

import (
	"context"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	_ = cancel
	server.Run()
	<-ctx.Done()
	stop(ctx)
}
func stop(ctx context.Context) {
	logger.Shutdown()
	server.Shutdown(ctx)
	os.Exit(0)
}
