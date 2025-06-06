// Code generated by MockGen. DO NOT EDIT.
// Source: ./pvz_handler.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	request "github.com/AnikinSimon/avito-test-backend/internal/models/dto/request"
	entity "github.com/AnikinSimon/avito-test-backend/internal/models/entity"
	gomock "github.com/golang/mock/gomock"
)

// MockPvzService is a mock of PvzService interface.
type MockPvzService struct {
	ctrl     *gomock.Controller
	recorder *MockPvzServiceMockRecorder
}

// MockPvzServiceMockRecorder is the mock recorder for MockPvzService.
type MockPvzServiceMockRecorder struct {
	mock *MockPvzService
}

// NewMockPvzService creates a new mock instance.
func NewMockPvzService(ctrl *gomock.Controller) *MockPvzService {
	mock := &MockPvzService{ctrl: ctrl}
	mock.recorder = &MockPvzServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPvzService) EXPECT() *MockPvzServiceMockRecorder {
	return m.recorder
}

// CreatePvz mocks base method.
func (m *MockPvzService) CreatePvz(arg0 context.Context, arg1 *request.CreatePvz) (*entity.Pvz, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePvz", arg0, arg1)
	ret0, _ := ret[0].(*entity.Pvz)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePvz indicates an expected call of CreatePvz.
func (mr *MockPvzServiceMockRecorder) CreatePvz(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePvz", reflect.TypeOf((*MockPvzService)(nil).CreatePvz), arg0, arg1)
}
