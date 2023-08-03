package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zelas91/metric-collector/internal/server/storages"
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
		handler *Handler
		want    want
		url     string
		method  string
	}{
		{
			name:    "Bad request #1",
			handler: NewHandler(storages.NewMemStorage()),
			url:     "/update/unknown/testCounter/100",
			method:  http.MethodPost,
			want:    want{code: 400, body: ""},
		},
		{
			name:    "Ok #3",
			handler: NewHandler(storages.NewMemStorage()),
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
			defer func() {
				err := w.Result().Body.Close()
				require.NoError(t, err, "Body close error")
			}()

			statusCode := w.Result().StatusCode
			read, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err, "Body read error")

			result := strings.TrimSpace(string(read))
			assert.Equal(t, test.want.code, statusCode, "status code not as expected")
			assert.Equal(t, test.want.body, result, "status code not as expected")
		})
	}

}
func TestGetMetric(t *testing.T) {
	handler := &Handler{MemStore: &storages.MemStorage{Gauge: map[string]types.Gauge{
		"cpu":    {Value: 0.85},
		"memory": {Value: 0.6},
	}, Counter: map[string]types.Counter{
		"requests": {Value: 100},
		"errors":   {Value: 5},
	}}}
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

			defer func() {
				err := res.Body.Close()
				if err != nil {
					require.NoError(t, err, "Body close error")
				}
			}()

			statusCode := res.StatusCode

			assert.Equal(t, test.want, statusCode, "status code not as expected")
		})
	}
}
func TestIsType(t *testing.T) {
	tests := []struct {
		name    string
		want    bool
		strType string
	}{
		{
			name:    "test isType Gauge yes #1",
			want:    true,
			strType: types.GaugeType,
		},
		{
			name:    "test isType no #2",
			want:    false,
			strType: "Gauges",
		},
		{
			name:    "test isType Counter ok #3",
			want:    true,
			strType: types.CounterType,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, isType(test.strType))
		})
	}
}

func TestIsValue(t *testing.T) {
	tests := []struct {
		name  string
		want  bool
		value string
	}{
		{
			name:  "test float64 is value #1",
			want:  true,
			value: "12.5",
		}, {
			name:  "test int64 is value #2",
			want:  true,
			value: "12",
		},
		{
			name:  "test invalid is value #3",
			want:  false,
			value: "none",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, isValue(test.value))
		})
	}
}
