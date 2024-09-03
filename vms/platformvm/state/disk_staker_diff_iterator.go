// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"encoding/binary"
	"fmt"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/ids"
)

const (
	// startDiffKey = [subnetID] + [inverseHeight]
	startDiffKeyLength = ids.IDLen + database.Uint64Size
	// diffKey = [subnetID] + [inverseHeight] + [shortNodeID]/[NodeID]
	minDiffKeyLength = startDiffKeyLength + ids.ShortNodeIDLen
	// diffKeyNodeIDOffset = [subnetIDLen] + [inverseHeightLen]
	diffKeyNodeIDOffset = ids.IDLen + database.Uint64Size

	// weightValue = [isNegative] + [weight]
	weightValueLength = database.BoolSize + database.Uint64Size
)

var (
	errUnexpectedDiffKeyLength     = fmt.Errorf("expected min diff key length of %d", minDiffKeyLength)
	errUnexpectedWeightValueLength = fmt.Errorf("expected weight value length %d", weightValueLength)
)

// marshalStartDiffKey is used to determine the starting key when iterating.
//
// Invariant: the result is a prefix of [marshalDiffKey] when called with the
// same arguments.
func marshalStartDiffKey(subnetID ids.ID, height uint64) []byte {
	key := make([]byte, startDiffKeyLength)
	copy(key, subnetID[:])
	packIterableHeight(key[ids.IDLen:], height)
	return key
}

func marshalDiffKey(subnetID ids.ID, height uint64, nodeID ids.NodeID) []byte {
	nodeIDBytes := nodeID.Bytes()
	keyLength := startDiffKeyLength + len(nodeIDBytes)
	key := make([]byte, keyLength)
	copy(key, subnetID[:])
	packIterableHeight(key[ids.IDLen:], height)
	copy(key[diffKeyNodeIDOffset:], nodeIDBytes)
	return key
}

func unmarshalDiffKey(key []byte) (ids.ID, uint64, ids.NodeID, error) {
	if len(key) < minDiffKeyLength {
		return ids.Empty, 0, ids.EmptyNodeID, errUnexpectedDiffKeyLength
	}
	var (
		subnetID ids.ID
		nodeID   ids.NodeID
		err      error
	)
	nodeID, err = ids.ParseNodeID(key[diffKeyNodeIDOffset:])
	if err != nil {
		return ids.Empty, 0, ids.EmptyNodeID, errUnexpectedDiffKeyLength
	}
	copy(subnetID[:], key)
	height := unpackIterableHeight(key[ids.IDLen:])
	return subnetID, height, nodeID, err
}

func marshalWeightDiff(diff *ValidatorWeightDiff) []byte {
	value := make([]byte, weightValueLength)
	if diff.Decrease {
		value[0] = database.BoolTrue
	}
	binary.BigEndian.PutUint64(value[database.BoolSize:], diff.Amount)
	return value
}

func unmarshalWeightDiff(value []byte) (*ValidatorWeightDiff, error) {
	if len(value) != weightValueLength {
		return nil, errUnexpectedWeightValueLength
	}
	return &ValidatorWeightDiff{
		Decrease: value[0] == database.BoolTrue,
		Amount:   binary.BigEndian.Uint64(value[database.BoolSize:]),
	}, nil
}

// Note: [height] is encoded as a bit flipped big endian number so that
// iterating lexicographically results in iterating in decreasing heights.
//
// Invariant: [key] has sufficient length
func packIterableHeight(key []byte, height uint64) {
	binary.BigEndian.PutUint64(key, ^height)
}

// Because we bit flip the height when constructing the key, we must remember to
// bip flip again here.
//
// Invariant: [key] has sufficient length
func unpackIterableHeight(key []byte) uint64 {
	return ^binary.BigEndian.Uint64(key)
}
