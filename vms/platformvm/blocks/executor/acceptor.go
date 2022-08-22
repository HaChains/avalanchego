// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package executor

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/ava-labs/avalanchego/snow/choices"
	"github.com/ava-labs/avalanchego/utils/window"
	"github.com/ava-labs/avalanchego/vms/platformvm/blocks"
	"github.com/ava-labs/avalanchego/vms/platformvm/metrics"
)

var _ blocks.Visitor = &acceptor{}

// acceptor handles the logic for accepting a block.
// All errors returned by this struct are fatal and should result in the chain
// being shutdown.
type acceptor struct {
	*backend
	metrics          metrics.Metrics
	recentlyAccepted *window.Window
}

func (a *acceptor) ApricotAbortBlock(b *blocks.ApricotAbortBlock) error {
	blkID := b.ID()
	parentID := b.Parent()
	a.ctx.Log.Verbo(
		"accepting block",
		zap.String("blockType", "apricot abort"),
		zap.Stringer("blkID", blkID),
		zap.Uint64("height", b.Height()),
		zap.Stringer("parentID", parentID),
	)

	parentState, ok := a.blkIDToState[parentID]
	if !ok {
		return fmt.Errorf("couldn't find state of block %s, parent of %s", parentID, blkID)
	}

	// Update metrics
	if a.bootstrapped.GetValue() {
		if parentState.initiallyPreferCommit {
			a.metrics.MarkOptionVoteLost()
		} else {
			a.metrics.MarkOptionVoteWon()
		}
	}

	return a.optionBlock(b)
}

func (a *acceptor) ApricotCommitBlock(b *blocks.ApricotCommitBlock) error {
	blkID := b.ID()
	parentID := b.Parent()
	a.ctx.Log.Verbo(
		"accepting block",
		zap.String("blockType", "apricot commit"),
		zap.Stringer("blkID", blkID),
		zap.Uint64("height", b.Height()),
		zap.Stringer("parentID", parentID),
	)

	parentState, ok := a.blkIDToState[parentID]
	if !ok {
		return fmt.Errorf("couldn't find state of block %s, parent of %s", parentID, blkID)
	}

	// Update metrics
	if a.bootstrapped.GetValue() {
		if parentState.initiallyPreferCommit {
			a.metrics.MarkOptionVoteWon()
		} else {
			a.metrics.MarkOptionVoteLost()
		}
	}

	return a.optionBlock(b)
}

// Note that:
//
// * We don't free the proposal block in this method.
//   It is freed when its child is accepted.
//   We need to keep this block's state in memory for its child to use.
//
// * We only update the metrics to reflect this block's
//   acceptance when its child is accepted.
//
// * We don't write this block to state here.
//   That is done when this block's child (a CommitBlock or AbortBlock) is accepted.
//   We do this so that in the event that the node shuts down, the proposal block
//   is not written to disk unless its child is.
//   (The VM's Shutdown method commits the database.)
//   The snowman.Engine requires that the last committed block is a decision block
func (a *acceptor) ApricotProposalBlock(b *blocks.ApricotProposalBlock) error {
	blkID := b.ID()
	a.ctx.Log.Verbo(
		"accepting block",
		zap.String("blockType", "apricot proposal"),
		zap.Stringer("blkID", blkID),
		zap.Uint64("height", b.Height()),
		zap.Stringer("parentID", b.Parent()),
	)

	// See comment for [lastAccepted].
	a.backend.lastAccepted = b.ID()
	return nil
}

func (a *acceptor) ApricotStandardBlock(b *blocks.ApricotStandardBlock) error {
	blkID := b.ID()
	defer a.free(blkID)

	a.ctx.Log.Verbo(
		"accepting block",
		zap.String("blockType", "apricot standard"),
		zap.Stringer("blkID", blkID),
		zap.Uint64("height", b.Height()),
		zap.Stringer("parentID", b.Parent()),
	)

	if err := a.commonAccept(b); err != nil {
		return err
	}

	blkState, ok := a.blkIDToState[blkID]
	if !ok {
		return fmt.Errorf("couldn't find state of block %s", blkID)
	}

	// Update the state to reflect the changes made in [onAcceptState].
	blkState.onAcceptState.Apply(a.state)

	defer a.state.Abort()
	batch, err := a.state.CommitBatch()
	if err != nil {
		return fmt.Errorf(
			"failed to commit VM's database for block %s: %w",
			blkID,
			err,
		)
	}

	// Note that this method writes [batch] to the database.
	if err := a.ctx.SharedMemory.Apply(blkState.atomicRequests, batch); err != nil {
		return fmt.Errorf("failed to apply vm's state to shared memory: %w", err)
	}

	if onAcceptFunc := blkState.onAcceptFunc; onAcceptFunc != nil {
		onAcceptFunc()
	}
	return nil
}

func (a *acceptor) ApricotAtomicBlock(b *blocks.ApricotAtomicBlock) error {
	blkID := b.ID()
	defer a.free(blkID)

	a.ctx.Log.Verbo(
		"accepting block",
		zap.String("blockType", "apricot atomic"),
		zap.Stringer("blkID", blkID),
		zap.Uint64("height", b.Height()),
		zap.Stringer("parentID", b.Parent()),
	)

	if err := a.commonAccept(b); err != nil {
		return err
	}

	blkState, ok := a.blkIDToState[blkID]
	if !ok {
		return fmt.Errorf("couldn't find state of block %s", blkID)
	}

	// Update the state to reflect the changes made in [onAcceptState].
	blkState.onAcceptState.Apply(a.state)

	defer a.state.Abort()
	batch, err := a.state.CommitBatch()
	if err != nil {
		return fmt.Errorf(
			"failed to commit VM's database for block %s: %w",
			blkID,
			err,
		)
	}

	// Note that this method writes [batch] to the database.
	if err := a.ctx.SharedMemory.Apply(blkState.atomicRequests, batch); err != nil {
		return fmt.Errorf(
			"failed to atomically accept tx %s in block %s: %w",
			b.Tx.ID(),
			blkID,
			err,
		)
	}
	return nil
}

func (a *acceptor) optionBlock(b blocks.Block) error {
	blkID := b.ID()
	parentID := b.Parent()

	defer func() {
		// Note: we assume this block's sibling doesn't
		// need the parent's state when it's rejected.
		a.free(parentID)
		a.free(blkID)
	}()

	// Note that the parent must be accepted first.
	parentState, ok := a.blkIDToState[parentID]
	if !ok {
		return fmt.Errorf("couldn't find state of block %s, parent of %s", parentID, blkID)
	}
	if err := a.commonAccept(parentState.statelessBlock); err != nil {
		return err
	}

	if err := a.commonAccept(b); err != nil {
		return err
	}

	blkState, ok := a.blkIDToState[blkID]
	if !ok {
		return fmt.Errorf("couldn't find state of block %s", blkID)
	}
	blkState.onAcceptState.Apply(a.state)
	return a.state.Commit()
}

func (a *acceptor) commonAccept(b blocks.Block) error {
	blkID := b.ID()

	if err := a.metrics.MarkAccepted(b); err != nil {
		return fmt.Errorf("failed to accept block %s: %w", blkID, err)
	}

	a.backend.lastAccepted = blkID
	a.state.SetLastAccepted(blkID)
	a.state.SetHeight(b.Height())
	a.state.AddStatelessBlock(b, choices.Accepted)
	a.recentlyAccepted.Add(blkID)
	return nil
}
