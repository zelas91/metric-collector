package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/payload"
	"github.com/zelas91/metric-collector/internal/server/types"
	"time"
)

var log = logger.New()

type ClientHTTP struct {
	Client *resty.Client
}

func NewClientHTTP() *ClientHTTP {
	client := resty.New()
	client.SetTimeout(1 * time.Second)
	return &ClientHTTP{Client: client}
}

func (c *ClientHTTP) UpdateMetrics(s *Stats, baseURL string) error {
	for name, value := range s.GetGauges() {
		val := float64(value)
		body, err := json.Marshal(payload.Metrics{
			ID:    name,
			MType: types.GaugeType,
			Value: &val,
		})
		if err != nil {
			return fmt.Errorf("json marshal eroor = %v", err)
		}
		resp, err := c.Client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(body).
			Post(baseURL)
		if err != nil {
			return fmt.Errorf("error post request %v", err)
		}
		if resp.StatusCode() != 200 {
			return errors.New("answer result is not correct")
		}
	}

	for name, value := range s.GetCounters() {
		val := int64(value)
		body, err := json.Marshal(payload.Metrics{
			ID:    name,
			MType: types.CounterType,
			Delta: &val,
		})
		if err != nil {
			return fmt.Errorf("json marshal eroor = %v", err)
		}
		resp, err := c.Client.R().
			SetBody(body).
			SetHeader("Content-Type", "application/json").
			Post(baseURL)
		if err != nil {
			return fmt.Errorf("error post request %v", err)
		}
		if resp.StatusCode() != 200 {
			return errors.New("answer result is not correct")
		}
	}
	return nil
}

func Run(ctx context.Context, pollInterval, reportInterval int, baseURL string) {
	go func(ctx context.Context) {
		s := NewStats()
		c := NewClientHTTP()

		tickerReport := time.NewTicker(time.Duration(reportInterval) * time.Second)
		defer tickerReport.Stop()
		tickerPoll := time.NewTicker(time.Duration(pollInterval) * time.Second)
		defer tickerPoll.Stop()

		for {
			select {
			case <-tickerReport.C:
				err := c.UpdateMetrics(s, baseURL)
				if err != nil {
					log.Debug(err)
				}
			case <-tickerPoll.C:
				s.ReadStats()
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
}

//
//func Run(pollInterval, reportInterval int, baseURL string) {
//	s := NewStats()
//	c := NewClientHTTP()
//	sigChan := make(chan os.Signal, 1)
//	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
//	go func() {
//		poll := time.Now().Add(time.Duration(pollInterval) * time.Second)
//		report := time.Now().Add(time.Duration(reportInterval) * time.Second)
//
//		for {
//
//			if time.Now().After(poll) {
//				poll = time.Now().Add(time.Duration(reportInterval) * time.Second)
//				s.ReadStats()
//			}
//
//			if time.Now().After(report) {
//				report = time.Now().Add(time.Duration(reportInterval) * time.Second)
//				err := c.UpdateMetrics(s, baseURL)
//				if err != nil {
//					logrus.Debug(err)
//				}
//			}
//
//			time.Sleep(500 * time.Microsecond)
//
//		}
//
//	}()
//	<-sigChan
//}
