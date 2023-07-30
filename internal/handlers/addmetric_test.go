package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zelas91/metric-collector/internal/advicerrors"
	"github.com/zelas91/metric-collector/internal/storages"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddMetric_MetricAdd(t *testing.T) {
	type want struct {
		code int
		body string
	}
	tests := []struct {
		name          string
		updateHandler MetricHandler
		want          want
		url           string
		method        string
	}{
		{
			name:          "Bad request #1",
			updateHandler: MetricHandler{mem: storages.NewMemStorage()},
			want:          want{code: 400, body: "bad request"},
			url:           "/update/unknown/testCounter/100",
			method:        http.MethodPost,
		},
		{
			name:          "Not found #2",
			updateHandler: MetricHandler{mem: storages.NewMemStorage()},
			want:          want{code: 404, body: "not found"},
			url:           "/",
			method:        http.MethodPost,
		},
		{
			name:          "Ok #3",
			updateHandler: MetricHandler{mem: storages.NewMemStorage()},
			want:          want{code: 200, body: ""},
			url:           "/update/counter/someMetric/527",
			method:        http.MethodPost,
		},
		{
			name:          "status method not allowed #4",
			updateHandler: MetricHandler{mem: storages.NewMemStorage()},
			want:          want{code: 405, body: "method not allowed"},
			url:           "/update/counter/someMetric/527",
			method:        http.MethodGet,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.url, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(advicerrors.Middleware(test.updateHandler.MetricAdd))
			h(w, request)
			read, err := io.ReadAll(w.Result().Body)
			require.NoError(t, err, "Body read error")
			err = w.Result().Body.Close()
			require.NoError(t, err, "Body close error")
			body := strings.TrimSpace(string(read))
			assert.Equal(t, test.want.code, w.Result().StatusCode, "status code not as expected")
			assert.Equal(t, test.want.body, body, "status code not as expected")
		})
	}

}
