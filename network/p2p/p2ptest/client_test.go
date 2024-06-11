// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package p2ptest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/network/p2p"
	"github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/utils/set"
)

func TestNewClient_AppGossip(t *testing.T) {
	require := require.New(t)
	ctx := context.Background()

	appGossipChan := make(chan struct{})
	testHandler := p2p.TestHandler{
		AppGossipF: func(ctx context.Context, nodeID ids.NodeID, gossipBytes []byte) {
			close(appGossipChan)
		},
	}

	client := NewClient(t, ctx, testHandler)
	require.NoError(client.AppGossip(ctx, common.SendConfig{}, []byte("foobar")))
	<-appGossipChan
}

func TestNewClient_AppRequest(t *testing.T) {
	tests := []struct {
		name        string
		appResponse []byte
		appErr      error
		appRequestF func(ctx context.Context, client *p2p.Client, onResponse p2p.AppResponseCallback) error
	}{
		{
			name:        "AppRequest - response",
			appResponse: []byte("foobar"),
			appRequestF: func(ctx context.Context, client *p2p.Client, onResponse p2p.AppResponseCallback) error {
				return client.AppRequest(ctx, set.Of(ids.GenerateTestNodeID()), []byte("foo"), onResponse)
			},
		},
		{
			name:   "AppRequest - error",
			appErr: errors.New("foobar"),
			appRequestF: func(ctx context.Context, client *p2p.Client, onResponse p2p.AppResponseCallback) error {
				return client.AppRequest(ctx, set.Of(ids.GenerateTestNodeID()), []byte("foo"), onResponse)
			},
		},
		{
			name:        "AppRequestAny - response",
			appResponse: []byte("foobar"),
			appRequestF: func(ctx context.Context, client *p2p.Client, onResponse p2p.AppResponseCallback) error {
				return client.AppRequestAny(ctx, []byte("foo"), onResponse)
			},
		},
		{
			name:   "AppRequestAny - error",
			appErr: errors.New("foobar"),
			appRequestF: func(ctx context.Context, client *p2p.Client, onResponse p2p.AppResponseCallback) error {
				return client.AppRequestAny(ctx, []byte("foo"), onResponse)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO remove when AppErrors are supported
			if tt.appErr != nil {
				t.Skip("sending app errors not supported yet")
			}

			require := require.New(t)
			ctx := context.Background()

			appRequestChan := make(chan struct{})
			testHandler := p2p.TestHandler{
				AppRequestF: func(context.Context, ids.NodeID, time.Time, []byte) ([]byte, error) {
					return tt.appResponse, tt.appErr
				},
			}

			client := NewClient(t, ctx, testHandler)
			require.NoError(tt.appRequestF(
				ctx,
				client,
				func(_ context.Context, _ ids.NodeID, responseBytes []byte, err error) {
					require.Equal(tt.appErr, err)
					require.Equal(tt.appResponse, responseBytes)
					close(appRequestChan)
				},
			))
			<-appRequestChan
		})
	}
}
