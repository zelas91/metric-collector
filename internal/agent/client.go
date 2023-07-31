package agent

import (
	"fmt"
	"github.com/zelas91/metric-collector/internal/server/utils/types"
	"net/http"
)

const (
	baseUrl = "http://localhost:8080/update"
)

type HttpClient struct {
	Client *http.Client
}

func (c *HttpClient) UpdateMetrics(s *Stats) {
	for k, v := range s.GetGauges() {
		url := fmt.Sprintf("%s/%s/%s/%f", baseUrl, types.GaugeType, k, v.Value)
		fmt.Println(url)
		resp, err := c.Client.Post(url, "text/plain", nil)
		fmt.Println(resp, err)
	}

	for k, v := range s.GetCounters() {
		url := fmt.Sprintf("%s/%s/%s/%d", baseUrl, types.CounterType, k, v.Value)
		fmt.Println(url)
		resp, err := c.Client.Post(url, "text/plain", nil)
		fmt.Println(resp, err)
	}

}
