// Code generated by MockGen. DO NOT EDIT.
// Source: ./service/toggle_disabler.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDisableToggle is a mock of DisableToggle interface.
type MockDisableToggle struct {
	ctrl     *gomock.Controller
	recorder *MockDisableToggleMockRecorder
}

// MockDisableToggleMockRecorder is the mock recorder for MockDisableToggle.
type MockDisableToggleMockRecorder struct {
	mock *MockDisableToggle
}

// NewMockDisableToggle creates a new mock instance.
func NewMockDisableToggle(ctrl *gomock.Controller) *MockDisableToggle {
	mock := &MockDisableToggle{ctrl: ctrl}
	mock.recorder = &MockDisableToggleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDisableToggle) EXPECT() *MockDisableToggleMockRecorder {
	return m.recorder
}

// Disable mocks base method.
func (m *MockDisableToggle) Disable(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Disable", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Disable indicates an expected call of Disable.
func (mr *MockDisableToggleMockRecorder) Disable(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disable", reflect.TypeOf((*MockDisableToggle)(nil).Disable), ctx, key)
}

// MockDisableToggleRepository is a mock of DisableToggleRepository interface.
type MockDisableToggleRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDisableToggleRepositoryMockRecorder
}

// MockDisableToggleRepositoryMockRecorder is the mock recorder for MockDisableToggleRepository.
type MockDisableToggleRepositoryMockRecorder struct {
	mock *MockDisableToggleRepository
}

// NewMockDisableToggleRepository creates a new mock instance.
func NewMockDisableToggleRepository(ctrl *gomock.Controller) *MockDisableToggleRepository {
	mock := &MockDisableToggleRepository{ctrl: ctrl}
	mock.recorder = &MockDisableToggleRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDisableToggleRepository) EXPECT() *MockDisableToggleRepositoryMockRecorder {
	return m.recorder
}

// Disable mocks base method.
func (m *MockDisableToggleRepository) Disable(ctx context.Context, key string, value bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Disable", ctx, key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// Disable indicates an expected call of Disable.
func (mr *MockDisableToggleRepositoryMockRecorder) Disable(ctx, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disable", reflect.TypeOf((*MockDisableToggleRepository)(nil).Disable), ctx, key, value)
}
