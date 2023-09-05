package main

import (
	"context"
	"github.com/zelas91/metric-collector/internal/agent"
	"github.com/zelas91/metric-collector/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := NewConfig()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	_ = cancel
	agent.Run(ctx, conf.PollInterval, conf.ReportInterval, conf.BaseURL)
	<-ctx.Done()
	stop()
}
func stop() {
	logger.Shutdown()
	os.Exit(0)
}
