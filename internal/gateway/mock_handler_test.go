// Code generated by MockGen. DO NOT EDIT.
// Source: internal/gateway/handler.go

// Package gateway_test is a generated GoMock package.
package gateway_test

import (
	gateway "broozkan/api-gateway/internal/gateway"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHandlerService is a mock of HandlerService interface.
type MockHandlerService struct {
	ctrl     *gomock.Controller
	recorder *MockHandlerServiceMockRecorder
}

// MockHandlerServiceMockRecorder is the mock recorder for MockHandlerService.
type MockHandlerServiceMockRecorder struct {
	mock *MockHandlerService
}

// NewMockHandlerService creates a new mock instance.
func NewMockHandlerService(ctrl *gomock.Controller) *MockHandlerService {
	mock := &MockHandlerService{ctrl: ctrl}
	mock.recorder = &MockHandlerServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandlerService) EXPECT() *MockHandlerServiceMockRecorder {
	return m.recorder
}

// Forward mocks base method.
func (m *MockHandlerService) Forward(ctx context.Context, body interface{}, serviceProvider gateway.Router) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Forward", ctx, body, serviceProvider)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Forward indicates an expected call of Forward.
func (mr *MockHandlerServiceMockRecorder) Forward(ctx, body, serviceProvider interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Forward", reflect.TypeOf((*MockHandlerService)(nil).Forward), ctx, body, serviceProvider)
}

// ResolveService mocks base method.
func (m *MockHandlerService) ResolveService(service string) gateway.Router {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveService", service)
	ret0, _ := ret[0].(gateway.Router)
	return ret0
}

// ResolveService indicates an expected call of ResolveService.
func (mr *MockHandlerServiceMockRecorder) ResolveService(service interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveService", reflect.TypeOf((*MockHandlerService)(nil).ResolveService), service)
}