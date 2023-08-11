package agent

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/zelas91/metric-collector/internal/server/types"
	"os"
	"os/signal"
	"syscall"
	"time"
)

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
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		tickerReport := time.NewTicker(time.Duration(reportInterval) * time.Second)
		defer tickerReport.Stop()
		tickerPoll := time.NewTicker(time.Duration(pollInterval) * time.Second)
		defer tickerPoll.Stop()
		for {
			select {
			case <-tickerReport.C:
				err := c.UpdateMetrics(s, baseURL)
				if err != nil {
					logrus.Debug(err)
				}
			case <-tickerPoll.C:
				s.ReadStats()
			}
		}
	}()

	<-sigChan
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
