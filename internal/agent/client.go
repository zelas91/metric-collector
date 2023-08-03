package agent

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/zelas91/metric-collector/internal/server/types"
	"sync"
	"time"
)

type ClientHttp struct {
	Client *resty.Client
}

func NewClientHttp() *ClientHttp {
	return &ClientHttp{Client: resty.New()}
}

func (c *ClientHttp) UpdateMetrics(s *Stats, baseUrl string) error {
	for name, value := range s.GetGauges() {
		resp, err := c.Client.R().SetPathParams(map[string]string{
			"type":  types.GaugeType,
			"name":  name,
			"value": fmt.Sprintf("%f", value.Value),
		}).SetHeader("Content-Type", "text/plain").Post(fmt.Sprintf("%s/{type}/{name}/{value}", baseUrl))
		if err != nil {
			return err
		}
		if resp.StatusCode() != 200 {
			return errors.New("answer result is not correct")
		}
	}

	for name, value := range s.GetCounters() {
		resp, err := c.Client.R().SetPathParams(map[string]string{
			"type":  types.GaugeType,
			"name":  name,
			"value": fmt.Sprintf("%d", value.Value),
		}).SetHeader("Content-Type", "text/plain").Post(fmt.Sprintf("%s/{type}/{name}/{value}", baseUrl))
		if err != nil {
			return err
		}
		if resp.StatusCode() != 200 {
			return errors.New("answer result is not correct")
		}
	}
	return nil
}

func Run(pollInterval, reportInterval time.Duration, baseURL string) {
	s := NewStats()
	c := NewClientHttp()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			<-time.After(pollInterval * time.Second)
			s.ReadStats()

		}
	}()
	go func() {
		for {
			<-time.After(reportInterval * time.Second)
			err := c.UpdateMetrics(s, baseURL)
			if err != nil {
				panic(err)
			}

		}
	}()
	wg.Wait()
}
