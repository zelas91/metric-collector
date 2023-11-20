package internal

import (
	"context"
	"encoding/json"
	"github.com/zelas91/metric-collector/internal/server/config"
	"github.com/zelas91/metric-collector/internal/server/controller"
	"github.com/zelas91/metric-collector/internal/server/repository"
	"github.com/zelas91/metric-collector/internal/server/service"
	"github.com/zelas91/metric-collector/internal/server/types"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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
	body, err := json.Marshal(metrics)
	if err != nil {
		log.Fatal(err)
	}
	w := httptest.NewRecorder()
	file := "/tmp/metrics-db.json"
	interval := 0
	h := controller.NewMetricHandler(service.NewMemService(context.Background(),
		repository.NewFileStorage(context.TODO(), &config.Config{FilePath: &file, StoreInterval: &interval}), &config.Config{})).InitRoutes(nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		request := httptest.NewRequest(http.MethodPost, "/updates", strings.NewReader(string(body)))
		request.Header = map[string][]string{"Content-Type": {"application/json"}}

		h.ServeHTTP(w, request)
	}
}

//func BenchmarkAddMetricJSON(b *testing.B) {
//	metrics := []repository.Metric{
//		{ID: "Test",
//			MType: types.GaugeType,
//			Value: new(float64)},
//		{ID: "Test2",
//			MType: types.GaugeType,
//			Value: new(float64)},
//		{ID: "Test3",
//			MType: types.GaugeType,
//			Value: new(float64)},
//		{ID: "Test4",
//			MType: types.GaugeType,
//			Value: new(float64)},
//		{ID: "Test5",
//			MType: types.CounterType,
//			Delta: new(int64)},
//	}
//	*metrics[0].Value = 20.75
//	*metrics[1].Value = 20.75
//	*metrics[2].Value = 20.75
//	*metrics[3].Value = 20.75
//	*metrics[4].Delta = 100
//	body, err := json.Marshal(metrics)
//	if err != nil {
//		log.Fatal(err)
//	}
//	mem := repository.NewMemStorage()
//	w := httptest.NewRecorder()
//	h := controller.NewMetricHandler(service.NewMemService(context.Background(),
//		mem, &config.Config{})).InitRoutes(nil)
//	b.ResetTimer()
//
//	for i := 0; i < b.N; i++ {
//		request := httptest.NewRequest(http.MethodPost, "/updates", strings.NewReader(string(body)))
//		request.Header = map[string][]string{"Content-Type": {"application/json"}}
//
//		h.ServeHTTP(w, request)
//	}
//}
