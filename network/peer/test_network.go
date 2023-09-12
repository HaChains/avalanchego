// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package peer

import (
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/proto/pb/p2p"
	"github.com/ava-labs/avalanchego/utils/ips"
)

var TestNetwork Network = testNetwork{}

type testNetwork struct{}

func (testNetwork) Connected(ids.GenericNodeID) {}

func (testNetwork) AllowConnection(ids.GenericNodeID) bool {
	return true
}

func (testNetwork) Track(ids.GenericNodeID, []*ips.ClaimedIPPort) ([]*p2p.PeerAck, error) {
	return nil, nil
}

func (testNetwork) MarkTracked(ids.GenericNodeID, []*p2p.PeerAck) error {
	return nil
}

func (testNetwork) Disconnected(ids.GenericNodeID) {}

func (testNetwork) Peers(ids.GenericNodeID) ([]ips.ClaimedIPPort, error) {
	return nil, nil
}
