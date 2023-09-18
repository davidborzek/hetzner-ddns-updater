// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/davidborzek/hetzner-ddns-updater/pkg/hetzner (interfaces: Client)
//
// Generated by this command:
//
//	mockgen -package=mock -destination=mock_gen.go github.com/davidborzek/hetzner-ddns-updater/pkg/hetzner Client
//
// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	hetzner "github.com/davidborzek/hetzner-ddns-updater/pkg/hetzner"
	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// UpdateRecord mocks base method.
func (m *MockClient) UpdateRecord(arg0 string, arg1 hetzner.Record) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRecord", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRecord indicates an expected call of UpdateRecord.
func (mr *MockClientMockRecorder) UpdateRecord(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRecord", reflect.TypeOf((*MockClient)(nil).UpdateRecord), arg0, arg1)
}
