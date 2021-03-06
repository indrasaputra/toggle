// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/toggle_deleter.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	entity "github.com/indrasaputra/toggle/entity"
)

// MockDeleteToggleDatabase is a mock of DeleteToggleDatabase interface.
type MockDeleteToggleDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDeleteToggleDatabaseMockRecorder
}

// MockDeleteToggleDatabaseMockRecorder is the mock recorder for MockDeleteToggleDatabase.
type MockDeleteToggleDatabaseMockRecorder struct {
	mock *MockDeleteToggleDatabase
}

// NewMockDeleteToggleDatabase creates a new mock instance.
func NewMockDeleteToggleDatabase(ctrl *gomock.Controller) *MockDeleteToggleDatabase {
	mock := &MockDeleteToggleDatabase{ctrl: ctrl}
	mock.recorder = &MockDeleteToggleDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeleteToggleDatabase) EXPECT() *MockDeleteToggleDatabaseMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockDeleteToggleDatabase) Delete(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDeleteToggleDatabaseMockRecorder) Delete(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDeleteToggleDatabase)(nil).Delete), ctx, key)
}

// GetByKey mocks base method.
func (m *MockDeleteToggleDatabase) GetByKey(ctx context.Context, key string) (*entity.Toggle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByKey", ctx, key)
	ret0, _ := ret[0].(*entity.Toggle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByKey indicates an expected call of GetByKey.
func (mr *MockDeleteToggleDatabaseMockRecorder) GetByKey(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByKey", reflect.TypeOf((*MockDeleteToggleDatabase)(nil).GetByKey), ctx, key)
}

// MockDeleteToggleCache is a mock of DeleteToggleCache interface.
type MockDeleteToggleCache struct {
	ctrl     *gomock.Controller
	recorder *MockDeleteToggleCacheMockRecorder
}

// MockDeleteToggleCacheMockRecorder is the mock recorder for MockDeleteToggleCache.
type MockDeleteToggleCacheMockRecorder struct {
	mock *MockDeleteToggleCache
}

// NewMockDeleteToggleCache creates a new mock instance.
func NewMockDeleteToggleCache(ctrl *gomock.Controller) *MockDeleteToggleCache {
	mock := &MockDeleteToggleCache{ctrl: ctrl}
	mock.recorder = &MockDeleteToggleCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeleteToggleCache) EXPECT() *MockDeleteToggleCacheMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockDeleteToggleCache) Delete(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockDeleteToggleCacheMockRecorder) Delete(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockDeleteToggleCache)(nil).Delete), ctx, key)
}
