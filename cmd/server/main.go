package main

import (
	"github.com/zelas91/metric-collector/internal/advicerrors"
	"github.com/zelas91/metric-collector/internal/handlers"
	"github.com/zelas91/metric-collector/internal/storages"
	"net/http"
)

func main() {
	hand := handlers.NewUpdateHandler(storages.NewMemStorage())
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", advicerrors.Middleware(hand.UpdateMetric))
	http.ListenAndServe(":8080", mux)
}
