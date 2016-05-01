// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/smancke/dockroute/dockrouted (interfaces: Backend)

package dockrouted

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of Backend interface
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *_MockBackendRecorder
}

// Recorder for MockBackend (not exported)
type _MockBackendRecorder struct {
	mock *MockBackend
}

func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &_MockBackendRecorder{mock}
	return mock
}

func (_m *MockBackend) EXPECT() *_MockBackendRecorder {
	return _m.recorder
}

func (_m *MockBackend) GetService(_param0 string) (string, error) {
	ret := _m.ctrl.Call(_m, "GetService", _param0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockBackendRecorder) GetService(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetService", arg0)
}
