// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package proposervm

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/database/memdb"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/ava-labs/avalanchego/snow/choices"
	"github.com/ava-labs/avalanchego/snow/consensus/snowman"
	"github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/snow/engine/snowman/block"
	"github.com/ava-labs/avalanchego/vms/proposervm/summary"

	statelessblock "github.com/ava-labs/avalanchego/vms/proposervm/block"
)

// Post Fork section
func TestBlockBackfillEnabledPostFork(t *testing.T) {
	require := require.New(t)
	toEngineCh := make(chan common.Message)
	innerVM, vm := setupBlockBackfillingVM(t, toEngineCh)
	defer func() {
		require.NoError(vm.Shutdown(context.Background()))
	}()

	// 1. Accept a State summary
	var (
		forkHeight                 = uint64(100)
		stateSummaryHeight         = uint64(2023)
		proVMParentStateSummaryBlk = ids.GenerateTestID()
	)

	innerSummary := &block.TestStateSummary{
		IDV:     ids.ID{'s', 'u', 'm', 'm', 'a', 'r', 'y', 'I', 'D'},
		HeightV: stateSummaryHeight,
		BytesV:  []byte{'i', 'n', 'n', 'e', 'r'},
	}
	innerStateSyncedBlk := &snowman.TestBlock{
		TestDecidable: choices.TestDecidable{
			IDV: ids.ID{'i', 'n', 'n', 'e', 'r', 'S', 'y', 'n', 'c', 'e', 'd'},
		},
		HeightV: innerSummary.Height(),
		BytesV:  []byte("inner state synced block"),
	}
	stateSummary := createPostForkStateSummary(t, vm, forkHeight, proVMParentStateSummaryBlk, innerVM, innerSummary, innerStateSyncedBlk)

	ctx := context.Background()
	innerSummary.AcceptF = func(ctx context.Context) (block.StateSyncMode, error) {
		return block.StateSyncStatic, nil
	}
	_, err := stateSummary.Accept(ctx)
	require.NoError(err)

	// 2. Check that block backfilling is enabled looking at innerVM
	innerVM.BackfillBlocksEnabledF = func(ctx context.Context) (ids.ID, uint64, error) {
		return ids.Empty, 0, block.ErrBlockBackfillingNotEnabled
	}
	_, _, err = vm.BackfillBlocksEnabled(ctx)
	require.ErrorIs(err, block.ErrBlockBackfillingNotEnabled)

	innerVM.BackfillBlocksEnabledF = func(ctx context.Context) (ids.ID, uint64, error) {
		return innerStateSyncedBlk.ID(), innerStateSyncedBlk.Height() - 1, nil
	}
	blkID, _, err := vm.BackfillBlocksEnabled(ctx)
	require.NoError(err)
	require.Equal(proVMParentStateSummaryBlk, blkID)

	innerVM.BackfillBlocksEnabledF = func(ctx context.Context) (ids.ID, uint64, error) {
		return ids.Empty, 0, block.ErrBlockBackfillingNotEnabled
	}
	_, _, err = vm.BackfillBlocksEnabled(ctx)
	require.ErrorIs(err, block.ErrBlockBackfillingNotEnabled)
}

func TestBlockBackfillPostForkSuccess(t *testing.T) {
	// setup VM with backfill enabled
	require := require.New(t)
	toEngineCh := make(chan common.Message)
	innerVM, vm := setupBlockBackfillingVM(t, toEngineCh)
	defer func() {
		require.NoError(vm.Shutdown(context.Background()))
	}()

	var (
		forkHeight     = uint64(100)
		blkCount       = 12
		startBlkHeight = uint64(1492)

		// create a list of consecutive blocks and build state summary of top of them
		proBlks, innerBlks = createTestBlocks(t, vm, forkHeight, blkCount, startBlkHeight)

		innerTopBlk        = innerBlks[len(innerBlks)-1]
		proTopBlk          = proBlks[len(proBlks)-1]
		stateSummaryHeight = innerTopBlk.Height() + 1
	)

	innerSummary := &block.TestStateSummary{
		IDV:     ids.ID{'s', 'u', 'm', 'm', 'a', 'r', 'y', 'I', 'D'},
		HeightV: stateSummaryHeight,
		BytesV:  []byte{'i', 'n', 'n', 'e', 'r'},
	}
	innerStateSyncedBlk := &snowman.TestBlock{
		TestDecidable: choices.TestDecidable{
			IDV: ids.ID{'i', 'n', 'n', 'e', 'r', 'S', 'y', 'n', 'c', 'e', 'd'},
		},
		ParentV: innerTopBlk.ID(),
		HeightV: innerSummary.Height(),
		BytesV:  []byte("inner state synced block"),
	}
	stateSummary := createPostForkStateSummary(t, vm, forkHeight, proTopBlk.ID(), innerVM, innerSummary, innerStateSyncedBlk)

	innerSummary.AcceptF = func(ctx context.Context) (block.StateSyncMode, error) {
		return block.StateSyncStatic, nil
	}

	ctx := context.Background()
	_, err := stateSummary.Accept(ctx)
	require.NoError(err)

	innerVM.BackfillBlocksEnabledF = func(ctx context.Context) (ids.ID, uint64, error) {
		return innerStateSyncedBlk.ID(), innerStateSyncedBlk.Height() - 1, nil
	}

	blkID, _, err := vm.BackfillBlocksEnabled(ctx)
	require.NoError(err)
	require.Equal(proTopBlk.ID(), blkID)

	// Backfill some blocks
	innerVM.ParseBlockF = func(_ context.Context, b []byte) (snowman.Block, error) {
		for _, blk := range innerBlks {
			if bytes.Equal(b, blk.Bytes()) {
				return blk, nil
			}
		}
		return nil, database.ErrNotFound
	}
	innerVM.GetBlockF = func(_ context.Context, blkID ids.ID) (snowman.Block, error) {
		for _, blk := range innerBlks {
			if blkID == blk.ID() {
				return blk, nil
			}
		}
		return nil, database.ErrNotFound
	}
	innerVM.BackfillBlocksF = func(_ context.Context, b [][]byte) (ids.ID, uint64, error) {
		lowestblk := innerBlks[0]
		for _, blk := range innerBlks {
			if blk.Height() < lowestblk.Height() {
				lowestblk = blk
			}
		}
		return lowestblk.Parent(), lowestblk.Height() - 1, nil
	}

	blkBytes := make([][]byte, 0, len(proBlks))
	for _, blk := range proBlks {
		blkBytes = append(blkBytes, blk.Bytes())
	}
	nextBlkID, nextBlkHeight, err := vm.BackfillBlocks(ctx, blkBytes)
	require.NoError(err)
	require.Equal(proBlks[0].Parent(), nextBlkID)
	require.Equal(proBlks[0].Height()-1, nextBlkHeight)

	// check proBlocks have been indexed
	for _, blk := range proBlks {
		blkID, err := vm.GetBlockIDAtHeight(ctx, blk.Height())
		require.NoError(err)
		require.Equal(blk.ID(), blkID)

		_, err = vm.GetBlock(ctx, blkID)
		require.NoError(err)
	}
}

