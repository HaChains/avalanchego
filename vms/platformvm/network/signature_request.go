// Copyright (C) 2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package network

import (
	"context"
	"fmt"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/network/p2p"
	"github.com/ava-labs/avalanchego/proto/pb/sdk"
	"github.com/ava-labs/avalanchego/snow/engine/common"
	"github.com/ava-labs/avalanchego/utils/crypto/bls"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/messages"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
	"google.golang.org/protobuf/proto"
)

const (
	ErrFailedToParse = iota
	ErrFailedToGetSig
	ErrFailedToMarshal
	ErrUnsupportedWarpMessageType
	ErrInvalidCodecVersion
	ErrFailedToSignMessage
	ErrInvalidSignatureLength
)

var _ p2p.Handler = (*signatureRequestHandler)(nil)

type signatureRequestHandler struct {
	p2p.NoOpHandler
	signer warp.Signer
}

// TODO: This should be allowed only for local networks
func (s signatureRequestHandler) AppRequest(
	_ context.Context,
	_ ids.NodeID,
	_ time.Time,
	requestBytes []byte,
) ([]byte, *common.AppError) {
	// Per ACP-118, the requestBytes are the serialized form of
	// sdk.SignatureRequest.
	req := new(sdk.SignatureRequest)
	if err := proto.Unmarshal(requestBytes, req); err != nil {
		return nil, &common.AppError{
			Code:    ErrFailedToParse,
			Message: "failed to unmarshal request: " + err.Error(),
		}
	}

	unsignedMessage, err := warp.ParseUnsignedMessage(req.Message)
	if err != nil {
		return nil, &common.AppError{
			Code:    ErrFailedToParse,
			Message: "failed to parse unsigned message: " + err.Error(),
		}
	}
	parsed, err := payload.Parse(unsignedMessage.Payload)
	if err != nil {
		return nil, &common.AppError{
			Code:    ErrFailedToParse,
			Message: "failed to parse payload: " + err.Error(),
		}
	}

	var sig []byte
	switch p := parsed.(type) {
	case *payload.AddressedCall:
		msg, err := payload.ParseAddressedCall(unsignedMessage.Payload)
		if err != nil {
			return nil, &common.AppError{
				Code:    ErrFailedToParse,
				Message: "failed to parse addressed call: " + err.Error(),
			}
		}
		// Check that the addressed call payload is a registered Warp message type
		var dst interface{}
		ver, err := messages.Codec.Unmarshal(msg.Payload, dst)
		if err != nil {
			return nil, &common.AppError{
				Code:    ErrUnsupportedWarpMessageType,
				Message: "unsupported warp message type",
			}
		}
		if ver != messages.CodecVersion {
			return nil, &common.AppError{
				Code:    ErrInvalidCodecVersion,
				Message: "invalid codec version",
			}
		}
		sig, err = s.signer.Sign(unsignedMessage)
		if err != nil {
			return nil, &common.AppError{
				Code:    ErrFailedToSignMessage,
				Message: "failed to sign message: " + err.Error(),
			}
		}
		if len(sig) != bls.SignatureLen {
			return nil, &common.AppError{
				Code:    ErrInvalidSignatureLength,
				Message: "invalid signature length",
			}
		}
	default:
		return nil, &common.AppError{
			Code:    ErrFailedToParse,
			Message: fmt.Sprintf("unknown payload type: %T", p),
		}
	}
	if err != nil {
		return nil, &common.AppError{
			Code:    ErrFailedToGetSig,
			Message: "failed to get signature: " + err.Error(),
		}
	}

	// Per ACP-118, the responseBytes are the serialized form of
	// sdk.SignatureResponse.
	resp := &sdk.SignatureResponse{Signature: sig[:]}
	respBytes, err := proto.Marshal(resp)
	if err != nil {
		return nil, &common.AppError{
			Code:    ErrFailedToMarshal,
			Message: "failed to marshal response: " + err.Error(),
		}
	}
	return respBytes, nil
}