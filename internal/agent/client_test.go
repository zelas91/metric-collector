package agent

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateMetrics(t *testing.T) {
	t.Run("test update metric #1", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			fmt.Println(request.URL)
			fmt.Println(request.URL.String())
			fmt.Println(request.RequestURI)
			parts := strings.Split(strings.TrimPrefix(request.RequestURI, "/"), "/")
			assert.Equal(t, 3, len(parts))
		}))
		defer server.Close()

		client := NewClientHttp()
		s := NewStats()
		s.ReadStats()
		err := client.UpdateMetrics(s, server.URL)

		assert.NoError(t, err)
	})

}