func TestBlockBackfillPostForkPartialSuccess(t *testing.T) {
	// setup VM with backfill enabled
	require := require.New(t)
	toEngineCh := make(chan common.Message)
	innerVM, vm := setupBlockBackfillingVM(t, toEngineCh)
	defer func() {
		require.NoError(vm.Shutdown(context.Background()))
	}()

	var (
		forkHeight     = uint64(100)
		blkCount       = 10
		startBlkHeight = uint64(1492)

		// create a list of consecutive blocks and build state summary of top of them
		proBlks, innerBlks = createTestBlocks(t, vm, forkHeight, blkCount, startBlkHeight)

		innerTopBlk        = innerBlks[len(innerBlks)-1]
		proTopBlk          = proBlks[len(proBlks)-1]
		stateSummaryHeight = innerTopBlk.Height() + 1
	)

	innerSummary := &block.TestStateSummary{
		IDV:     ids.ID{'s', 'u', 'm', 'm', 'a', 'r', 'y', 'I', 'D'},
		HeightV: stateSummaryHeight,
		BytesV:  []byte{'i', 'n', 'n', 'e', 'r'},
	}
	innerStateSyncedBlk := &snowman.TestBlock{
		TestDecidable: choices.TestDecidable{
			IDV: ids.ID{'i', 'n', 'n', 'e', 'r', 'S', 'y', 'n', 'c', 'e', 'd'},
		},
		ParentV: innerTopBlk.ID(),
		HeightV: innerSummary.Height(),
		BytesV:  []byte("inner state synced block"),
	}
	stateSummary := createPostForkStateSummary(t, vm, forkHeight, proTopBlk.ID(), innerVM, innerSummary, innerStateSyncedBlk)

	innerSummary.AcceptF = func(ctx context.Context) (block.StateSyncMode, error) {
		return block.StateSyncStatic, nil
	}

	ctx := context.Background()
	_, err := stateSummary.Accept(ctx)
	require.NoError(err)

	innerVM.BackfillBlocksEnabledF = func(ctx context.Context) (ids.ID, uint64, error) {
		return innerStateSyncedBlk.ID(), innerStateSyncedBlk.Height() - 1, nil
	}

	blkID, height, err := vm.BackfillBlocksEnabled(ctx)
	require.NoError(err)
	require.Equal(proTopBlk.ID(), blkID)
	require.Equal(proTopBlk.Height(), height)

	// Backfill some blocks
	innerVM.ParseBlockF = func(_ context.Context, b []byte) (snowman.Block, error) {
		for _, blk := range innerBlks {
			if bytes.Equal(b, blk.Bytes()) {
				return blk, nil
			}
		}
		return nil, database.ErrNotFound
	}

	// simulate that lower half of backfilled blocks won't be accepted by innerVM
	idx := len(innerBlks) / 2
	innerVM.GetBlockF = func(_ context.Context, blkID ids.ID) (snowman.Block, error) {
		for _, blk := range innerBlks {
			if blkID != blk.ID() {
				continue
			}
			// if it's one of the lower half blocks, assume it's not stored
			// since it was rejected
			if blk.Height() <= innerBlks[idx].Height() {
				return nil, database.ErrNotFound
			}
			return blk, nil
		}
		return nil, database.ErrNotFound
	}

	innerVM.BackfillBlocksF = func(_ context.Context, b [][]byte) (ids.ID, uint64, error) {
		// assume lowest half blocks fails verification
		return innerBlks[idx].ID(), innerBlks[idx].Height(), nil
	}

	blkBytes := make([][]byte, 0, len(proBlks))
	for _, blk := range proBlks {
		blkBytes = append(blkBytes, blk.Bytes())
	}
	nextBlkID, nextBlkHeight, err := vm.BackfillBlocks(ctx, blkBytes)
	require.NoError(err)
	require.Equal(proBlks[idx].ID(), nextBlkID)
	require.Equal(proBlks[idx].Height(), nextBlkHeight)

	// check only upper half of blocks have been indexed
	for i, blk := range proBlks {
		if i <= idx {
			_, err := vm.GetBlockIDAtHeight(ctx, blk.Height())
			require.ErrorIs(err, database.ErrNotFound)

			_, err = vm.GetBlock(ctx, blk.ID())
			require.ErrorIs(err, database.ErrNotFound)
		} else {
			blkID, err := vm.GetBlockIDAtHeight(ctx, blk.Height())
			require.NoError(err)
			require.Equal(blk.ID(), blkID)

			_, err = vm.GetBlock(ctx, blkID)
			require.NoError(err)
		}
	}
}

