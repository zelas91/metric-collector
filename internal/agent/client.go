package agent

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/zelas91/metric-collector/internal/logger"
	"github.com/zelas91/metric-collector/internal/server/repository"
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

func createGauges(s *Stats) []repository.Metric {
	gauges := s.GetGauges()
	metrics := make([]repository.Metric, 0, len(gauges))
	for name, value := range gauges {
		val := value
		metrics = append(metrics, repository.Metric{
			ID:    name,
			Value: &val,
			MType: types.GaugeType,
		})
	}
	return metrics
}

func createCounters(s *Stats) []repository.Metric {
	counters := s.GetCounters()
	metrics := make([]repository.Metric, 0, len(counters))
	for name, value := range counters {
		val := value
		metrics = append(metrics, repository.Metric{
			ID:    name,
			Delta: &val,
			MType: types.CounterType,
		})
	}
	return metrics
}

type effectorUpdateMetrics func(s *Stats, baseURL string) error

func retryUpdateMetrics(effector effectorUpdateMetrics, ch <-chan time.Time) effectorUpdateMetrics {
	return func(s *Stats, baseURL string) error {
		var delay time.Duration
		retries := 3
		for r := 1; ; r++ {
			delay += 1 * time.Second
			select {
			case <-time.After(delay):
			case <-ch:
				return nil

			}

			if err := effector(s, baseURL); err == nil || r >= retries {
				return err
			}
		}
	}
}
func (c *ClientHTTP) UpdateMetrics(s *Stats, baseURL string) error {

	gauges := createGauges(s)
	counters := createCounters(s)

	metrics := append(gauges, counters...)

	body, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("update metrics marshal err :%w", err)
	}
	gzipBody, err := gzipCompress(body)
	if err != nil {
		return fmt.Errorf("error compress body %w", err)
	}
	resp, err := c.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(gzipBody).
		Post(baseURL)

	if err != nil {
		return fmt.Errorf("error post request %w", err)
	}
	if resp.StatusCode() != 200 {
		return errors.New("answer result is not correct")
	}
	return nil
}

func gzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, fmt.Errorf("gzip Compress error=%v", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to gzip : %v", err)

	}
	if err = w.Close(); err != nil {
		return nil, fmt.Errorf("gzip writer close error : %v", err)
	}
	return buf.Bytes(), err
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
				if err := c.UpdateMetrics(s, baseURL); err != nil {
					log.Errorf("update metrics err: %v", err)
					r := retryUpdateMetrics(c.UpdateMetrics, tickerReport.C)
					if err = r(s, baseURL); err != nil {
						log.Errorf("retry err: %v", err)
					}
				}

			case <-tickerPoll.C:
				s.ReadStats()
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
}
