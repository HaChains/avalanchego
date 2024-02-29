// Copyright (C) 2019-2024, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package network

import (
	"time"

	"github.com/ava-labs/avalanchego/utils/units"
)

var DefaultConfig = Config{
	MaxValidatorSetStaleness:                    time.Minute,
	TargetGossipSize:                            20 * units.KiB,
	PushGossipFirstNumValidators:                100,
	PushGossipFirstNumPeers:                     0,
	PushGossipFollowupNumValidators:             10,
	PushGossipFollowupNumPeers:                  0,
	PushGossipDiscardedCacheSize:                1024,
	PushGossipMaxRegossipFrequency:              10 * time.Second,
	PushGossipFrequency:                         500 * time.Millisecond,
	PullGossipPollSize:                          1,
	PullGossipFrequency:                         1500 * time.Millisecond,
	PullGossipThrottlingPeriod:                  10 * time.Second,
	PullGossipThrottlingLimit:                   2,
	ExpectedBloomFilterElements:                 8 * 1024,
	ExpectedBloomFilterFalsePositiveProbability: .01,
	MaxBloomFilterFalsePositiveProbability:      .05,
}

type Config struct {
	// MaxValidatorSetStaleness limits how old of a validator set the network
	// will use for peer sampling and rate limiting.
	MaxValidatorSetStaleness time.Duration `json:"max-validator-set-staleness"`
	// TargetGossipSize is the number of bytes that will be attempted to be
	// sent when pushing transactions and when responded to transaction pull
	// requests.
	TargetGossipSize int `json:"target-gossip-size"`
	// PushGossipFirstNumValidators is the number of validators to push
	// transactions to in the first round of gossip.
	PushGossipFirstNumValidators int `json:"push-gossip-first-num-validators"`
	// PushGossipFirstNumPeers is the number of peers to push transactions to in
	// the first round of gossip.
	PushGossipFirstNumPeers int `json:"push-gossip-first-num-peers"`
	// PushGossipFollowupNumValidators is the number of validators to push
	// transactions to after the first round of gossip.
	PushGossipFollowupNumValidators int `json:"push-gossip-followup-num-validators"`
	// PushGossipFollowupNumPeers is the number of peers to push transactions to
	// after the first round of gossip.
	PushGossipFollowupNumPeers int `json:"push-gossip-followup-num-peers"`
	// PushGossipDiscardedCacheSize is the number of txIDs to cache to avoid
	// pushing transactions that were recently dropped from the mempool.
	PushGossipDiscardedCacheSize int `json:"push-gossip-discarded-cache-size"`
	// PushGossipMaxRegossipFrequency is the limit for how frequently a
	// transaction will be push gossiped.
	PushGossipMaxRegossipFrequency time.Duration `json:"push-gossip-max-regossip-frequency"`
	// PushGossipFrequency is how frequently rounds of push gossip are
	// performed.
	PushGossipFrequency time.Duration `json:"push-gossip-frequency"`
	// PullGossipPollSize is the number of validators to sample when performing
	// a round of pull gossip.
	PullGossipPollSize int `json:"pull-gossip-poll-size"`
	// PullGossipFrequency is how frequently rounds of pull gossip are
	// performed.
	PullGossipFrequency time.Duration `json:"pull-gossip-frequency"`
	// PullGossipThrottlingPeriod is how large of a window the throttler should
	// use.
	PullGossipThrottlingPeriod time.Duration `json:"pull-gossip-throttling-period"`
	// PullGossipThrottlingLimit is the number of pull querys that are allowed
	// by a validator in every throttling window.
	PullGossipThrottlingLimit int `json:"pull-gossip-throttling-limit"`
	// ExpectedBloomFilterElements is the number of elements to expect when
	// creating a new bloom filter. The larger this number is, the larger the
	// bloom filter will be.
	ExpectedBloomFilterElements int `json:"expected-bloom-filter-elements"`
	// ExpectedBloomFilterFalsePositiveProbability is the expected probability
	// of a false positive after having inserted ExpectedBloomFilterElements
	// into a bloom filter. The smaller this number is, the larger the bloom
	// filter will be.
	ExpectedBloomFilterFalsePositiveProbability float64 `json:"expected-bloom-filter-false-positive-probability"`
	// MaxBloomFilterFalsePositiveProbability is used to determine when the
	// bloom filter should be refreshed. Once the expected probability of a
	// false positive exceeds this value, the bloom filter will be regenerated.
	// The smaller this number is, the more frequently that the bloom filter
	// will be regenerated.
	MaxBloomFilterFalsePositiveProbability float64 `json:"max-bloom-filter-false-positive-probability"`
}
