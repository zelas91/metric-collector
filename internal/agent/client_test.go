package agent

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/zelas91/metric-collector/internal/server/config"
	"github.com/zelas91/metric-collector/internal/server/controller"
	"github.com/zelas91/metric-collector/internal/server/repository"
	"github.com/zelas91/metric-collector/internal/server/service"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetrics(t *testing.T) {
	t.Run("test update metric #1", func(t *testing.T) {

		handler := controller.NewMetricHandler(service.NewMemService(context.Background(), repository.NewMemStorage(), &config.Config{}))
		server := httptest.NewServer(handler.InitRoutes(nil))
		defer server.Close()

		client := NewClientHTTP()
		s := NewStats()
		s.ReadStats()
		err := client.UpdateMetrics(s, server.URL+"/updates", "")
		assert.NoError(t, err)
	})

}
