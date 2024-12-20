// Code generated by MockGen. DO NOT EDIT.
// Source: module.go

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockHandler is a mock of Handler interface.
type MockHandler struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerMockRecorder
}

// MockHandlerMockRecorder is the mock recorder for MockHandler.
type MockHandlerMockRecorder struct {
	mock *MockHandler
}

// NewMockHandler creates a new mock instance.
func NewMockHandler(ctrl *gomock.Controller) *MockHandler {
	mock := &MockHandler{ctrl: ctrl}
	mock.recorder = &MockHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandler) EXPECT() *MockHandlerMockRecorder {
	return m.recorder
}

// Register mocks base method.
func (m *MockHandler) Register(arg0 *gin.Engine) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Register", arg0)
}

// Register indicates an expected call of Register.
func (mr *MockHandlerMockRecorder) Register(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockHandler)(nil).Register), arg0)
}

// MocksessionService is a mock of sessionService interface.
type MocksessionService struct {
	ctrl     *gomock.Controller
	recorder *MocksessionServiceMockRecorder
}

// MocksessionServiceMockRecorder is the mock recorder for MocksessionService.
type MocksessionServiceMockRecorder struct {
	mock *MocksessionService
}

// NewMocksessionService creates a new mock instance.
func NewMocksessionService(ctrl *gomock.Controller) *MocksessionService {
	mock := &MocksessionService{ctrl: ctrl}
	mock.recorder = &MocksessionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocksessionService) EXPECT() *MocksessionServiceMockRecorder {
	return m.recorder
}

// CheckSession mocks base method.
func (m *MocksessionService) CheckSession(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckSession", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckSession indicates an expected call of CheckSession.
func (mr *MocksessionServiceMockRecorder) CheckSession(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckSession", reflect.TypeOf((*MocksessionService)(nil).CheckSession), arg0)
}

// CreateSession mocks base method.
func (m *MocksessionService) CreateSession() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession")
	ret0, _ := ret[0].(string)
	return ret0
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MocksessionServiceMockRecorder) CreateSession() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MocksessionService)(nil).CreateSession))
}

// Delete mocks base method.
func (m *MocksessionService) Delete(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MocksessionServiceMockRecorder) Delete(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MocksessionService)(nil).Delete), key)
}