// Across Fork section
func TestBlockBackfillEnabledAcrossFork(t *testing.T) {
	require := require.New(t)
	toEngineCh := make(chan common.Message)
	innerVM, vm := setupBlockBackfillingVM(t, toEngineCh)
	defer func() {
		require.NoError(vm.Shutdown(context.Background()))
	}()

	// 1. Accept a State summary
	var (
		forkHeight                 = uint64(50)
		stateSummaryHeight         = forkHeight + 1
		proVMParentStateSummaryBlk = ids.GenerateTestID()
	)

	innerSummary := &block.TestStateSummary{
		IDV:     ids.ID{'s', 'u', 'm', 'm', 'a', 'r', 'y', 'I', 'D'},
		HeightV: stateSummaryHeight,
		BytesV:  []byte{'i', 'n', 'n', 'e', 'r'},
	}
	innerStateSyncedBlk := &snowman.TestBlock{
		TestDecidable: choices.TestDecidable{
			IDV: ids.ID{'i', 'n', 'n', 'e', 'r', 'S', 'y', 'n', 'c', 'e', 'd'},
		},
		HeightV: innerSummary.Height(),
		BytesV:  []byte("inner state synced block"),
	}
	stateSummary := createPostForkStateSummary(t, vm, forkHeight, proVMParentStateSummaryBlk, innerVM, innerSummary, innerStateSyncedBlk)

	innerSummary.AcceptF = func(ctx context.Context) (block.StateSyncMode, error) {
		return block.StateSyncStatic, nil
	}
	ctx := context.Background()
	_, err := stateSummary.Accept(ctx)
	require.NoError(err)

	// 2. Check that block backfilling is enabled looking at innerVM
	innerVM.BackfillBlocksEnabledF = func(ctx context.Context) (ids.ID, uint64, error) {
		return ids.Empty, 0, block.ErrBlockBackfillingNotEnabled
	}
	_, _, err = vm.BackfillBlocksEnabled(ctx)
	require.ErrorIs(err, block.ErrBlockBackfillingNotEnabled)

	innerVM.BackfillBlocksEnabledF = func(ctx context.Context) (ids.ID, uint64, error) {
		return innerStateSyncedBlk.ID(), innerStateSyncedBlk.Height() - 1, nil
	}
	blkID, _, err := vm.BackfillBlocksEnabled(ctx)
	require.NoError(err)
	require.Equal(proVMParentStateSummaryBlk, blkID)

	innerVM.BackfillBlocksEnabledF = func(ctx context.Context) (ids.ID, uint64, error) {
		return ids.Empty, 0, block.ErrBlockBackfillingNotEnabled
	}
	_, _, err = vm.BackfillBlocksEnabled(ctx)
	require.ErrorIs(err, block.ErrBlockBackfillingNotEnabled)
}

