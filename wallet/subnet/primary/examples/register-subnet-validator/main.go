// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/ava-labs/avalanchego/api/info"
	"github.com/ava-labs/avalanchego/genesis"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/message"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	"github.com/ava-labs/avalanchego/vms/secp256k1fx"
	"github.com/ava-labs/avalanchego/wallet/subnet/primary"
)

func main() {
	key := genesis.EWOQKey
	uri := "http://localhost:9650"
	kc := secp256k1fx.NewKeychain(key)
<<<<<<< HEAD
	subnetID := ids.FromStringOrPanic("2eZYSgCU738xN7aRw47NsBUPqnKkoqJMYUJexTsX19VdTNSZc9")
	chainID := ids.FromStringOrPanic("JCuZD3rg8dEEeY5cJwQxwSj6Shqo8DtRynBCt1nXkBXCouZVw")
=======
	subnetID := ids.FromStringOrPanic("2DeHa7Qb6sufPkmQcFWG2uCd4pBPv9WB6dkzroiMQhd1NSRtof")
	chainID := ids.FromStringOrPanic("2BMFrJ9xeh5JdwZEx6uuFcjfZC2SV2hdbMT8ee5HrvjtfJb5br")
>>>>>>> implement-acp-77--set-subnet-validator-weight-tx
	addressHex := ""
	weight := uint64(1)

	address, err := hex.DecodeString(addressHex)
	if err != nil {
		log.Fatalf("failed to decode address %q: %s\n", addressHex, err)
	}

	ctx := context.Background()
	infoClient := info.NewClient(uri)

	skBytes, err := os.ReadFile("/Users/stephen/.avalanchego/staking/signer.key")
	if err != nil {
		log.Fatalf("failed to read signer key: %s\n", err)
	}

	sk, err := bls.SecretKeyFromBytes(skBytes)
	if err != nil {
		log.Fatalf("failed to parse secret key: %s\n", err)
	}

	nodeInfoStartTime := time.Now()
	nodeID, nodePoP, err := infoClient.GetNodeID(ctx)
	if err != nil {
		log.Fatalf("failed to fetch node IDs: %s\n", err)
	}
	log.Printf("fetched node ID %s in %s\n", nodeID, time.Since(nodeInfoStartTime))

	// MakeWallet fetches the available UTXOs owned by [kc] on the network that
	// [uri] is hosting and registers [subnetID].
	walletSyncStartTime := time.Now()
	wallet, err := primary.MakeWallet(ctx, &primary.WalletConfig{
		URI:          uri,
		AVAXKeychain: kc,
		EthKeychain:  kc,
	})
	if err != nil {
		log.Fatalf("failed to initialize wallet: %s\n", err)
	}
	log.Printf("synced wallet in %s\n", time.Since(walletSyncStartTime))

	// Get the P-chain wallet
	pWallet := wallet.P()
	context := pWallet.Builder().Context()

	addressedCallPayload, err := message.NewRegisterSubnetValidator(
		subnetID,
		nodeID,
		nodePoP.PublicKey,
		uint64(time.Now().Add(5*time.Minute).Unix()),
		message.PChainOwner{},
		message.PChainOwner{},
		weight,
	)
	if err != nil {
		log.Fatalf("failed to create RegisterSubnetValidator message: %s\n", err)
	}
	addressedCallPayloadJSON, err := json.MarshalIndent(addressedCallPayload, "", "\t")
	if err != nil {
		log.Fatalf("failed to marshal RegisterSubnetValidator message: %s\n", err)
	}
	log.Println(string(addressedCallPayloadJSON))

	addressedCall, err := payload.NewAddressedCall(
		address,
		addressedCallPayload.Bytes(),
	)
	if err != nil {
		log.Fatalf("failed to create AddressedCall message: %s\n", err)
	}

	unsignedWarp, err := warp.NewUnsignedMessage(
		context.NetworkID,
		chainID,
		addressedCall.Bytes(),
	)
	if err != nil {
		log.Fatalf("failed to create unsigned Warp message: %s\n", err)
	}

	signers := set.NewBits()
	signers.Add(0) // [signers] has weight from [vdr[0]]

	unsignedBytes := unsignedWarp.Bytes()
	sig := bls.Sign(sk, unsignedBytes)
	sigBytes := [bls.SignatureLen]byte{}
	copy(sigBytes[:], bls.SignatureToBytes(sig))

	warp, err := warp.NewMessage(
		unsignedWarp,
		&warp.BitSetSignature{
			Signers:   signers.Bytes(),
			Signature: sigBytes,
		},
	)
	if err != nil {
		log.Fatalf("failed to create Warp message: %s\n", err)
	}

	registerSubnetValidatorStartTime := time.Now()
	registerSubnetValidatorTx, err := pWallet.IssueRegisterSubnetValidatorTx(
		units.Avax,
		nodePoP.ProofOfPossession,
		warp.Bytes(),
	)
	if err != nil {
		log.Fatalf("failed to issue register subnet validator transaction: %s\n", err)
	}

	validationID := addressedCallPayload.ValidationID()
	log.Printf("registered new subnet validator %s to subnet %s with txID %s as validationID %s in %s\n", nodeID, subnetID, registerSubnetValidatorTx.ID(), validationID, time.Since(registerSubnetValidatorStartTime))
}
