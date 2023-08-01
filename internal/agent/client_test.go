package agent

import (
	"github.com/stretchr/testify/assert"
	"github.com/zelas91/metric-collector/internal/server/advicerrors"
	"github.com/zelas91/metric-collector/internal/server/handlers"
	"github.com/zelas91/metric-collector/internal/server/storages"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetrics(t *testing.T) {
	t.Run("test update metric #1", func(t *testing.T) {

		handler := advicerrors.AdviceHandler(handlers.NewMetricHandler(storages.NewMemStorage()).MetricAdd)
		server := httptest.NewServer(handler)
		defer server.Close()

		client := NewClientHttp()
		s := NewStats()
		s.ReadStats()
		err := client.UpdateMetrics(s, server.URL+"/update")
		assert.NoError(t, err)
	})

}