func TestBlockBackfillAcrossForkSuccess(t *testing.T) {
	// setup VM with backfill enabled
	require := require.New(t)
	toEngineCh := make(chan common.Message)
	innerVM, vm := setupBlockBackfillingVM(t, toEngineCh)
	defer func() {
		require.NoError(vm.Shutdown(context.Background()))
	}()

	var (
		forkHeight     = uint64(100)
		blkCount       = 4
		startBlkHeight = forkHeight - uint64(blkCount)/2

		// create a list of consecutive blocks and build state summary of top of them
		proBlks, innerBlks = createTestBlocks(t, vm, forkHeight, blkCount, startBlkHeight)

		innerTopBlk        = innerBlks[len(innerBlks)-1]
		proTopBlk          = proBlks[len(proBlks)-1]
		stateSummaryHeight = innerTopBlk.Height() + 1
	)

	innerSummary := &block.TestStateSummary{
		IDV:     ids.ID{'s', 'u', 'm', 'm', 'a', 'r', 'y', 'I', 'D'},
		HeightV: stateSummaryHeight,
		BytesV:  []byte{'i', 'n', 'n', 'e', 'r'},
	}
	innerStateSyncedBlk := &snowman.TestBlock{
		TestDecidable: choices.TestDecidable{
			IDV: ids.ID{'i', 'n', 'n', 'e', 'r', 'S', 'y', 'n', 'c', 'e', 'd'},
		},
		ParentV: innerTopBlk.ID(),
		HeightV: innerSummary.Height(),
		BytesV:  []byte("inner state synced block"),
	}
	stateSummary := createPostForkStateSummary(t, vm, forkHeight, proTopBlk.ID(), innerVM, innerSummary, innerStateSyncedBlk)

	innerSummary.AcceptF = func(ctx context.Context) (block.StateSyncMode, error) {
		return block.StateSyncStatic, nil
	}

	ctx := context.Background()
	_, err := stateSummary.Accept(ctx)
	require.NoError(err)

	innerVM.BackfillBlocksEnabledF = func(ctx context.Context) (ids.ID, uint64, error) {
		return innerStateSyncedBlk.ID(), innerStateSyncedBlk.Height() - 1, nil
	}

	blkID, _, err := vm.BackfillBlocksEnabled(ctx)
	require.NoError(err)
	require.Equal(proTopBlk.ID(), blkID)

	// Backfill some blocks
	innerVM.ParseBlockF = func(_ context.Context, b []byte) (snowman.Block, error) {
		for _, blk := range innerBlks {
			if bytes.Equal(b, blk.Bytes()) {
				return blk, nil
			}
		}
		return nil, database.ErrNotFound
	}
	innerVM.GetBlockF = func(_ context.Context, blkID ids.ID) (snowman.Block, error) {
		for _, blk := range innerBlks {
			if blkID == blk.ID() {
				return blk, nil
			}
		}
		return nil, database.ErrNotFound
	}
	innerVM.BackfillBlocksF = func(_ context.Context, b [][]byte) (ids.ID, uint64, error) {
		lowestblk := innerBlks[0]
		for _, blk := range innerBlks {
			if blk.Height() < lowestblk.Height() {
				lowestblk = blk
			}
		}
		return lowestblk.Parent(), lowestblk.Height() - 1, nil
	}
	innerVM.GetBlockIDAtHeightF = func(ctx context.Context, height uint64) (ids.ID, error) {
		for _, blk := range innerBlks {
			if height == blk.Height() {
				return blk.ID(), nil
			}
		}
		return ids.Empty, database.ErrNotFound
	}

	blkBytes := make([][]byte, 0, len(proBlks))
	for _, blk := range proBlks {
		blkBytes = append(blkBytes, blk.Bytes())
	}
	nextBlkID, nextBlkHeight, err := vm.BackfillBlocks(ctx, blkBytes)
	require.NoError(err)
	require.Equal(proBlks[0].Parent(), nextBlkID)
	require.Equal(proBlks[0].Height()-1, nextBlkHeight)

	// check proBlocks have been indexed
	for _, blk := range proBlks {
		blkID, err := vm.GetBlockIDAtHeight(ctx, blk.Height())
		require.NoError(err)
		require.Equal(blk.ID(), blkID)

		_, err = vm.GetBlock(ctx, blkID)
		require.NoError(err)
	}
}

