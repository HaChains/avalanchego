// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ava-labs/avalanchego/vms/proposervm/scheduler (interfaces: Scheduler)
//
// Generated by this command:
//
//	mockgen -package=scheduler -destination=vms/proposervm/scheduler/mock_scheduler.go github.com/ava-labs/avalanchego/vms/proposervm/scheduler Scheduler
//

// Package scheduler is a generated GoMock package.
package scheduler

import (
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockScheduler is a mock of Scheduler interface.
type MockScheduler struct {
	ctrl     *gomock.Controller
	recorder *MockSchedulerMockRecorder
}

// MockSchedulerMockRecorder is the mock recorder for MockScheduler.
type MockSchedulerMockRecorder struct {
	mock *MockScheduler
}

// NewMockScheduler creates a new mock instance.
func NewMockScheduler(ctrl *gomock.Controller) *MockScheduler {
	mock := &MockScheduler{ctrl: ctrl}
	mock.recorder = &MockSchedulerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScheduler) EXPECT() *MockSchedulerMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockScheduler) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockSchedulerMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockScheduler)(nil).Close))
}

// Dispatch mocks base method.
func (m *MockScheduler) Dispatch(arg0 time.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Dispatch", arg0)
}

// Dispatch indicates an expected call of Dispatch.
func (mr *MockSchedulerMockRecorder) Dispatch(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dispatch", reflect.TypeOf((*MockScheduler)(nil).Dispatch), arg0)
}

// SetBuildBlockTime mocks base method.
func (m *MockScheduler) SetBuildBlockTime(arg0 time.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetBuildBlockTime", arg0)
}

// SetBuildBlockTime indicates an expected call of SetBuildBlockTime.
func (mr *MockSchedulerMockRecorder) SetBuildBlockTime(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBuildBlockTime", reflect.TypeOf((*MockScheduler)(nil).SetBuildBlockTime), arg0)
}