package agent

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/zelas91/metric-collector/internal/server/types"
	"net/http"
	"time"
)

type ClientHTTP struct {
	Client *resty.Client
}

func NewClientHTTP() *ClientHTTP {
	return &ClientHTTP{Client: resty.NewWithClient(&http.Client{Timeout: 1 * time.Second})}
}

func (c *ClientHTTP) UpdateMetrics(s *Stats, baseURL string) error {
	for name, value := range s.GetGauges() {
		resp, err := c.Client.R().SetPathParams(map[string]string{
			"type":  types.GaugeType,
			"name":  name,
			"value": fmt.Sprintf("%f", value),
		}).SetHeader("Content-Type", "text/plain").Post(fmt.Sprintf("%s/{type}/{name}/{value}", baseURL))
		if err != nil {
			return fmt.Errorf("error post request %v", err)
		}
		if resp.StatusCode() != 200 {
			return errors.New("answer result is not correct")
		}
	}

	for name, value := range s.GetCounters() {
		resp, err := c.Client.R().SetPathParams(map[string]string{
			"type":  types.GaugeType,
			"name":  name,
			"value": fmt.Sprintf("%d", value),
		}).SetHeader("Content-Type", "text/plain").Post(fmt.Sprintf("%s/{type}/{name}/{value}", baseURL))
		if err != nil {
			return fmt.Errorf("error post request %v", err)
		}
		if resp.StatusCode() != 200 {
			return errors.New("answer result is not correct")
		}
	}
	return nil
}

func Run(pollInterval, reportInterval int, baseURL string) {
	s := NewStats()
	c := NewClientHTTP()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(time.Duration(reportInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err := c.UpdateMetrics(s, baseURL)
				if err != nil {
					logrus.Debug(err)
				}
			default:
				s.ReadStats()
				time.Sleep(time.Duration(pollInterval) * time.Second)
			}
		}
	}()

	<-ctx.Done()
}