func TestBlockBackfillAcrossForkPartialSuccess(t *testing.T) {
	// setup VM with backfill enabled
	require := require.New(t)
	toEngineCh := make(chan common.Message)
	innerVM, vm := setupBlockBackfillingVM(t, toEngineCh)
	defer func() {
		require.NoError(vm.Shutdown(context.Background()))
	}()

	var (
		forkHeight     = uint64(100)
		blkCount       = 8
		startBlkHeight = forkHeight - uint64(blkCount)/2

		// simulate that the bottom [idxFailure] blocks will fail
		// being pushed in innerVM
		idxFailure = 3

		// create a list of consecutive blocks and build state summary of top of them
		proBlks, innerBlks = createTestBlocks(t, vm, forkHeight, blkCount, startBlkHeight)

		innerTopBlk        = innerBlks[len(innerBlks)-1]
		proTopBlk          = proBlks[len(proBlks)-1]
		stateSummaryHeight = innerTopBlk.Height() + 1
	)

	innerSummary := &block.TestStateSummary{
		IDV:     ids.ID{'s', 'u', 'm', 'm', 'a', 'r', 'y', 'I', 'D'},
		HeightV: stateSummaryHeight,
		BytesV:  []byte{'i', 'n', 'n', 'e', 'r'},
	}
	innerStateSyncedBlk := &snowman.TestBlock{
		TestDecidable: choices.TestDecidable{
			IDV: ids.ID{'i', 'n', 'n', 'e', 'r', 'S', 'y', 'n', 'c', 'e', 'd'},
		},
		ParentV: innerTopBlk.ID(),
		HeightV: innerSummary.Height(),
		BytesV:  []byte("inner state synced block"),
	}
	stateSummary := createPostForkStateSummary(t, vm, forkHeight, proTopBlk.ID(), innerVM, innerSummary, innerStateSyncedBlk)

	innerSummary.AcceptF = func(ctx context.Context) (block.StateSyncMode, error) {
		return block.StateSyncStatic, nil
	}

	ctx := context.Background()
	_, err := stateSummary.Accept(ctx)
	require.NoError(err)

	innerVM.BackfillBlocksEnabledF = func(ctx context.Context) (ids.ID, uint64, error) {
		return innerStateSyncedBlk.ID(), innerStateSyncedBlk.Height() - 1, nil
	}

	blkID, height, err := vm.BackfillBlocksEnabled(ctx)
	require.NoError(err)
	require.Equal(proTopBlk.ID(), blkID)
	require.Equal(proTopBlk.Height(), height)

	// Backfill some blocks
	innerVM.ParseBlockF = func(_ context.Context, b []byte) (snowman.Block, error) {
		for _, blk := range innerBlks {
			if bytes.Equal(b, blk.Bytes()) {
				return blk, nil
			}
		}
		return nil, database.ErrNotFound
	}

	// simulate that lower half of backfilled blocks won't be accepted by innerVM
	innerVM.GetBlockF = func(_ context.Context, blkID ids.ID) (snowman.Block, error) {
		for _, blk := range innerBlks {
			if blkID != blk.ID() {
				continue
			}
			// if it's one of the lower half blocks, assume it's not stored
			// since it was rejected
			if blk.Height() <= innerBlks[idxFailure].Height() {
				return nil, database.ErrNotFound
			}
			return blk, nil
		}
		return nil, database.ErrNotFound
	}
	innerVM.GetBlockIDAtHeightF = func(ctx context.Context, height uint64) (ids.ID, error) {
		for _, blk := range innerBlks {
			if height != blk.Height() {
				continue
			}
			// if it's one of the lower half blocks, assume it's not stored
			// since it was rejected
			if blk.Height() <= innerBlks[idxFailure].Height() {
				return ids.Empty, database.ErrNotFound
			}
			return blk.ID(), nil
		}
		return ids.Empty, database.ErrNotFound
	}

	innerVM.BackfillBlocksF = func(_ context.Context, b [][]byte) (ids.ID, uint64, error) {
		// assume lowest half blocks fails verification
		return innerBlks[idxFailure].ID(), innerBlks[idxFailure].Height(), nil
	}

	blkBytes := make([][]byte, 0, len(proBlks))
	for _, blk := range proBlks {
		blkBytes = append(blkBytes, blk.Bytes())
	}
	nextBlkID, nextBlkHeight, err := vm.BackfillBlocks(ctx, blkBytes)
	require.NoError(err)
	require.Equal(proBlks[idxFailure].ID(), nextBlkID)
	require.Equal(proBlks[idxFailure].Height(), nextBlkHeight)

	// check only upper half of blocks have been indexed
	for i, blk := range proBlks {
		if i <= idxFailure {
			_, err := vm.GetBlockIDAtHeight(ctx, blk.Height())
			require.ErrorIs(err, database.ErrNotFound)

			_, err = vm.GetBlock(ctx, blk.ID())
			require.ErrorIs(err, database.ErrNotFound)
		} else {
			blkID, err := vm.GetBlockIDAtHeight(ctx, blk.Height())
			require.NoError(err)
			require.Equal(blk.ID(), blkID)

			_, err = vm.GetBlock(ctx, blkID)
			require.NoError(err)
		}
	}
}

// Pre Fork section
func TestBlockBackfillEnabledPreFork(t *testing.T) {
	require := require.New(t)
	toEngineCh := make(chan common.Message)
	innerVM, vm := setupBlockBackfillingVM(t, toEngineCh)
	defer func() {
		require.NoError(vm.Shutdown(context.Background()))
	}()

	// 1. Accept a State summary
	stateSummary := &block.TestStateSummary{
		IDV:     ids.ID{'s', 'u', 'm', 'm', 'a', 'r', 'y', 'I', 'D'},
		HeightV: 100,
		BytesV:  []byte{'i', 'n', 'n', 'e', 'r'},
	}
	innerStateSyncedBlk := &snowman.TestBlock{
		TestDecidable: choices.TestDecidable{
			IDV: ids.ID{'i', 'n', 'n', 'e', 'r', 'S', 'y', 'n', 'c', 'e', 'd'},
		},
		HeightV: stateSummary.Height(),
		BytesV:  []byte("inner state synced block"),
	}

	stateSummary.AcceptF = func(ctx context.Context) (block.StateSyncMode, error) {
		return block.StateSyncStatic, nil
	}

	ctx := context.Background()
	_, err := stateSummary.Accept(ctx)
	require.NoError(err)

	// 2. Check that block backfilling is enabled looking at innerVM
	innerVM.BackfillBlocksEnabledF = func(ctx context.Context) (ids.ID, uint64, error) {
		return ids.Empty, 0, block.ErrBlockBackfillingNotEnabled
	}
	_, _, err = vm.BackfillBlocksEnabled(ctx)
	require.ErrorIs(err, block.ErrBlockBackfillingNotEnabled)

	innerVM.BackfillBlocksEnabledF = func(_ context.Context) (ids.ID, uint64, error) {
		return innerStateSyncedBlk.ID(), innerStateSyncedBlk.Height() - 1, nil
	}
	innerVM.GetBlockIDAtHeightF = func(ctx context.Context, height uint64) (ids.ID, error) {
		if height == innerStateSyncedBlk.Height() {
			return innerStateSyncedBlk.ID(), nil
		}
		return ids.Empty, database.ErrNotFound
	}
	blkID, _, err := vm.BackfillBlocksEnabled(ctx)
	require.NoError(err)
	require.Equal(innerStateSyncedBlk.Parent(), blkID)

	innerVM.BackfillBlocksEnabledF = func(ctx context.Context) (ids.ID, uint64, error) {
		return ids.Empty, 0, block.ErrBlockBackfillingNotEnabled
	}
	_, _, err = vm.BackfillBlocksEnabled(ctx)
	require.ErrorIs(err, block.ErrBlockBackfillingNotEnabled)
}

