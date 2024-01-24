package main

import (
	"context"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server"
	_ "google.golang.org/grpc/encoding/gzip"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := NewConfig()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	_ = cancel
	server.Run(ctx, cfg)
	<-ctx.Done()
	stop(ctx)
}
func stop(ctx context.Context) {
	server.Shutdown(ctx)
	logger.Shutdown()
	os.Exit(0)
}
