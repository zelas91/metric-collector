package server

import (
	"fmt"
	"github.com/zelas91/metric-collector/internal/server/advicerrors"
	"github.com/zelas91/metric-collector/internal/server/handlers"
	"github.com/zelas91/metric-collector/internal/server/storages"
	"net/http"
)

func Run(port string) error {
	hand := handlers.NewMetricHandler(storages.NewMemStorage())
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", advicerrors.AdviceHandler(hand.MetricAdd))
	return http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}
