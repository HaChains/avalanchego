// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package messages

import (
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
)

// RegisterSubnetValidator is signed when the subnet wants to add a validator with
// the given weight to the subnet.
type RegisterSubnetValidator struct {
	SubnetID  ids.ID                 `serialize:"true"`
	NodeID    ids.ID                 `serialize:"true"`
	Weight    uint64                 `serialize:"true"`
	BlsPubKey [bls.PublicKeyLen]byte `serialize:"true"`
	Expiry    uint64                 `serialize:"true"`

	bytes []byte
}

// NewRegisterSubnetValidator creates a new *RegisterSubnetValidator and initializes it.
func NewRegisterSubnetValidator(
	subnetID ids.ID,
	nodeID ids.ID,
	weight uint64,
	blsPubKey [bls.PublicKeyLen]byte,
	expiry uint64,
) (*RegisterSubnetValidator, error) {
	bhp := &RegisterSubnetValidator{
		SubnetID:  subnetID,
		NodeID:    nodeID,
		Weight:    weight,
		BlsPubKey: blsPubKey,
		Expiry:    expiry,
	}
	return bhp, initialize(bhp)
}

// ParseRegisterSubnetValidator converts a slice of bytes into an initialized RegisterSubnetValidator.
func ParseRegisterSubnetValidator(b []byte) (*RegisterSubnetValidator, error) {
	payloadIntf, err := Parse(b)
	if err != nil {
		return nil, err
	}
	payload, ok := payloadIntf.(*RegisterSubnetValidator)
	if !ok {
		return nil, fmt.Errorf("%w: %T", errWrongType, payloadIntf)
	}
	return payload, nil
}

// Bytes returns the binary representation of this payload. It assumes that the
// payload is initialized from either NewRegisterSubnetValidator or Parse.
func (b *RegisterSubnetValidator) Bytes() []byte {
	return b.bytes
}

func (b *RegisterSubnetValidator) initialize(bytes []byte) {
	b.bytes = bytes
}
