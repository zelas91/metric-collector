package utils

import (
	"github.com/zelas91/metric-collector/internal/advicerrors"
	"net/http"
	"strconv"
)

func isValue(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}
func isType(mType string) bool {
	switch mType {

	case "counter":
	case "gauge":
	default:
		return false
	}
	return true

}
func CheckUpdateMetric(method string, parts []string) *advicerrors.AppError {
	if method != http.MethodPost {
		return advicerrors.NewErrMethodNotAllowed("method not allowed")
	}
	if len(parts) < 5 {
		return advicerrors.NewErrNotFound("not found")
	}
	metricTypeStr := parts[2]
	value := parts[4]
	if !isType(metricTypeStr) {
		return advicerrors.NewErrBadRequest("bad request")
	}
	if !isValue(value) {
		return advicerrors.NewErrBadRequest("bad request")
	}
	return nil
}