func TestBlockBackfillPreForkSuccess(t *testing.T) {
	// setup VM with backfill enabled
	require := require.New(t)
	toEngineCh := make(chan common.Message)
	innerVM, vm := setupBlockBackfillingVM(t, toEngineCh)
	defer func() {
		require.NoError(vm.Shutdown(context.Background()))
	}()

	var (
		forkHeight     = uint64(2000)
		blkCount       = 8
		startBlkHeight = uint64(100)

		// create a list of consecutive blocks and build state summary of top of them
		// proBlks should all be preForkBlocks
		proBlks, innerBlks = createTestBlocks(t, vm, forkHeight, blkCount, startBlkHeight)

		innerTopBlk        = innerBlks[len(innerBlks)-1]
		preForkTopBlk      = proBlks[len(proBlks)-1]
		stateSummaryHeight = innerTopBlk.Height() + 1
	)

	stateSummary := &block.TestStateSummary{
		IDV:     ids.ID{'s', 'u', 'm', 'm', 'a', 'r', 'y', 'I', 'D'},
		HeightV: stateSummaryHeight,
		BytesV:  []byte{'i', 'n', 'n', 'e', 'r'},
	}
	innerStateSyncedBlk := &snowman.TestBlock{
		TestDecidable: choices.TestDecidable{
			IDV: ids.ID{'i', 'n', 'n', 'e', 'r', 'S', 'y', 'n', 'c', 'e', 'd'},
		},
		ParentV: innerTopBlk.ID(),
		HeightV: stateSummary.Height(),
		BytesV:  []byte("inner state synced block"),
	}
	stateSummary.AcceptF = func(ctx context.Context) (block.StateSyncMode, error) {
		return block.StateSyncStatic, nil
	}

	ctx := context.Background()
	_, err := stateSummary.Accept(ctx)
	require.NoError(err)

	innerVM.LastAcceptedF = func(context.Context) (ids.ID, error) {
		return innerStateSyncedBlk.ID(), nil
	}
	innerVM.GetBlockF = func(_ context.Context, blkID ids.ID) (snowman.Block, error) {
		switch blkID {
		case innerStateSyncedBlk.ID():
			return innerStateSyncedBlk, nil
		default:
			return nil, database.ErrNotFound
		}
	}

	innerVM.BackfillBlocksEnabledF = func(ctx context.Context) (ids.ID, uint64, error) {
		return innerStateSyncedBlk.ID(), innerStateSyncedBlk.Height() - 1, nil
	}
	innerVM.GetBlockIDAtHeightF = func(ctx context.Context, height uint64) (ids.ID, error) {
		if height == innerStateSyncedBlk.Height() {
			return innerStateSyncedBlk.ID(), nil
		}
		return ids.Empty, database.ErrNotFound
	}

	blkID, _, err := vm.BackfillBlocksEnabled(ctx)
	require.NoError(err)
	require.Equal(preForkTopBlk.ID(), blkID)

	// Backfill some blocks
	innerVM.ParseBlockF = func(_ context.Context, b []byte) (snowman.Block, error) {
		for _, blk := range innerBlks {
			if bytes.Equal(b, blk.Bytes()) {
				return blk, nil
			}
		}
		return nil, database.ErrNotFound
	}
	innerVM.GetBlockF = func(_ context.Context, blkID ids.ID) (snowman.Block, error) {
		for _, blk := range innerBlks {
			if blkID == blk.ID() {
				return blk, nil
			}
		}
		return nil, database.ErrNotFound
	}
	innerVM.GetBlockIDAtHeightF = func(ctx context.Context, height uint64) (ids.ID, error) {
		for _, blk := range innerBlks {
			if height == blk.Height() {
				return blk.ID(), nil
			}
		}
		return ids.Empty, database.ErrNotFound
	}
	innerVM.BackfillBlocksF = func(_ context.Context, b [][]byte) (ids.ID, uint64, error) {
		lowestblk := innerBlks[0]
		for _, blk := range innerBlks {
			if blk.Height() < lowestblk.Height() {
				lowestblk = blk
			}
		}
		return lowestblk.Parent(), lowestblk.Height() - 1, nil
	}

	blkBytes := make([][]byte, 0, len(proBlks))
	for _, blk := range proBlks {
		blkBytes = append(blkBytes, blk.Bytes())
	}
	nextBlkID, nextBlkHeight, err := vm.BackfillBlocks(ctx, blkBytes)
	require.NoError(err)
	require.Equal(proBlks[0].Parent(), nextBlkID)
	require.Equal(proBlks[0].Height()-1, nextBlkHeight)

	// check proBlocks have been indexed
	for _, blk := range proBlks {
		blkID, err := vm.GetBlockIDAtHeight(ctx, blk.Height())
		require.NoError(err)
		require.Equal(blk.ID(), blkID)

		_, err = vm.GetBlock(ctx, blkID)
		require.NoError(err)
	}
}

