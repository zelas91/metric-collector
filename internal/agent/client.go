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

func (c *ClientHTTP) UpdateMetrics(s *Stats, baseURL string) error {
	for name, value := range s.GetGauges() {
		body, err := json.Marshal(repository.Metric{
			ID:    name,
			MType: types.GaugeType,
			Value: &value,
		})

		if err != nil {
			return fmt.Errorf("json marshal eroor = %v", err)
		}
		gzipBody, err := gzipCompress(body)
		if err != nil {
			return fmt.Errorf("error compress body %v", err)
		}
		resp, err := c.Client.R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
			SetBody(gzipBody).
			Post(baseURL)
		if err != nil {
			return fmt.Errorf("error post request %v", err)
		}
		if resp.StatusCode() != 200 {
			return errors.New("answer result is not correct")
		}
	}

	for name, value := range s.GetCounters() {
		body, err := json.Marshal(repository.Metric{
			ID:    name,
			MType: types.CounterType,
			Delta: &value,
		})
		if err != nil {
			return fmt.Errorf("json marshal eroor = %v", err)
		}

		gzipBody, err := gzipCompress(body)
		if err != nil {
			return fmt.Errorf("error compress body %v", err)
		}
		resp, err := c.Client.R().
			SetBody(gzipBody).
			SetHeader("Content-Type", "application/json").
			SetHeader("Content-Encoding", "gzip").
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
				err := c.UpdateMetrics(s, baseURL)
				if err != nil {
					log.Error(err)
				}
			case <-tickerPoll.C:
				s.ReadStats()
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
}
