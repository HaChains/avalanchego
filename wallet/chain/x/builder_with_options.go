// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package x

import (
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/avm/txs"
	"github.com/ava-labs/avalanchego/vms/avm/txs/fees"
	"github.com/ava-labs/avalanchego/vms/components/avax"
	"github.com/ava-labs/avalanchego/vms/components/verify"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/chain/x/backends"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary/common"
)

var _ backends.Builder = (*builderWithOptions)(nil)

type builderWithOptions struct {
	backends.Builder
	options []common.Option
}

// NewBuilderWithOptions returns a new transaction builder that will use the
// given options by default.
//
//   - [builder] is the builder that will be called to perform the underlying
//     operations.
//   - [options] will be provided to the builder in addition to the options
//     provided in the method calls.
func NewBuilderWithOptions(builder backends.Builder, options ...common.Option) backends.Builder {
	return &builderWithOptions{
		Builder: builder,
		options: options,
	}
}

func (b *builderWithOptions) GetFTBalance(
	options ...common.Option,
) (map[ids.ID]uint64, error) {
	return b.Builder.GetFTBalance(
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) GetImportableBalance(
	chainID ids.ID,
	options ...common.Option,
) (map[ids.ID]uint64, error) {
	return b.Builder.GetImportableBalance(
		chainID,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewBaseTx(
	outputs []*avax.TransferableOutput,
	feeCalc *fees.Calculator,
	options ...common.Option,
) (*txs.BaseTx, error) {
	return b.Builder.NewBaseTx(
		outputs,
		feeCalc,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewCreateAssetTx(
	name string,
	symbol string,
	denomination byte,
	initialState map[uint32][]verify.State,
	feeCalc *fees.Calculator,
	options ...common.Option,
) (*txs.CreateAssetTx, error) {
	return b.Builder.NewCreateAssetTx(
		name,
		symbol,
		denomination,
		initialState,
		feeCalc,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewOperationTx(
	operations []*txs.Operation,
	feeCalc *fees.Calculator,
	options ...common.Option,
) (*txs.OperationTx, error) {
	return b.Builder.NewOperationTx(
		operations,
		feeCalc,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewOperationTxMintFT(
	outputs map[ids.ID]*secp256k1fx.TransferOutput,
	feeCalc *fees.Calculator,
	options ...common.Option,
) (*txs.OperationTx, error) {
	return b.Builder.NewOperationTxMintFT(
		outputs,
		feeCalc,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewOperationTxMintNFT(
	assetID ids.ID,
	payload []byte,
	owners []*secp256k1fx.OutputOwners,
	feeCalc *fees.Calculator,
	options ...common.Option,
) (*txs.OperationTx, error) {
	return b.Builder.NewOperationTxMintNFT(
		assetID,
		payload,
		owners,
		feeCalc,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewOperationTxMintProperty(
	assetID ids.ID,
	owner *secp256k1fx.OutputOwners,
	feeCalc *fees.Calculator,
	options ...common.Option,
) (*txs.OperationTx, error) {
	return b.Builder.NewOperationTxMintProperty(
		assetID,
		owner,
		feeCalc,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewOperationTxBurnProperty(
	assetID ids.ID,
	feeCalc *fees.Calculator,
	options ...common.Option,
) (*txs.OperationTx, error) {
	return b.Builder.NewOperationTxBurnProperty(
		assetID,
		feeCalc,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewImportTx(
	chainID ids.ID,
	to *secp256k1fx.OutputOwners,
	feeCalc *fees.Calculator,
	options ...common.Option,
) (*txs.ImportTx, error) {
	return b.Builder.NewImportTx(
		chainID,
		to,
		feeCalc,
		common.UnionOptions(b.options, options)...,
	)
}

func (b *builderWithOptions) NewExportTx(
	chainID ids.ID,
	outputs []*avax.TransferableOutput,
	feeCalc *fees.Calculator,
	options ...common.Option,
) (*txs.ExportTx, error) {
	return b.Builder.NewExportTx(
		chainID,
		outputs,
		feeCalc,
		common.UnionOptions(b.options, options)...,
	)
}
