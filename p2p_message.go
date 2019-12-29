package tss

import (
	"fmt"
	"sync"

	btss "github.com/binance-chain/tss-lib/tss"
	"github.com/libp2p/go-libp2p-core/peer"
)

// THORChainTSSMessageType  represent the messgae type used in THORChain TSS
type THORChainTSSMessageType uint8

const (
	// TSSKeyGenMsg is the message directly generated by tss-lib package
	TSSKeyGenMsg THORChainTSSMessageType = iota
	// TSSKeySignMsg is the message directly generated by tss lib for sign
	TSSKeySignMsg
	// TSSKeyGenVerMsg is the message we create on top to make sure everyone received the same message
	TSSKeyGenVerMsg
	// TSSKeySignVerMsg is the message we create to make sure every party receive the same broadcast message
	TSSKeySignVerMsg
	// Unknown message
	Unknown
)

// String implement fmt.Stringer
func (msgType THORChainTSSMessageType) String() string {
	switch msgType {
	case TSSKeyGenMsg:
		return "TSSKeyGenMsg"
	case TSSKeySignMsg:
		return "TSSKeySignMsg"
	case TSSKeyGenVerMsg:
		return "TSSKeyGenVerMsg"
	case TSSKeySignVerMsg:
		return "TSSKeySignVerMsg"
	default:
		return "Unknown"
	}
}

// WrappedMessage is a message with type in it
type WrappedMessage struct {
	MessageType THORChainTSSMessageType `json:"message_type"`
	Payload     []byte                  `json:"payload"`
}

// BroadcastWrappedMessage used to send to broadcast channel
type BroadcastWrapppedMessage struct {
	peers      []peer.ID
	wrappedMsg *WrappedMessage
}

// BroadcastConfirmMessage is used to broadcast to all parties what message they receive
type BroadcastConfirmMessage struct {
	P2PID string `json:"P2PID"`
	Key   string `json:"key"`
	Hash  string `json:"hash"`
}

// WireMessage the message that produced by tss-lib package
type WireMessage struct {
	Routing   *btss.MessageRouting `json:"routing"`
	RoundInfo string               `json:"round_info"`
	Message   []byte               `json:"message"`
}

// GetCacheKey return the key we used to cache it locally
func (m *WireMessage) GetCacheKey() string {
	return fmt.Sprintf("%s-%s", m.Routing.From.Id, m.RoundInfo)
}

// LocalCacheItem used to cache the unconfirmed broadcast message
type LocalCacheItem struct {
	Msg           *WireMessage
	Hash          string
	lock          *sync.Mutex
	ConfirmedList map[string]string
}

// UpdateConfirmList add the given party's hash into the confirm list
func (l *LocalCacheItem) UpdateConfirmList(P2PID, hash string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.ConfirmedList[P2PID] = hash
}

// TotalConfirmParty number of parties that already confirmed their hash
func (l *LocalCacheItem) TotalConfirmParty() int {
	l.lock.Lock()
	defer l.lock.Unlock()
	return len(l.ConfirmedList)
}
