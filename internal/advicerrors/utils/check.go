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
func CheckUpdateMetric(method string, parts []string) error {
	if method != http.MethodPost {
		return advicerrors.ErrMethodNotAllowed
	}
	if len(parts) < 5 {
		return advicerrors.ErrNotFound
	}
	metricTypeStr := parts[2]
	value := parts[4]
	if !isType(metricTypeStr) {
		return advicerrors.ErrBadRequest
	}
	if !isValue(value) {
		return advicerrors.ErrBadRequest
	}
	return nil
}
