// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package state

import (
	"crypto"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/staking"
	"github.com/ava-labs/avalanchego/upgrade"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/vms/proposervm/block"
)

func testBlockState(require *require.Assertions, bs BlockState) {
	parentID := ids.ID{1}
	timestamp := upgrade.InitiallyActiveTime
	pChainHeight := uint64(2)
	innerBlockBytes := []byte{3}
	chainID := ids.ID{4}
	networkID := uint32(5)

	tlsCert, err := staking.NewTLSCert()
	require.NoError(err)

	cert, err := staking.ParseCertificate(tlsCert.Leaf.Raw)
	require.NoError(err)
	key := tlsCert.PrivateKey.(crypto.Signer)

	var parentBlockSig []byte
	var blsSignKey *bls.SecretKey

	b, err := block.Build(
		parentID,
		timestamp,
		pChainHeight,
		cert,
		innerBlockBytes,
		chainID,
		key,
		block.NextBlockVRFSig(
			parentBlockSig,
			blsSignKey,
			chainID,
			networkID),
	)
	require.NoError(err)

	_, err = bs.GetBlock(b.ID())
	require.Equal(database.ErrNotFound, err)

	_, err = bs.GetBlock(b.ID())
	require.Equal(database.ErrNotFound, err)

	require.NoError(bs.PutBlock(b))

	fetchedBlock, err := bs.GetBlock(b.ID())
	require.NoError(err)
	require.Equal(b.Bytes(), fetchedBlock.Bytes())

	fetchedBlock, err = bs.GetBlock(b.ID())
	require.NoError(err)
	require.Equal(b.Bytes(), fetchedBlock.Bytes())
}

func TestBlockState(t *testing.T) {
	a := require.New(t)

	db := memdb.New()
	bs := NewBlockState(db)

	testBlockState(a, bs)
}

func TestMeteredBlockState(t *testing.T) {
	a := require.New(t)

	db := memdb.New()
	bs, err := NewMeteredBlockState(db, "", prometheus.NewRegistry())
	a.NoError(err)

	testBlockState(a, bs)
}