func TestBlockBackfillPreForkPartialSuccess(t *testing.T) {
	// setup VM with backfill enabled
	require := require.New(t)
	toEngineCh := make(chan common.Message)
	innerVM, vm := setupBlockBackfillingVM(t, toEngineCh)
	defer func() {
		require.NoError(vm.Shutdown(context.Background()))
	}()

	var (
		forkHeight     = uint64(2000)
		blkCount       = 10
		startBlkHeight = uint64(100)

		// create a list of consecutive blocks and build state summary of top of them
		// proBlks should all be preForkBlocks
		proBlks, innerBlks = createTestBlocks(t, vm, forkHeight, blkCount, startBlkHeight)

		innerTopBlk        = innerBlks[len(innerBlks)-1]
		preForkTopBlk      = proBlks[len(proBlks)-1]
		stateSummaryHeight = innerTopBlk.Height() + 1
	)

	stateSummary := &block.TestStateSummary{
		IDV:     ids.ID{'s', 'u', 'm', 'm', 'a', 'r', 'y', 'I', 'D'},
		HeightV: stateSummaryHeight,
		BytesV:  []byte{'i', 'n', 'n', 'e', 'r'},
	}
	innerStateSyncedBlk := &snowman.TestBlock{
		TestDecidable: choices.TestDecidable{
			IDV: ids.ID{'i', 'n', 'n', 'e', 'r', 'S', 'y', 'n', 'c', 'e', 'd'},
		},
		ParentV: innerTopBlk.ID(),
		HeightV: stateSummary.Height(),
		BytesV:  []byte("inner state synced block"),
	}
	stateSummary.AcceptF = func(ctx context.Context) (block.StateSyncMode, error) {
		return block.StateSyncStatic, nil
	}

	ctx := context.Background()
	_, err := stateSummary.Accept(ctx)
	require.NoError(err)

	innerVM.LastAcceptedF = func(context.Context) (ids.ID, error) {
		return innerStateSyncedBlk.ID(), nil
	}
	innerVM.GetBlockF = func(_ context.Context, blkID ids.ID) (snowman.Block, error) {
		switch blkID {
		case innerStateSyncedBlk.ID():
			return innerStateSyncedBlk, nil
		default:
			return nil, database.ErrNotFound
		}
	}

	innerVM.BackfillBlocksEnabledF = func(ctx context.Context) (ids.ID, uint64, error) {
		return innerStateSyncedBlk.ID(), innerStateSyncedBlk.Height() - 1, nil
	}
	innerVM.GetBlockIDAtHeightF = func(ctx context.Context, height uint64) (ids.ID, error) {
		if height == innerStateSyncedBlk.Height() {
			return innerStateSyncedBlk.ID(), nil
		}
		return ids.Empty, database.ErrNotFound
	}

	blkID, _, err := vm.BackfillBlocksEnabled(ctx)
	require.NoError(err)
	require.Equal(preForkTopBlk.ID(), blkID)

	// Backfill some blocks
	innerVM.ParseBlockF = func(_ context.Context, b []byte) (snowman.Block, error) {
		for _, blk := range innerBlks {
			if bytes.Equal(b, blk.Bytes()) {
				return blk, nil
			}
		}
		return nil, database.ErrNotFound
	}
	// simulate that lower half of backfilled blocks won't be accepted by innerVM
	idx := len(innerBlks) / 2
	innerVM.GetBlockF = func(_ context.Context, blkID ids.ID) (snowman.Block, error) {
		for _, blk := range innerBlks {
			if blkID != blk.ID() {
				continue
			}
			// if it's one of the lower half blocks, assume it's not stored
			// since it was rejected
			if blk.Height() <= innerBlks[idx].Height() {
				return nil, database.ErrNotFound
			}
			return blk, nil
		}
		return nil, database.ErrNotFound
	}
	innerVM.GetBlockIDAtHeightF = func(ctx context.Context, height uint64) (ids.ID, error) {
		for _, blk := range innerBlks {
			if height != blk.Height() {
				continue
			}
			// if it's one of the lower half blocks, assume it's not stored
			// since it was rejected
			if blk.Height() <= innerBlks[idx].Height() {
				return ids.Empty, database.ErrNotFound
			}
			return blk.ID(), nil
		}
		return ids.Empty, database.ErrNotFound
	}
	innerVM.BackfillBlocksF = func(_ context.Context, b [][]byte) (ids.ID, uint64, error) {
		// assume lowest half blocks fails verification
		return innerBlks[idx].ID(), innerBlks[idx].Height(), nil
	}

	blkBytes := make([][]byte, 0, len(proBlks))
	for _, blk := range proBlks {
		blkBytes = append(blkBytes, blk.Bytes())
	}
	nextBlkID, nextBlkHeight, err := vm.BackfillBlocks(ctx, blkBytes)
	require.NoError(err)
	require.Equal(proBlks[idx].ID(), nextBlkID)
	require.Equal(proBlks[idx].Height(), nextBlkHeight)

	// check only upper half of blocks have been indexed
	for i, blk := range proBlks {
		if i <= idx {
			_, err := vm.GetBlockIDAtHeight(ctx, blk.Height())
			require.ErrorIs(err, database.ErrNotFound)

			_, err = vm.GetBlock(ctx, blk.ID())
			require.ErrorIs(err, database.ErrNotFound)
		} else {
			blkID, err := vm.GetBlockIDAtHeight(ctx, blk.Height())
			require.NoError(err)
			require.Equal(blk.ID(), blkID)

			_, err = vm.GetBlock(ctx, blkID)
			require.NoError(err)
		}
	}
}

