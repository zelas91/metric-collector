package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/zelas91/metric-collector/internal/agent"
	"github.com/zelas91/metric-collector/internal/logger"
)

func main() {
	conf := NewConfig()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	_ = cancel
	agent.Run(ctx, conf.PollInterval, conf.ReportInterval, conf.BaseURL, conf.Key, *conf.RateLimit)
	log.Info("start agent")
	<-ctx.Done()
	stop()
}
func stop() {
	log.Info("stop agent")
	logger.Shutdown()
	os.Exit(0)
}
