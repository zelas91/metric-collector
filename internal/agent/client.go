package agent

import (
	"errors"
	"fmt"
	"github.com/zelas91/metric-collector/internal/server/utils/types"
	"net/http"
	"sync"
	"time"
)

type ClientHttp struct {
	Client *http.Client
}

func NewClientHttp() *ClientHttp {
	return &ClientHttp{Client: &http.Client{
		Timeout: 1 * time.Second,
	}}
}

func (c *ClientHttp) UpdateMetrics(s *Stats, baseUrl string) error {
	for k, v := range s.GetGauges() {
		url := fmt.Sprintf("%s/%s/%s/%f", baseUrl, types.GaugeType, k, v.Value)
		resp, err := c.Client.Post(url, "text/plain", nil)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			panic(errors.New("answer result is not correct"))
		}
	}

	for k, v := range s.GetCounters() {
		url := fmt.Sprintf("%s/%s/%s/%d", baseUrl, types.CounterType, k, v.Value)
		resp, err := c.Client.Post(url, "text/plain", nil)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			panic(errors.New("answer result is not correct"))
		}
	}
	return nil
}

func Run(pollInterval, reportInterval time.Duration, baseurl string) {
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
			err := c.UpdateMetrics(s, baseurl)
			if err != nil {
				panic(err)
			}

		}
	}()
	wg.Wait()
}
