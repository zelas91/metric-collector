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
	agent.Run(ctx, conf.PollInterval, conf.ReportInterval, conf.BaseURL,
		conf.Key, *conf.RateLimit, conf.CryptoCertPath, *conf.Mode)
	log.Info("start agent")
	<-ctx.Done()
	stop()
}

func stop() {
	log.Info("stop agent")
	logger.Shutdown()
	os.Exit(0)
}
