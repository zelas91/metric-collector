package controller

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zelas91/metric-collector/internal/server/config"
	"github.com/zelas91/metric-collector/internal/server/controller/middleware"
	"github.com/zelas91/metric-collector/internal/server/repository"
	"github.com/zelas91/metric-collector/internal/server/service"
	mock_service "github.com/zelas91/metric-collector/internal/server/service/mocks"
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
			handler: NewMetricHandler(service.NewMemService(context.Background(), repository.NewMemStorage(), &config.Config{})),
			url:     "/update/unknown/testCounter/100",
			method:  http.MethodPost,
			want:    want{code: 400, body: ""},
		},
		{
			name:    "Ok #3",
			handler: NewMetricHandler(service.NewMemService(context.Background(), repository.NewMemStorage(), &config.Config{})),
			want:    want{code: 200, body: ""},
			url:     "/update/counter/someMetric/527",
			method:  http.MethodPost,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.url, nil)
			w := httptest.NewRecorder()
			h := test.handler.InitRoutes(nil)

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
	serv := service.NewMemService(context.Background(), repository.NewMemStorage(), &config.Config{})
	_, _ = serv.AddMetric(context.Background(), "cpu", types.GaugeType, "0.85")
	_, _ = serv.AddMetric(context.Background(), "memory", types.GaugeType, "0.6")
	_, _ = serv.AddMetric(context.Background(), "requests", types.CounterType, "100")
	_, _ = serv.AddMetric(context.Background(), "errors", types.CounterType, "5")
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

			h := handler.InitRoutes(nil)
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
		body    repository.Metric
		header  http.Header
		want    want
		url     string
		method  string
	}{
		{
			name:    "StatusUnsupportedMediaType #1",
			handler: NewMetricHandler(service.NewMemService(context.Background(), repository.NewMemStorage(), &config.Config{})),
			url:     "/update/",
			body: repository.Metric{
				ID: "Test",
			},
			header: map[string][]string{"Content-Type": {"text/html"}},
			method: http.MethodPost,
			want:   want{code: http.StatusUnsupportedMediaType, body: "{\"message\":\"incorrect media type \"}"},
		},
		{
			name:    "Ok  gauge #2",
			handler: NewMetricHandler(service.NewMemService(context.Background(), repository.NewMemStorage(), &config.Config{})),
			want:    want{code: http.StatusOK, body: "{\"id\":\"Test\",\"type\":\"gauge\",\"value\":20.12}"},
			url:     "/update/",
			body: repository.Metric{
				ID:    "Test",
				MType: types.GaugeType,
				Value: &gaugeVal,
			},
			header: map[string][]string{"Content-Type": {"application/json"}},
			method: http.MethodPost,
		},
		{
			name:    "Bad request  #3",
			handler: NewMetricHandler(service.NewMemService(context.Background(), repository.NewMemStorage(), &config.Config{})),
			want:    want{code: http.StatusBadRequest, body: "{\"message\":\"counter delta not found\"}"},
			url:     "/update/",
			body: repository.Metric{
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
			h := test.handler.InitRoutes(nil)

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
	serv := service.NewMemService(context.Background(), repository.NewMemStorage(), &config.Config{})
	_, _ = serv.AddMetric(context.Background(), "cpu", types.GaugeType, "0.85")
	_, _ = serv.AddMetric(context.Background(), "memory", types.GaugeType, "0.6")
	_, _ = serv.AddMetric(context.Background(), "requests", types.CounterType, "100")
	_, _ = serv.AddMetric(context.Background(), "errors", types.CounterType, "5")
	handler := NewMetricHandler(serv)
	type result struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name   string
		want   result
		url    string
		body   repository.Metric
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
			body: repository.Metric{
				ID:    "cpu",
				MType: types.GaugeType,
			},
			header: map[string][]string{"Content-Type": {"application/json"}},
			method: http.MethodPost,
		},
		{
			name: "Get not found",
			want: result{
				statusCode: http.StatusNotFound,
				body:       "{\"message\":\"not found metrics\"}",
			},
			body: repository.Metric{
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
			h := handler.InitRoutes(nil)
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

func BenchmarkAddMetricJSONFile(b *testing.B) {
	metrics := []repository.Metric{
		{ID: "Test",
			MType: types.GaugeType,
			Value: new(float64)},
		{ID: "Test2",
			MType: types.GaugeType,
			Value: new(float64)},
		{ID: "Test3",
			MType: types.GaugeType,
			Value: new(float64)},
		{ID: "Test4",
			MType: types.GaugeType,
			Value: new(float64)},
		{ID: "Test5",
			MType: types.CounterType,
			Delta: new(int64)},
	}
	*metrics[0].Value = 20.75
	*metrics[1].Value = 20.75
	*metrics[2].Value = 20.75
	*metrics[3].Value = 20.75
	*metrics[4].Delta = 100
	bodyJSON, err := json.Marshal(metrics)
	if err != nil {
		log.Fatal(err)
	}
	var body bytes.Buffer
	gz := gzip.NewWriter(&body)
	if _, err := gz.Write(bodyJSON); err != nil {
		log.Fatal(err)
	}
	if err = gz.Close(); err != nil {
		log.Fatal(err)
	}
	w := httptest.NewRecorder()
	file := "/tmp/metrics-db.json"
	interval := 0
	h := NewMetricHandler(service.NewMemService(context.Background(),
		repository.NewFileStorage(context.TODO(), &config.Config{FilePath: &file,
			StoreInterval: &interval}), &config.Config{})).InitRoutes(nil)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		request := httptest.NewRequest(http.MethodPost, "/updates", strings.NewReader(body.String()))
		request.Header = map[string][]string{"Content-Type": {"application/json"}, "Content-Encoding": {"gzip"}}
		b.StartTimer()
		h.ServeHTTP(w, request)
	}
}

func BenchmarkAddMetricJSON(b *testing.B) {
	metrics := []repository.Metric{
		{ID: "Test",
			MType: types.GaugeType,
			Value: new(float64)},
		{ID: "Test2",
			MType: types.GaugeType,
			Value: new(float64)},
		{ID: "Test3",
			MType: types.GaugeType,
			Value: new(float64)},
		{ID: "Test4",
			MType: types.GaugeType,
			Value: new(float64)},
		{ID: "Test5",
			MType: types.CounterType,
			Delta: new(int64)},
	}
	*metrics[0].Value = 20.75
	*metrics[1].Value = 20.75
	*metrics[2].Value = 20.75
	*metrics[3].Value = 20.75
	*metrics[4].Delta = 100
	bodyJSON, err := json.Marshal(metrics)
	if err != nil {
		log.Fatal(err)
	}
	var body bytes.Buffer
	gz := gzip.NewWriter(&body)
	if _, err := gz.Write(bodyJSON); err != nil {
		log.Fatal(err)
	}
	if err = gz.Close(); err != nil {
		log.Fatal(err)
	}
	mem := repository.NewMemStorage()
	w := httptest.NewRecorder()
	h := NewMetricHandler(service.NewMemService(context.Background(),
		mem, &config.Config{})).InitRoutes(nil)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		request := httptest.NewRequest(http.MethodPost, "/updates", strings.NewReader(body.String()))
		request.Header = map[string][]string{"Content-Type": {"application/json"}, "Content-Encoding": {"gzip"}}
		b.StartTimer()
		h.ServeHTTP(w, request)
	}

}

func TestGetMetrics(t *testing.T) {
	val := 29.23
	tests := []struct {
		name         string
		mockBehavior func(s *mock_service.MockService)
		method       string
		url          string
		statusCode   int
	}{{
		name:       "OK",
		statusCode: http.StatusOK,
		mockBehavior: func(s *mock_service.MockService) {
			s.EXPECT().GetMetrics(gomock.Any()).Return([]repository.Metric{
				{ID: "CPUZ",
					MType: types.GaugeType,
					Value: &val},
			})
		},
		url:    "/",
		method: http.MethodGet,
	}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			serv := mock_service.NewMockService(ctrl)
			test.mockBehavior(serv)
			handler := NewMetricHandler(serv)

			request := httptest.NewRequest(test.method, test.url, nil)
			//request.Header = test.header
			w := httptest.NewRecorder()
			h := handler.InitRoutes(nil)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			statusCode := res.StatusCode
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err, "Body read error")
			assert.Equal(t, test.statusCode, statusCode)

		})
	}
}
func createGinContextDecompress(b *testing.B, body string) *gin.Context {
	w := httptest.NewRecorder()
	b.StopTimer()
	c, _ := gin.CreateTestContext(w)
	req, _ := http.NewRequest("post", "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	c.Request = req
	b.StartTimer()
	return c
}
func BenchmarkGzipDecompressMiddleware(b *testing.B) {
	metrics := []repository.Metric{
		{ID: "Test",
			MType: types.GaugeType,
			Value: new(float64)},
		{ID: "Test2",
			MType: types.GaugeType,
			Value: new(float64)},
		{ID: "Test3",
			MType: types.GaugeType,
			Value: new(float64)},
		{ID: "Test4",
			MType: types.GaugeType,
			Value: new(float64)},
		{ID: "Test5",
			MType: types.CounterType,
			Delta: new(int64)},
	}
	*metrics[0].Value = 20.75
	*metrics[1].Value = 20.75
	*metrics[2].Value = 20.75
	*metrics[3].Value = 20.75
	*metrics[4].Delta = 100
	bodyJSON, err := json.Marshal(metrics)
	if err != nil {
		log.Fatal(err)
	}
	var body bytes.Buffer
	gz := gzip.NewWriter(&body)
	gz.Write(bodyJSON)
	gz.Close()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		middleware.GzipDecompress(createGinContextDecompress(b, body.String()))
	}
}
func createGinContextCompress(b *testing.B) *gin.Context {
	b.StopTimer()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")

	c.Request = req
	b.StartTimer()
	return c
}
func BenchmarkGzipCompressMiddleware(b *testing.B) {

	for i := 0; i < b.N; i++ {
		middleware.GzipCompress(createGinContextCompress(b))
	}
}
