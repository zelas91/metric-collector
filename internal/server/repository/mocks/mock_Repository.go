// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	repository "github.com/zelas91/metric-collector/internal/server/repository"
)

// MockStorageRepository is a mock of StorageRepository interface.
type MockStorageRepository struct {
	ctrl     *gomock.Controller
	recorder *MockStorageRepositoryMockRecorder
}

// MockStorageRepositoryMockRecorder is the mock recorder for MockStorageRepository.
type MockStorageRepositoryMockRecorder struct {
	mock *MockStorageRepository
}

// NewMockStorageRepository creates a new mock instance.
func NewMockStorageRepository(ctrl *gomock.Controller) *MockStorageRepository {
	mock := &MockStorageRepository{ctrl: ctrl}
	mock.recorder = &MockStorageRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageRepository) EXPECT() *MockStorageRepositoryMockRecorder {
	return m.recorder
}

// AddMetric mocks base method.
func (m *MockStorageRepository) AddMetric(metrics repository.Metric) *repository.Metric {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMetric", metrics)
	ret0, _ := ret[0].(*repository.Metric)
	return ret0
}

// AddMetric indicates an expected call of AddMetric.
func (mr *MockStorageRepositoryMockRecorder) AddMetric(metrics interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMetric", reflect.TypeOf((*MockStorageRepository)(nil).AddMetric), metrics)
}

// GetMetric mocks base method.
func (m *MockStorageRepository) GetMetric(name string) (*repository.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetric", name)
	ret0, _ := ret[0].(*repository.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetric indicates an expected call of GetMetric.
func (mr *MockStorageRepositoryMockRecorder) GetMetric(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetric", reflect.TypeOf((*MockStorageRepository)(nil).GetMetric), name)
}

// GetMetrics mocks base method.
func (m *MockStorageRepository) GetMetrics() []repository.Metric {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetrics")
	ret0, _ := ret[0].([]repository.Metric)
	return ret0
}

// GetMetrics indicates an expected call of GetMetrics.
func (mr *MockStorageRepositoryMockRecorder) GetMetrics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetrics", reflect.TypeOf((*MockStorageRepository)(nil).GetMetrics))
}
