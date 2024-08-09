// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package block

import (
	"fmt"
	"sync"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/upgrade"
)

type ParseResult struct {
	Block Block
	Err   error
}

// ParseBlocks parses the given raw blocks into tuples of (Block, error).
// Each ParseResult is returned in the same order as its corresponding bytes in the input.
func ParseBlocks(blks [][]byte, upgradeConfig upgrade.Config, chainID ids.ID) []ParseResult {
	results := make([]ParseResult, len(blks))

	var wg sync.WaitGroup
	wg.Add(len(blks))

	for i, blk := range blks {
		go func(i int, blkBytes []byte) {
			defer wg.Done()
			results[i].Block, results[i].Err = Parse(blkBytes, upgradeConfig, chainID)
		}(i, blk)
	}

	wg.Wait()

	return results
}

// Parse a block and verify that the signature attached to the block is valid
// for the certificate provided in the block.
func Parse(bytes []byte, upgradeConfig upgrade.Config, chainID ids.ID) (Block, error) {
	block, err := ParseWithoutVerification(bytes, upgradeConfig)
	if err != nil {
		return nil, err
	}
	return block, block.verify(chainID)
}

// ParseWithoutVerification parses a block without verifying that the signature
// on the block is correct.
func ParseWithoutVerification(bytes []byte, upgradeConfig upgrade.Config) (Block, error) {
	var block Block
	parsedVersion, err := Codec.Unmarshal(bytes, &block)
	if err != nil {
		return nil, err
	}
	if parsedVersion > CodecVersion {
		return nil, fmt.Errorf("expected codec version up to %d but got %d", CodecVersion, parsedVersion)
	}
	if err = block.initialize(bytes); err != nil {
		return block, err
	}
	if statelessBlock, ok := block.(*statelessBlock); ok {
		if !upgradeConfig.IsEtnaActivated(statelessBlock.Timestamp()) && parsedVersion == CodecVersion {
			return nil, fmt.Errorf("stateless block encoded using a codec version %d that is not yet enabled", CodecVersion)
		}
	}

	return block, nil
}

func ParseHeader(bytes []byte) (Header, error) {
	header := statelessHeader{}
	parsedVersion, err := Codec.Unmarshal(bytes, &header)
	if err != nil {
		return nil, err
	}
	if parsedVersion != CodecVersionV0 {
		return nil, fmt.Errorf("expected codec version %d but got %d", CodecVersionV0, parsedVersion)
	}
	header.bytes = bytes
	return &header, nil
}
