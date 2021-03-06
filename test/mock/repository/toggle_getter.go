// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/repository/toggle_getter.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	entity "github.com/indrasaputra/toggle/entity"
)

// MockGetToggleDatabase is a mock of GetToggleDatabase interface.
type MockGetToggleDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockGetToggleDatabaseMockRecorder
}

// MockGetToggleDatabaseMockRecorder is the mock recorder for MockGetToggleDatabase.
type MockGetToggleDatabaseMockRecorder struct {
	mock *MockGetToggleDatabase
}

// NewMockGetToggleDatabase creates a new mock instance.
func NewMockGetToggleDatabase(ctrl *gomock.Controller) *MockGetToggleDatabase {
	mock := &MockGetToggleDatabase{ctrl: ctrl}
	mock.recorder = &MockGetToggleDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGetToggleDatabase) EXPECT() *MockGetToggleDatabaseMockRecorder {
	return m.recorder
}

// GetAll mocks base method.
func (m *MockGetToggleDatabase) GetAll(ctx context.Context, limit uint) ([]*entity.Toggle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, limit)
	ret0, _ := ret[0].([]*entity.Toggle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockGetToggleDatabaseMockRecorder) GetAll(ctx, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockGetToggleDatabase)(nil).GetAll), ctx, limit)
}

// GetByKey mocks base method.
func (m *MockGetToggleDatabase) GetByKey(ctx context.Context, key string) (*entity.Toggle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByKey", ctx, key)
	ret0, _ := ret[0].(*entity.Toggle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByKey indicates an expected call of GetByKey.
func (mr *MockGetToggleDatabaseMockRecorder) GetByKey(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByKey", reflect.TypeOf((*MockGetToggleDatabase)(nil).GetByKey), ctx, key)
}

// MockGetToggleCache is a mock of GetToggleCache interface.
type MockGetToggleCache struct {
	ctrl     *gomock.Controller
	recorder *MockGetToggleCacheMockRecorder
}

// MockGetToggleCacheMockRecorder is the mock recorder for MockGetToggleCache.
type MockGetToggleCacheMockRecorder struct {
	mock *MockGetToggleCache
}

// NewMockGetToggleCache creates a new mock instance.
func NewMockGetToggleCache(ctrl *gomock.Controller) *MockGetToggleCache {
	mock := &MockGetToggleCache{ctrl: ctrl}
	mock.recorder = &MockGetToggleCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGetToggleCache) EXPECT() *MockGetToggleCacheMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockGetToggleCache) Get(ctx context.Context, key string) (*entity.Toggle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key)
	ret0, _ := ret[0].(*entity.Toggle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockGetToggleCacheMockRecorder) Get(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockGetToggleCache)(nil).Get), ctx, key)
}

// Set mocks base method.
func (m *MockGetToggleCache) Set(ctx context.Context, toggle *entity.Toggle) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, toggle)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockGetToggleCacheMockRecorder) Set(ctx, toggle interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockGetToggleCache)(nil).Set), ctx, toggle)
}
