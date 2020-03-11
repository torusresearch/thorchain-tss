package tss

import (
	"encoding/base64"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	bkeygen "github.com/binance-chain/tss-lib/ecdsa/keygen"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/libp2p/go-libp2p-core/peer"
	maddr "github.com/multiformats/go-multiaddr"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tcrypto "github.com/tendermint/tendermint/crypto"

	"gitlab.com/thorchain/tss/go-tss/common"
	"gitlab.com/thorchain/tss/go-tss/keygen"
	"gitlab.com/thorchain/tss/go-tss/keysign"
	"gitlab.com/thorchain/tss/go-tss/messages"
	"gitlab.com/thorchain/tss/go-tss/p2p"
	"gitlab.com/thorchain/tss/go-tss/storage"
)

// TssServer is the structure that can provide all keysign and key gen features
type TssServer struct {
	conf              common.TssConfig
	logger            zerolog.Logger
	Status            common.TssStatus
	p2pCommunication  *p2p.Communication
	localNodePubKey   string
	preParams         *bkeygen.LocalPreParams
	wg                sync.WaitGroup
	tssKeyGenLocker   *sync.Mutex
	stopChan          chan struct{}
	homeBase          string
	partyCoordinator  *p2p.PartyCoordinator
	stateManager      storage.LocalStateManager
	signatureNotifier *keysign.SignatureNotifier
}

// NewTss create a new instance of Tss
func NewTss(
	bootstrapPeers []maddr.Multiaddr,
	p2pPort int,
	priKey tcrypto.PrivKey,
	rendezvous,
	baseFolder string,
	conf common.TssConfig,
	preParams *bkeygen.LocalPreParams,
) (*TssServer, error) {
	pubKey, err := sdk.Bech32ifyAccPub(priKey.PubKey())
	if err != nil {
		return nil, fmt.Errorf("fail to genearte the key: %w", err)
	}

	comm, err := p2p.NewCommunication(rendezvous, bootstrapPeers, p2pPort)
	if err != nil {
		return nil, fmt.Errorf("fail to create communication layer: %w", err)
	}

	// When using the keygen party it is recommended that you pre-compute the
	// "safe primes" and Paillier secret beforehand because this can take some
	// time.
	// This code will generate those parameters using a concurrency limit equal
	// to the number of available CPU cores.
	if preParams == nil || !preParams.Validate() {
		preParams, err = bkeygen.GeneratePreParams(conf.PreParamTimeout)
		if err != nil {
			return nil, fmt.Errorf("fail to generate pre parameters: %w", err)
		}
	}
	if !preParams.Validate() {
		return nil, errors.New("invalid preparams")
	}

	priKeyRawBytes, err := getPriKeyRawBytes(priKey)
	if err != nil {
		return nil, fmt.Errorf("fail to get private key")
	}
	if err := comm.Start(priKeyRawBytes); nil != err {
		return nil, fmt.Errorf("fail to start p2p network: %w", err)
	}
	pc := p2p.NewPartyCoordinator(comm.GetHost())
	stateManager, err := storage.NewFileStateMgr(baseFolder)
	if err != nil {
		return nil, fmt.Errorf("fail to create file state manager")
	}
	sn := keysign.NewSignatureNotifier(comm.GetHost())
	tssServer := TssServer{
		conf:   conf,
		logger: log.With().Str("module", "tss").Logger(),
		Status: common.TssStatus{
			Starttime: time.Now(),
		},
		p2pCommunication:  comm,
		localNodePubKey:   pubKey,
		preParams:         preParams,
		tssKeyGenLocker:   &sync.Mutex{},
		stopChan:          make(chan struct{}),
		partyCoordinator:  pc,
		stateManager:      stateManager,
		signatureNotifier: sn,
	}

	return &tssServer, nil
}

// Start Tss server
func (t *TssServer) Start() error {
	log.Info().Msg("Starting the TSS servers")
	t.Status.Starttime = time.Now()
	t.signatureNotifier.Start()
	go t.p2pCommunication.ProcessBroadcast()
	return nil
}

// Stop Tss server
func (t *TssServer) Stop() {
	close(t.stopChan)
	// stop the p2p and finish the p2p wait group
	err := t.p2pCommunication.Stop()
	if err != nil {
		t.logger.Error().Msgf("error in shutdown the p2p server")
	}
	t.signatureNotifier.Stop()
	t.partyCoordinator.Stop()
	log.Info().Msg("The Tss and p2p server has been stopped successfully")
}

func (t *TssServer) requestToMsgId(request interface{}) (string, error) {
	var dat []byte
	switch value := request.(type) {
	case keygen.Request:
		keyAccumulation := ""
		keys := value.Keys
		sort.Strings(keys)
		for _, el := range keys {
			keyAccumulation += el
		}
		dat = []byte(keyAccumulation)
	case keysign.Request:
		msgToSign, err := base64.StdEncoding.DecodeString(value.Message)
		if err != nil {
			t.logger.Error().Err(err).Msg("error in decode the keysign req")
			return "", err
		}
		dat = msgToSign
	default:
		t.logger.Error().Msg("unknown request type")
		return "", errors.New("unknown request type")
	}
	return common.MsgToHashString(dat)
}

func (t *TssServer) joinParty(msgID string, messageToSign []byte, keys []string) (*messages.JoinPartyResponse, peer.ID, error) {
	noResponse := &messages.JoinPartyResponse{
		Type: messages.JoinPartyResponse_Unknown,
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	peerIDs, err := GetPeerIDsFromPubKeys(keys)
	if err != nil {
		return noResponse, "", fmt.Errorf("fail to convert pub key to peer id: %w", err)
	}
	totalNodes := int32(len(keys))
	leader, err := p2p.LeaderNode(messageToSign, totalNodes)
	if err != nil {
		return noResponse, "", fmt.Errorf("fail to get leader node: %w", err)
	}

	leaderPeerID, err := GetPeerIDFromPubKey(keys[leader])
	if err != nil {
		return noResponse, leaderPeerID, fmt.Errorf("fail to get peer id from node pubkey: %w", err)
	}
	t.logger.Info().Msgf("leader peer: %s", leaderPeerID)
	joinPartyReq := &messages.JoinPartyRequest{
		ID: msgID,
	}
	msg, err := t.partyCoordinator.JoinPartyWithRetry(leaderPeerID, joinPartyReq, peerIDs, totalNodes)
	return msg, leaderPeerID, err
}

// GetLocalPeerID return the local peer
func (t *TssServer) GetLocalPeerID() string {
	return t.p2pCommunication.GetLocalPeerID()
}

// GetStatus return the TssStatus
func (t *TssServer) GetStatus() common.TssStatus {
	return t.Status
}
