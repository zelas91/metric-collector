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
	"github.com/zelas91/metric-collector/internal/utils"
	"sync"
	"time"
)

var (
	log = logger.New()
)

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

type effectorUpdateMetrics func(s *Stats, baseURL, key string) error

func retryUpdateMetrics(effector effectorUpdateMetrics, exit <-chan time.Time) effectorUpdateMetrics {
	return func(s *Stats, baseURL, key string) error {
		retries := 3
		for r := 1; ; r++ {
			delay := time.Duration(r) * time.Second
			select {
			case <-time.After(delay):
			case <-exit:
				return errors.New("retry deadline exceeded")
			}
			if err := effector(s, baseURL, key); err == nil || r >= retries {
				return err
			}
		}
	}
}
func (c *ClientHTTP) UpdateMetrics(s *Stats, baseURL, key string) error {

	gauges := createGauges(s)
	counters := createCounters(s)
	metrics := append(gauges, counters...)

	headers := make(map[string]string)

	body, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("update metrics marshal err :%w", err)
	}

	gzipBody, err := gzipCompress(body)
	if err != nil {
		return fmt.Errorf("error compress body %w", err)
	}

	hash, err := utils.GenerateHash(gzipBody, key)

	if err != nil {
		if !errors.Is(err, utils.ErrInvalidKey) {
			return fmt.Errorf("update metrics genetate hash err:%w", err)
		}
		log.Errorf("Invalid hash key")
	}

	if hash != nil {
		headers["HashSHA256"] = *hash
	}
	headers["Content-Type"] = "application/json"
	headers["Content-Encoding"] = "gzip"

	resp, err := c.Client.R().
		SetHeaders(headers).
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

func Run(ctx context.Context, pollInterval, reportInterval int, baseURL, key string) {
	go func(ctx context.Context) {
		s := NewStats()
		c := NewClientHTTP()

		tickerReport := time.NewTicker(time.Duration(reportInterval) * time.Second)
		defer tickerReport.Stop()
		tickerPoll := time.NewTicker(time.Duration(pollInterval) * time.Second)
		defer tickerPoll.Stop()
		wg := new(sync.WaitGroup)
		for {
			select {
			case <-tickerReport.C:
				for i := 0; i < 2; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						if err := c.UpdateMetrics(s, baseURL, key); err != nil {
							log.Errorf("update metrics err: %v", err)
						}
					}()
				}
				wg.Wait()

			case <-tickerPoll.C:
				s.ReadStats()
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
}

func Run2(ctx context.Context, pollInterval, reportInterval int, baseURL, key string) {
	s := NewStats()
	c := NewClientHTTP()

	tickerReport := time.NewTicker(time.Duration(reportInterval) * time.Second)
	//defer tickerReport.Stop()
	tickerPoll := time.NewTicker(time.Duration(pollInterval) * time.Second)
	//defer tickerPoll.Stop()

	for w := 0; w < 5; w++ {
		w := w
		go func(ctx context.Context, tickerReport <-chan time.Time) {
			log.Info(w)
			for {
				select {
				case <-tickerReport:
					log.Info("ticket")
					if err := c.UpdateMetrics(s, baseURL, key); err != nil {
						log.Errorf("update metrics err: %v", err)
						r := retryUpdateMetrics(c.UpdateMetrics, tickerReport)
						if err = r(s, baseURL, key); err != nil {
							log.Errorf("retry err: %v", err)
						}
					}

				case <-ctx.Done():
					return

				}
			}
		}(ctx, tickerReport.C)
	}

	go func(ctx context.Context, tickerPoll <-chan time.Time) {
		for {
			select {
			case <-tickerPoll:
				s.ReadStats()
			case <-ctx.Done():
				return

			}
		}
	}(ctx, tickerPoll.C)
	//for {
	//	select {
	//	case <-tickerReport.C:
	//		if err := c.UpdateMetrics(s, baseURL, key); err != nil {
	//			log.Errorf("update metrics err: %v", err)
	//			r := retryUpdateMetrics(c.UpdateMetrics, tickerReport.C)
	//			if err = r(s, baseURL, key); err != nil {
	//				log.Errorf("retry err: %v", err)
	//			}
	//		}
	//
	//	case <-tickerPoll.C:
	//		s.ReadStats()
	//	case <-ctx.Done():
	//		return
	//	}
	//}
}