func createTestBlocks(
	t *testing.T,
	vm *VM,
	forkHeight uint64,
	blkCount int,
	startBlkHeight uint64,
) (
	[]snowman.Block, // proposerVM blocks
	[]snowman.Block, // inner VM blocks
) {
	require := require.New(t)
	var (
		latestInnerBlkID = ids.GenerateTestID()

		dummyBlkTime      = time.Now()
		dummyPChainHeight = startBlkHeight / 2

		innerBlks = make([]snowman.Block, 0, blkCount)
		proBlks   = make([]snowman.Block, 0, blkCount)
	)
	for idx := 0; idx < blkCount; idx++ {
		blkHeight := startBlkHeight + uint64(idx)

		rndBytes := ids.GenerateTestID()
		innerBlkTop := &snowman.TestBlock{
			TestDecidable: choices.TestDecidable{
				IDV:     ids.GenerateTestID(),
				StatusV: choices.Processing,
			},
			BytesV:  rndBytes[:],
			ParentV: latestInnerBlkID,
			HeightV: startBlkHeight + uint64(idx),
		}
		latestInnerBlkID = innerBlkTop.ID()
		innerBlks = append(innerBlks, innerBlkTop)

		if blkHeight < forkHeight {
			proBlks = append(proBlks, &preForkBlock{
				vm:    vm,
				Block: innerBlkTop,
			})
		} else {
			latestProBlkID := ids.GenerateTestID()
			if len(proBlks) != 0 {
				latestProBlkID = proBlks[len(proBlks)-1].ID()
			}
			statelessChild, err := statelessblock.BuildUnsigned(
				latestProBlkID,
				dummyBlkTime,
				dummyPChainHeight,
				innerBlkTop.Bytes(),
			)
			require.NoError(err)
			proBlkTop := &postForkBlock{
				SignedBlock: statelessChild,
				postForkCommonComponents: postForkCommonComponents{
					vm:       vm,
					innerBlk: innerBlkTop,
					status:   choices.Processing,
				},
			}
			proBlks = append(proBlks, proBlkTop)
		}
	}
	return proBlks, innerBlks
}

func createPostForkStateSummary(
	t *testing.T,
	vm *VM,
	forkHeight uint64,
	proVMParentStateSummaryBlk ids.ID,
	innerVM *fullVM,
	innerSummary *block.TestStateSummary,
	innerBlk *snowman.TestBlock,
) block.StateSummary {
	require := require.New(t)

	pchainHeight := innerBlk.Height() / 2
	slb, err := statelessblock.Build(
		proVMParentStateSummaryBlk,
		innerBlk.Timestamp(),
		pchainHeight,
		vm.stakingCertLeaf,
		innerBlk.Bytes(),
		vm.ctx.ChainID,
		vm.stakingLeafSigner,
	)
	require.NoError(err)

	statelessSummary, err := summary.Build(forkHeight, slb.Bytes(), innerSummary.Bytes())
	require.NoError(err)

	innerVM.ParseStateSummaryF = func(ctx context.Context, summaryBytes []byte) (block.StateSummary, error) {
		require.Equal(innerSummary.BytesV, summaryBytes)
		return innerSummary, nil
	}
	innerVM.ParseBlockF = func(_ context.Context, b []byte) (snowman.Block, error) {
		require.Equal(innerBlk.Bytes(), b)
		return innerBlk, nil
	}

	summary, err := vm.ParseStateSummary(context.Background(), statelessSummary.Bytes())
	require.NoError(err)
	return summary
}

func setupBlockBackfillingVM(
	t *testing.T,
	toEngineCh chan<- common.Message,
) (
	*fullVM,
	*VM,
) {
	require := require.New(t)

	innerVM := &fullVM{
		TestVM: &block.TestVM{
			TestVM: common.TestVM{
				T: t,
			},
		},
		TestStateSyncableVM: &block.TestStateSyncableVM{
			T: t,
		},
	}

	// signal height index is complete
	innerVM.VerifyHeightIndexF = func(context.Context) error {
		return nil
	}

	// load innerVM expectations
	innerGenesisBlk := &snowman.TestBlock{
		TestDecidable: choices.TestDecidable{
			IDV: ids.ID{'i', 'n', 'n', 'e', 'r', 'G', 'e', 'n', 'e', 's', 'i', 's', 'I', 'D'},
		},
		HeightV: 0,
		BytesV:  []byte("genesis state"),
	}

	innerVM.InitializeF = func(_ context.Context, _ *snow.Context, _ database.Database,
		_ []byte, _ []byte, _ []byte, ch chan<- common.Message,
		_ []*common.Fx, _ common.AppSender,
	) error {
		return nil
	}
	innerVM.VerifyHeightIndexF = func(context.Context) error {
		return nil
	}
	innerVM.LastAcceptedF = func(context.Context) (ids.ID, error) {
		return innerGenesisBlk.ID(), nil
	}
	innerVM.GetBlockF = func(context.Context, ids.ID) (snowman.Block, error) {
		return innerGenesisBlk, nil
	}

	// createVM
	vm := New(
		innerVM,
		time.Time{},
		0,
		DefaultMinBlockDelay,
		DefaultNumHistoricalBlocks,
		pTestSigner,
		pTestCert,
	)

	ctx := snow.DefaultContextTest()
	ctx.NodeID = ids.NodeIDFromCert(pTestCert)

	require.NoError(vm.Initialize(
		context.Background(),
		ctx,
		memdb.New(),
		innerGenesisBlk.Bytes(),
		nil,
		nil,
		toEngineCh,
		nil,
		nil,
	))

	return innerVM, vm
}
