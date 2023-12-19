// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	repository "github.com/zelas91/metric-collector/internal/server/repository"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AddMetric mocks base method.
func (m *MockService) AddMetric(ctx context.Context, name, mType, value string) (*repository.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMetric", ctx, name, mType, value)
	ret0, _ := ret[0].(*repository.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMetric indicates an expected call of AddMetric.
func (mr *MockServiceMockRecorder) AddMetric(ctx, name, mType, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMetric", reflect.TypeOf((*MockService)(nil).AddMetric), ctx, name, mType, value)
}

// AddMetricJSON mocks base method.
func (m *MockService) AddMetricJSON(ctx context.Context, metric repository.Metric) (*repository.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMetricJSON", ctx, metric)
	ret0, _ := ret[0].(*repository.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddMetricJSON indicates an expected call of AddMetricJSON.
func (mr *MockServiceMockRecorder) AddMetricJSON(ctx, metric interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMetricJSON", reflect.TypeOf((*MockService)(nil).AddMetricJSON), ctx, metric)
}

// AddMetrics mocks base method.
func (m *MockService) AddMetrics(ctx context.Context, metrics []repository.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMetrics", ctx, metrics)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddMetrics indicates an expected call of AddMetrics.
func (mr *MockServiceMockRecorder) AddMetrics(ctx, metrics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMetrics", reflect.TypeOf((*MockService)(nil).AddMetrics), ctx, metrics)
}

// GetMetric mocks base method.
func (m *MockService) GetMetric(ctx context.Context, name string) (*repository.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetric", ctx, name)
	ret0, _ := ret[0].(*repository.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetric indicates an expected call of GetMetric.
func (mr *MockServiceMockRecorder) GetMetric(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetric", reflect.TypeOf((*MockService)(nil).GetMetric), ctx, name)
}

// GetMetrics mocks base method.
func (m *MockService) GetMetrics(ctx context.Context) []repository.Metric {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetrics", ctx)
	ret0, _ := ret[0].([]repository.Metric)
	return ret0
}

// GetMetrics indicates an expected call of GetMetrics.
func (mr *MockServiceMockRecorder) GetMetrics(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetrics", reflect.TypeOf((*MockService)(nil).GetMetrics), ctx)
}
