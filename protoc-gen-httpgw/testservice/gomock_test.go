// Automatically generated by MockGen. DO NOT EDIT!
// Source: srv.httpgw.go

package testservice

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
)

// Mock of SrvSetterInterface interface
type MockSrvSetterInterface struct {
	ctrl     *gomock.Controller
	recorder *_MockSrvSetterInterfaceRecorder
}

// Recorder for MockSrvSetterInterface (not exported)
type _MockSrvSetterInterfaceRecorder struct {
	mock *MockSrvSetterInterface
}

func NewMockSrvSetterInterface(ctrl *gomock.Controller) *MockSrvSetterInterface {
	mock := &MockSrvSetterInterface{ctrl: ctrl}
	mock.recorder = &_MockSrvSetterInterfaceRecorder{mock}
	return mock
}

func (_m *MockSrvSetterInterface) EXPECT() *_MockSrvSetterInterfaceRecorder {
	return _m.recorder
}

func (_m *MockSrvSetterInterface) Set(_param0 context.Context, _param1 *ReqSet) (*RespSet, error) {
	ret := _m.ctrl.Call(_m, "Set", _param0, _param1)
	ret0, _ := ret[0].(*RespSet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockSrvSetterInterfaceRecorder) Set(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Set", arg0, arg1)
}

// Mock of SrvWaiterInterface interface
type MockSrvWaiterInterface struct {
	ctrl     *gomock.Controller
	recorder *_MockSrvWaiterInterfaceRecorder
}

// Recorder for MockSrvWaiterInterface (not exported)
type _MockSrvWaiterInterfaceRecorder struct {
	mock *MockSrvWaiterInterface
}

func NewMockSrvWaiterInterface(ctrl *gomock.Controller) *MockSrvWaiterInterface {
	mock := &MockSrvWaiterInterface{ctrl: ctrl}
	mock.recorder = &_MockSrvWaiterInterfaceRecorder{mock}
	return mock
}

func (_m *MockSrvWaiterInterface) EXPECT() *_MockSrvWaiterInterfaceRecorder {
	return _m.recorder
}

func (_m *MockSrvWaiterInterface) Wait(_param0 context.Context, _param1 *ReqWait) (*Empty, error) {
	ret := _m.ctrl.Call(_m, "Wait", _param0, _param1)
	ret0, _ := ret[0].(*Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockSrvWaiterInterfaceRecorder) Wait(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Wait", arg0, arg1)
}
