package controller

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zelas91/metric-collector/internal/server/config"
	"github.com/zelas91/metric-collector/internal/server/payload"
	"github.com/zelas91/metric-collector/internal/server/repository"
	"github.com/zelas91/metric-collector/internal/server/service"
	"github.com/zelas91/metric-collector/internal/server/types"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddMetric(t *testing.T) {
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name    string
		handler *MetricHandler
		want    want
		url     string
		method  string
	}{
		{
			name:    "Bad request #1",
			handler: NewMetricHandler(service.NewMetricsService(repository.NewMemStorage(), &config.Config{}, context.Background())),
			url:     "/update/unknown/testCounter/100",
			method:  http.MethodPost,
			want:    want{code: 400, body: ""},
		},
		{
			name:    "Ok #3",
			handler: NewMetricHandler(service.NewMetricsService(repository.NewMemStorage(), &config.Config{}, context.Background())),
			want:    want{code: 200, body: ""},
			url:     "/update/counter/someMetric/527",
			method:  http.MethodPost,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.url, nil)
			w := httptest.NewRecorder()
			h := test.handler.InitRoutes()

			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			statusCode := res.StatusCode
			read, err := io.ReadAll(res.Body)
			require.NoError(t, err, "Body read error")

			result := strings.TrimSpace(string(read))
			assert.Equal(t, test.want.code, statusCode, "status code not as expected")
			assert.Equal(t, test.want.body, result, "status code not as expected")
		})
	}

}
func TestGetMetric(t *testing.T) {
	serv := service.NewMetricsService(repository.NewMemStorage(), &config.Config{}, context.Background())
	_ = serv.AddMetric("cpu", types.GaugeType, "0.85")
	_ = serv.AddMetric("memory", types.GaugeType, "0.6")
	_ = serv.AddMetric("requests", types.CounterType, "100")
	_ = serv.AddMetric("errors", types.CounterType, "5")
	handler := NewMetricHandler(serv)
	tests := []struct {
		name string
		want int
		url  string
	}{
		{
			name: "Get OK",
			want: http.StatusOK,
			url:  "/value/gauge/memory",
		},
		{
			name: "Get not found",
			want: http.StatusNotFound,
			url:  "/value/gauge/zpu",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.url, nil)
			w := httptest.NewRecorder()

			h := handler.InitRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			statusCode := res.StatusCode

			assert.Equal(t, test.want, statusCode, "status code not as expected")
		})
	}
}

func TestAddMetricJSON(t *testing.T) {
	type want struct {
		code int
		body string
	}
	gaugeVal := 20.12
	tests := []struct {
		name    string
		handler *MetricHandler
		body    payload.Metrics
		header  http.Header
		want    want
		url     string
		method  string
	}{
		{
			name:    "StatusUnsupportedMediaType #1",
			handler: NewMetricHandler(service.NewMetricsService(repository.NewMemStorage(), &config.Config{}, context.Background())),
			url:     "/update/",
			body: payload.Metrics{
				ID: "Test",
			},
			header: map[string][]string{"Content-Type": {"text/html"}},
			method: http.MethodPost,
			want:   want{code: http.StatusUnsupportedMediaType, body: "{\"message\":\"incorrect media type \"}"},
		},
		{
			name:    "Ok  gauge #2",
			handler: NewMetricHandler(service.NewMetricsService(repository.NewMemStorage(), &config.Config{}, context.Background())),
			want:    want{code: http.StatusOK, body: "{\"id\":\"Test\",\"type\":\"gauge\",\"value\":20.12}"},
			url:     "/update/",
			body: payload.Metrics{
				ID:    "Test",
				MType: types.GaugeType,
				Value: &gaugeVal,
			},
			header: map[string][]string{"Content-Type": {"application/json"}},
			method: http.MethodPost,
		},
		{
			name:    "Bad request  #3",
			handler: NewMetricHandler(service.NewMetricsService(repository.NewMemStorage(), &config.Config{}, context.Background())),
			want:    want{code: http.StatusBadRequest, body: "{\"message\":\"counter delta not found\"}"},
			url:     "/update/",
			body: payload.Metrics{
				ID:    "Test",
				MType: types.CounterType,
				Value: &gaugeVal,
			},
			header: map[string][]string{"Content-Type": {"application/json"}},
			method: http.MethodPost,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.body)
			require.NoError(t, err, "Body write error")

			request := httptest.NewRequest(test.method, test.url, strings.NewReader(string(body)))
			request.Header = test.header
			w := httptest.NewRecorder()
			h := test.handler.InitRoutes()

			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			statusCode := res.StatusCode
			read, err := io.ReadAll(res.Body)
			require.NoError(t, err, "Body read error")

			result := strings.TrimSpace(string(read))
			assert.Equal(t, test.want.code, statusCode, "status code not as expected")
			assert.Equal(t, test.want.body, result, "status code not as expected")
		})
	}
}
func TestGetMetricJSON(t *testing.T) {
	serv := service.NewMetricsService(repository.NewMemStorage(), &config.Config{}, context.Background())
	_ = serv.AddMetric("cpu", types.GaugeType, "0.85")
	_ = serv.AddMetric("memory", types.GaugeType, "0.6")
	_ = serv.AddMetric("requests", types.CounterType, "100")
	_ = serv.AddMetric("errors", types.CounterType, "5")
	handler := NewMetricHandler(serv)
	type result struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name   string
		want   result
		url    string
		body   payload.Metrics
		header http.Header
		method string
	}{
		{
			name: "Get OK",
			want: result{
				statusCode: http.StatusOK,
				body:       "{\"id\":\"cpu\",\"type\":\"gauge\",\"value\":0.85}",
			},
			url: "/value/",
			body: payload.Metrics{
				ID:    "cpu",
				MType: types.GaugeType,
			},
			header: map[string][]string{"Content-Type": {"application/json"}},
			method: http.MethodPost,
		},
		{
			name: "Get not found",
			want: result{
				statusCode: http.StatusOK,
				body:       "{\"id\":\"cpuz\",\"type\":\"gauge\",\"value\":0}",
			},
			body: payload.Metrics{
				ID:    "cpuz",
				MType: types.GaugeType,
			},
			url:    "/value/",
			method: http.MethodPost,
			header: map[string][]string{"Content-Type": {"application/json"}},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.body)
			require.NoError(t, err, "Body write error")

			request := httptest.NewRequest(test.method, test.url, strings.NewReader(string(body)))
			request.Header = test.header
			w := httptest.NewRecorder()
			h := handler.InitRoutes()
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			statusCode := res.StatusCode
			read, err := io.ReadAll(res.Body)
			require.NoError(t, err, "Body read error")

			result := strings.TrimSpace(string(read))
			assert.Equal(t, test.want.statusCode, statusCode, "status code not as expected")
			assert.Equal(t, test.want.body, result, "status code not as expected")
		})
	}
}
