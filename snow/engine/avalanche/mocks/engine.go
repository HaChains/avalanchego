// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	consensusavalanche "github.com/ava-labs/avalanchego/snow/consensus/avalanche"
	common "github.com/ava-labs/avalanchego/snow/engine/common"

	context "context"

	ids "github.com/ava-labs/avalanchego/ids"

	mock "github.com/stretchr/testify/mock"

	snow "github.com/ava-labs/avalanchego/snow"

	time "time"

	version "github.com/ava-labs/avalanchego/version"
)

// Engine is an autogenerated mock type for the Engine type
type Engine struct {
	mock.Mock
}

// Accepted provides a mock function with given fields: ctx, validatorID, requestID, containerIDs
func (_m *Engine) Accepted(ctx context.Context, validatorID ids.NodeID, requestID uint32, containerIDs []ids.ID) error {
	ret := _m.Called(ctx, validatorID, requestID, containerIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32, []ids.ID) error); ok {
		r0 = rf(ctx, validatorID, requestID, containerIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AcceptedFrontier provides a mock function with given fields: ctx, validatorID, requestID, containerIDs
func (_m *Engine) AcceptedFrontier(ctx context.Context, validatorID ids.NodeID, requestID uint32, containerIDs []ids.ID) error {
	ret := _m.Called(ctx, validatorID, requestID, containerIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32, []ids.ID) error); ok {
		r0 = rf(ctx, validatorID, requestID, containerIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AcceptedStateSummary provides a mock function with given fields: ctx, validatorID, requestID, summaryIDs
func (_m *Engine) AcceptedStateSummary(ctx context.Context, validatorID ids.NodeID, requestID uint32, summaryIDs []ids.ID) error {
	ret := _m.Called(ctx, validatorID, requestID, summaryIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32, []ids.ID) error); ok {
		r0 = rf(ctx, validatorID, requestID, summaryIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Ancestors provides a mock function with given fields: ctx, validatorID, requestID, containers
func (_m *Engine) Ancestors(ctx context.Context, validatorID ids.NodeID, requestID uint32, containers [][]byte) error {
	ret := _m.Called(ctx, validatorID, requestID, containers)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32, [][]byte) error); ok {
		r0 = rf(ctx, validatorID, requestID, containers)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AppGossip provides a mock function with given fields: ctx, nodeID, msg
func (_m *Engine) AppGossip(ctx context.Context, nodeID ids.NodeID, msg []byte) error {
	ret := _m.Called(ctx, nodeID, msg)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, []byte) error); ok {
		r0 = rf(ctx, nodeID, msg)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AppRequest provides a mock function with given fields: ctx, nodeID, requestID, deadline, request
func (_m *Engine) AppRequest(ctx context.Context, nodeID ids.NodeID, requestID uint32, deadline time.Time, request []byte) error {
	ret := _m.Called(ctx, nodeID, requestID, deadline, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32, time.Time, []byte) error); ok {
		r0 = rf(ctx, nodeID, requestID, deadline, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AppRequestFailed provides a mock function with given fields: ctx, nodeID, requestID
func (_m *Engine) AppRequestFailed(ctx context.Context, nodeID ids.NodeID, requestID uint32) error {
	ret := _m.Called(ctx, nodeID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32) error); ok {
		r0 = rf(ctx, nodeID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// AppResponse provides a mock function with given fields: ctx, nodeID, requestID, response
func (_m *Engine) AppResponse(ctx context.Context, nodeID ids.NodeID, requestID uint32, response []byte) error {
	ret := _m.Called(ctx, nodeID, requestID, response)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32, []byte) error); ok {
		r0 = rf(ctx, nodeID, requestID, response)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Chits provides a mock function with given fields: ctx, validatorID, requestID, containerIDs
func (_m *Engine) Chits(ctx context.Context, validatorID ids.NodeID, requestID uint32, containerIDs []ids.ID) error {
	ret := _m.Called(ctx, validatorID, requestID, containerIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32, []ids.ID) error); ok {
		r0 = rf(ctx, validatorID, requestID, containerIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Connected provides a mock function with given fields: id, nodeVersion
func (_m *Engine) Connected(id ids.NodeID, nodeVersion *version.Application) error {
	ret := _m.Called(id, nodeVersion)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.NodeID, *version.Application) error); ok {
		r0 = rf(id, nodeVersion)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Context provides a mock function with given fields:
func (_m *Engine) Context() *snow.ConsensusContext {
	ret := _m.Called()

	var r0 *snow.ConsensusContext
	if rf, ok := ret.Get(0).(func() *snow.ConsensusContext); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*snow.ConsensusContext)
		}
	}

	return r0
}

// CrossChainAppRequest provides a mock function with given fields: ctx, chainID, requestID, deadline, request
func (_m *Engine) CrossChainAppRequest(ctx context.Context, chainID ids.ID, requestID uint32, deadline time.Time, request []byte) error {
	ret := _m.Called(ctx, chainID, requestID, deadline, request)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.ID, uint32, time.Time, []byte) error); ok {
		r0 = rf(ctx, chainID, requestID, deadline, request)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CrossChainAppRequestFailed provides a mock function with given fields: ctx, chainID, requestID
func (_m *Engine) CrossChainAppRequestFailed(ctx context.Context, chainID ids.ID, requestID uint32) error {
	ret := _m.Called(ctx, chainID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.ID, uint32) error); ok {
		r0 = rf(ctx, chainID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CrossChainAppResponse provides a mock function with given fields: ctx, chainID, requestID, response
func (_m *Engine) CrossChainAppResponse(ctx context.Context, chainID ids.ID, requestID uint32, response []byte) error {
	ret := _m.Called(ctx, chainID, requestID, response)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.ID, uint32, []byte) error); ok {
		r0 = rf(ctx, chainID, requestID, response)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Disconnected provides a mock function with given fields: id
func (_m *Engine) Disconnected(id ids.NodeID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(ids.NodeID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: ctx, validatorID, requestID, containerID
func (_m *Engine) Get(ctx context.Context, validatorID ids.NodeID, requestID uint32, containerID ids.ID) error {
	ret := _m.Called(ctx, validatorID, requestID, containerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32, ids.ID) error); ok {
		r0 = rf(ctx, validatorID, requestID, containerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAccepted provides a mock function with given fields: ctx, validatorID, requestID, containerIDs
func (_m *Engine) GetAccepted(ctx context.Context, validatorID ids.NodeID, requestID uint32, containerIDs []ids.ID) error {
	ret := _m.Called(ctx, validatorID, requestID, containerIDs)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32, []ids.ID) error); ok {
		r0 = rf(ctx, validatorID, requestID, containerIDs)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAcceptedFailed provides a mock function with given fields: ctx, validatorID, requestID
func (_m *Engine) GetAcceptedFailed(ctx context.Context, validatorID ids.NodeID, requestID uint32) error {
	ret := _m.Called(ctx, validatorID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32) error); ok {
		r0 = rf(ctx, validatorID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAcceptedFrontier provides a mock function with given fields: ctx, validatorID, requestID
func (_m *Engine) GetAcceptedFrontier(ctx context.Context, validatorID ids.NodeID, requestID uint32) error {
	ret := _m.Called(ctx, validatorID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32) error); ok {
		r0 = rf(ctx, validatorID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAcceptedFrontierFailed provides a mock function with given fields: ctx, validatorID, requestID
func (_m *Engine) GetAcceptedFrontierFailed(ctx context.Context, validatorID ids.NodeID, requestID uint32) error {
	ret := _m.Called(ctx, validatorID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32) error); ok {
		r0 = rf(ctx, validatorID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAcceptedStateSummary provides a mock function with given fields: ctx, validatorID, requestID, keys
func (_m *Engine) GetAcceptedStateSummary(ctx context.Context, validatorID ids.NodeID, requestID uint32, keys []uint64) error {
	ret := _m.Called(ctx, validatorID, requestID, keys)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32, []uint64) error); ok {
		r0 = rf(ctx, validatorID, requestID, keys)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAcceptedStateSummaryFailed provides a mock function with given fields: ctx, validatorID, requestID
func (_m *Engine) GetAcceptedStateSummaryFailed(ctx context.Context, validatorID ids.NodeID, requestID uint32) error {
	ret := _m.Called(ctx, validatorID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32) error); ok {
		r0 = rf(ctx, validatorID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAncestors provides a mock function with given fields: ctx, validatorID, requestID, containerID
func (_m *Engine) GetAncestors(ctx context.Context, validatorID ids.NodeID, requestID uint32, containerID ids.ID) error {
	ret := _m.Called(ctx, validatorID, requestID, containerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32, ids.ID) error); ok {
		r0 = rf(ctx, validatorID, requestID, containerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAncestorsFailed provides a mock function with given fields: ctx, validatorID, requestID
func (_m *Engine) GetAncestorsFailed(ctx context.Context, validatorID ids.NodeID, requestID uint32) error {
	ret := _m.Called(ctx, validatorID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32) error); ok {
		r0 = rf(ctx, validatorID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetFailed provides a mock function with given fields: ctx, validatorID, requestID
func (_m *Engine) GetFailed(ctx context.Context, validatorID ids.NodeID, requestID uint32) error {
	ret := _m.Called(ctx, validatorID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32) error); ok {
		r0 = rf(ctx, validatorID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetStateSummaryFrontier provides a mock function with given fields: ctx, validatorID, requestID
func (_m *Engine) GetStateSummaryFrontier(ctx context.Context, validatorID ids.NodeID, requestID uint32) error {
	ret := _m.Called(ctx, validatorID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32) error); ok {
		r0 = rf(ctx, validatorID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetStateSummaryFrontierFailed provides a mock function with given fields: ctx, validatorID, requestID
func (_m *Engine) GetStateSummaryFrontierFailed(ctx context.Context, validatorID ids.NodeID, requestID uint32) error {
	ret := _m.Called(ctx, validatorID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32) error); ok {
		r0 = rf(ctx, validatorID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetVM provides a mock function with given fields:
func (_m *Engine) GetVM() common.VM {
	ret := _m.Called()

	var r0 common.VM
	if rf, ok := ret.Get(0).(func() common.VM); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.VM)
		}
	}

	return r0
}

// GetVtx provides a mock function with given fields: vtxID
func (_m *Engine) GetVtx(vtxID ids.ID) (consensusavalanche.Vertex, error) {
	ret := _m.Called(vtxID)

	var r0 consensusavalanche.Vertex
	if rf, ok := ret.Get(0).(func(ids.ID) consensusavalanche.Vertex); ok {
		r0 = rf(vtxID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(consensusavalanche.Vertex)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(ids.ID) error); ok {
		r1 = rf(vtxID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Gossip provides a mock function with given fields:
func (_m *Engine) Gossip() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Halt provides a mock function with given fields:
func (_m *Engine) Halt() {
	_m.Called()
}

// HealthCheck provides a mock function with given fields:
func (_m *Engine) HealthCheck() (interface{}, error) {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Notify provides a mock function with given fields: _a0
func (_m *Engine) Notify(_a0 common.Message) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Message) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PullQuery provides a mock function with given fields: ctx, validatorID, requestID, containerID
func (_m *Engine) PullQuery(ctx context.Context, validatorID ids.NodeID, requestID uint32, containerID ids.ID) error {
	ret := _m.Called(ctx, validatorID, requestID, containerID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32, ids.ID) error); ok {
		r0 = rf(ctx, validatorID, requestID, containerID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PushQuery provides a mock function with given fields: ctx, validatorID, requestID, container
func (_m *Engine) PushQuery(ctx context.Context, validatorID ids.NodeID, requestID uint32, container []byte) error {
	ret := _m.Called(ctx, validatorID, requestID, container)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32, []byte) error); ok {
		r0 = rf(ctx, validatorID, requestID, container)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Put provides a mock function with given fields: ctx, validatorID, requestID, container
func (_m *Engine) Put(ctx context.Context, validatorID ids.NodeID, requestID uint32, container []byte) error {
	ret := _m.Called(ctx, validatorID, requestID, container)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32, []byte) error); ok {
		r0 = rf(ctx, validatorID, requestID, container)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// QueryFailed provides a mock function with given fields: ctx, validatorID, requestID
func (_m *Engine) QueryFailed(ctx context.Context, validatorID ids.NodeID, requestID uint32) error {
	ret := _m.Called(ctx, validatorID, requestID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32) error); ok {
		r0 = rf(ctx, validatorID, requestID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Shutdown provides a mock function with given fields:
func (_m *Engine) Shutdown() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields: startReqID
func (_m *Engine) Start(startReqID uint32) error {
	ret := _m.Called(startReqID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint32) error); ok {
		r0 = rf(startReqID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StateSummaryFrontier provides a mock function with given fields: ctx, validatorID, requestID, summary
func (_m *Engine) StateSummaryFrontier(ctx context.Context, validatorID ids.NodeID, requestID uint32, summary []byte) error {
	ret := _m.Called(ctx, validatorID, requestID, summary)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ids.NodeID, uint32, []byte) error); ok {
		r0 = rf(ctx, validatorID, requestID, summary)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Timeout provides a mock function with given fields:
func (_m *Engine) Timeout() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewEngine interface {
	mock.TestingT
	Cleanup(func())
}

// NewEngine creates a new instance of Engine. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewEngine(t mockConstructorTestingTNewEngine) *Engine {
	mock := &Engine{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
