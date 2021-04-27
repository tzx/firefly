// Copyright © 2021 Kaleido, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package blockchain

import (
	"context"
)

// Plugin is the interface implemented by each blockchain plugin
type Plugin interface {

	// ConfigInterface returns the structure into which to marshal the plugin config
	ConfigInterface() interface{}

	// Init initializes the plugin, with the config marshaled into the return of ConfigInterface
	// Returns the supported featureset of the interface
	Init(ctx context.Context, config interface{}, events Events) (*Capabilities, error)

	// SubmitBroadcastBatch sequences a broadcast globally to all viewers of the blockchain
	// The returned tracking ID will be used to correlate with any subsequent transaction tracking updates
	SubmitBroadcastBatch(identity string, broadcast BroadcastBatch) (txTrackingID string, err error)
}

// BlockchainEvents is the interface provided to the blockchain plugin, to allow it to pass events back to firefly.
//
// All blockchain-sequenced events MUST be delivered to the same firefly core instance (within a cluster), to allow
// deterministic ordering of the event delivery to subscribed event handlers (apps etc.).
// If that firefly core instance terminates/disconnects from the blockchain remote agent, the stream should
// fail-over to another instance. Then all all un-confirmed messages will be replayed (in the correct sequence) to that node.
//
// One example of how this can be acheived, is with a singleton instance of the event stream runtime (with HA failover)
// and a WebSocket connection from the firefly core runtime to that instance. Then the remote event stream runtime can
// choose exactly one of those connected WebSockets to dispatch events to, and in the case the websocket disconnects
// pick the next available connection and re-deliver anything that was missed.
type Events interface {
	// TransactionUpdate notifies firefly of an update to a transaction. Only success/failure and errorMessage (for errors) are modeled.
	// additionalInfo can be used to add opaque protocol specific JSON from the plugin (protocol transaction ID etc.)
	// Note this is an optional hook information, and stored separately to the confirmation of the actual event that was being submitted/sequenced.
	// Only the party submitting the transaction will see this data.
	TransactionUpdate(txTrackingID string, txState TransactionState, errorMessage string, additionalInfo map[string]interface{})

	// SequencedBroadcastBatch notifies on the arrival of a sequenced batch of broadcast messages, which might have been
	// submitted by us, or by any other authorized party in the network.
	// additionalInfo can be used to add opaque protocol specific JSON from the plugin (block numbers etc.)
	SequencedBroadcastBatch(batch BroadcastBatch, additionalInfo map[string]interface{})
}

// BlockchainCapabilities the supported featureset of the blockchain
// interface implemented by the plugin, with the specified config
type Capabilities struct {
	// GlobalSequencer means submitting an ordered piece of data visible to all
	// participants of the network (requires an all-participant chain)
	GlobalSequencer bool
}

// TransactionState is the only architecturally significant thing that Firefly tracks on blockchain transactions.
// All other data is consider protocol specific, and hence stored as opaque data.
type TransactionState string

const (
	// TransactionStateSubmitted the transaction has been submitted
	TransactionStateSubmitted TransactionState = "submitted"
	// TransactionStateSubmitted the transaction is considered final per the rules of the blockchain technnology
	TransactionStateConfirmed TransactionState = "confirmed"
	// TransactionStateSubmitted the transaction has encountered, and is unlikely to ever become final on the blockchain. However, it is not impossible it will still be mined.
	TransactionStateFailed TransactionState = "error"
)

// BroadcastBatch is the set of data pinned to the blockchain for a batch of broadcasts.
// Broadcasts are batched where possible, as the storage of the off-chain data is expensive as it must be propagated to all members
// of the network (via a technology like IPFS).
type BroadcastBatch struct {

	// Timestamp is the time of submission, from the pespective of the submitter
	Timestamp uint64

	// BatchPaylodRef is a 32 byte fixed length binary value that can be passed to the storage interface to retrieve the payload
	BatchPaylodRef HexUUID

	// BatchID is the id of the batch - writing this in plain text to the blockchain makes for easy correlation on-chain/off-chain
	BatchID Bytes32
}