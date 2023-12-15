package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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
		conf.Key, *conf.RateLimit, readPublicKey(conf.CryptoCertPath))
	log.Info("start agent")
	<-ctx.Done()
	stop()
}
func readPublicKey(path string) *rsa.PublicKey {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Errorf(`read file="%s" , err=%v`, path, err)
		return nil
	}
	block, _ := pem.Decode(file)
	pub, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Errorf("parse certificate err=%v", err)
		return nil
	}

	if key, ok := pub.PublicKey.(*rsa.PublicKey); ok {
		return key
	}
	return nil

}
func stop() {
	log.Info("stop agent")
	logger.Shutdown()
	os.Exit(0)
}
