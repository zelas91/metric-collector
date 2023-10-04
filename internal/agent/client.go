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
	"time"
)

var (
	log = logger.New()
)

type ClientHTTP struct {
	client *resty.Client
}

func NewClientHTTP() *ClientHTTP {
	client := resty.New()
	client.SetTimeout(2 * time.Second)
	return &ClientHTTP{client: client}
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

func createMemoryAndCpu(s *Stats) []repository.Metric {
	memory := s.GetMemoryAndCPU()
	metrics := make([]repository.Metric, 0, len(memory))
	for name, value := range memory {
		val := value
		metrics = append(metrics, repository.Metric{
			ID:    name,
			Value: &val,
			MType: types.GaugeType,
		})
	}
	return metrics
}

type effectorUpdateMetrics func(client *resty.Client, header map[string]string, body []byte, url string) error

func retryUpdateMetrics(effector effectorUpdateMetrics, exit <-chan time.Time) effectorUpdateMetrics {
	return func(client *resty.Client, header map[string]string, body []byte, url string) error {
		retries := 3
		for r := 1; ; r++ {
			delay := time.Duration(r) * time.Second
			select {
			case <-time.After(delay):
			case <-exit:
				return errors.New("retry deadline exceeded")
			}
			if err := effector(client, header, body, url); err == nil || r >= retries {
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

	resp, err := c.client.R().
		SetHeaders(headers).
		SetBody(gzipBody).
		EnableTrace().
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

func readStats(s *Stats, ch chan<- []repository.Metric) {
	s.ReadStats()
	ch <- append(createCounters(s), createGauges(s)...)
}

func Run(ctx context.Context, pollInterval, reportInterval int, baseURL, key string, rateLimit int) {
	s := NewStats()

	tickerReport := time.NewTicker(time.Duration(reportInterval) * time.Second)
	tickerPoll := time.NewTicker(time.Duration(pollInterval) * time.Second)
	tickerPollCpuAndMemory := time.NewTicker(1 * time.Second)

	reportChan := make(chan []repository.Metric, 64)
	updChan := make(chan []repository.Metric, 64)

	for w := 0; w < rateLimit; w++ {
		go updateMetrics(baseURL, key, updChan, tickerReport.C)
	}

	go func() {
		for {
			select {
			case <-tickerReport.C:
				copyChannel(reportChan, updChan)
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-tickerPoll.C:
				readStats(s, reportChan)
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-tickerPollCpuAndMemory.C:
				reportChan <- createMemoryAndCpu(s)
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		if <-ctx.Done(); true {
			close(updChan)
			close(reportChan)
		}
	}()
}
func copyChannel(src <-chan []repository.Metric, dst chan<- []repository.Metric) {
	for {
		select {
		case value := <-src:
			dst <- value
		default:
			return
		}
	}
}

func updateMetrics(baseURL, key string, report <-chan []repository.Metric, exit <-chan time.Time) {
	client := resty.New()
	client.SetTimeout(1 * time.Second)
	for m := range report {
		headers := make(map[string]string)

		body, err := json.Marshal(m)
		if err != nil {
			log.Errorf("update metrics marshal err :%v", err)
			continue
		}

		gzipBody, err := gzipCompress(body)
		if err != nil {
			log.Errorf("error compress body %v", err)
			continue
		}

		hash, err := utils.GenerateHash(gzipBody, key)

		if err != nil {
			if !errors.Is(err, utils.ErrInvalidKey) {
				log.Errorf("update metrics genetate hash err:%v", err)
				continue
			}
			log.Errorf("Invalid hash key")
		}

		if hash != nil {
			headers["HashSHA256"] = *hash
		}
		headers["Content-Type"] = "application/json"
		headers["Content-Encoding"] = "gzip"

		if err = requestPost(client, headers, gzipBody, baseURL); err != nil {
			r := retryUpdateMetrics(requestPost, exit)
			if err = r(client, headers, gzipBody, baseURL); err != nil {
				log.Errorf("retry err: %v", err)
			}
		}

	}
}

func requestPost(client *resty.Client, header map[string]string, body []byte, url string) error {
	resp, err := client.R().SetHeaders(header).
		SetBody(body).
		Post(url)
	if err != nil {
		return fmt.Errorf("error post request %w", err)

	}
	if resp.StatusCode() != 200 {
		return errors.New("answer result is not correct")

	}
	return nil
}

//func Run2(ctx context.Context, pollInterval, reportInterval int, baseURL, key string) {
//	s := NewStats()
//	c := NewClientHTTP()
//
//	tickerReport := time.NewTicker(time.Duration(reportInterval) * time.Second)
//	//defer tickerReport.Stop()
//	tickerPoll := time.NewTicker(time.Duration(pollInterval) * time.Second)
//	//defer tickerPoll.Stop()
//
//	for w := 0; w < 5; w++ {
//		w := w
//		go func(ctx context.Context, tickerReport <-chan time.Time) {
//			log.Info(w)
//			for {
//				select {
//				case <-tickerReport:
//					log.Info("ticket")
//					if err := c.UpdateMetrics(s, baseURL, key); err != nil {
//						log.Errorf("update metrics err: %v", err)
//						r := retryUpdateMetrics(c.UpdateMetrics, tickerReport)
//						if err = r(s, baseURL, key); err != nil {
//							log.Errorf("retry err: %v", err)
//						}
//					}
//
//				case <-ctx.Done():
//					return
//
//				}
//			}
//		}(ctx, tickerReport.C)
//	}
//
//	go func(ctx context.Context, tickerPoll <-chan time.Time) {
//		for {
//			select {
//			case <-tickerPoll:
//				s.ReadStats()
//			case <-ctx.Done():
//				return
//
//			}
//		}
//	}(ctx, tickerPoll.C)
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
//}
