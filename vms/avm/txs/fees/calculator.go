// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package fees

import "github.com/ava-labs/avalanchego/vms/avm/txs"

// Calculator is the interfaces that any fee Calculator must implement
type Calculator interface {
	CalculateFee(tx *txs.Tx) (uint64, error)
}
