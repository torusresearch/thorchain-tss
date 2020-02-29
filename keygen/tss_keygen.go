package keygen

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/binance-chain/tss-lib/crypto"
	bkeygen "github.com/binance-chain/tss-lib/ecdsa/keygen"
	btss "github.com/binance-chain/tss-lib/tss"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	cryptokey "github.com/tendermint/tendermint/crypto"

	"gitlab.com/thorchain/tss/go-tss/common"
	"gitlab.com/thorchain/tss/go-tss/p2p"
)

type TssKeyGen struct {
	logger          zerolog.Logger
	priKey          cryptokey.PrivKey
	preParams       *bkeygen.LocalPreParams
	tssCommonStruct *common.TssCommon
	stopChan        chan struct{} // channel to indicate whether we should stop
	homeBase        string
	syncMsg         chan *p2p.Message
	localParty      *btss.PartyID
	keygenCurrent   *string
	commStopChan    chan struct{}
}

func NewTssKeyGen(homeBase, localP2PID string, conf common.TssConfig, privKey cryptokey.PrivKey, broadcastChan chan *p2p.BroadcastMsgChan, stopChan chan struct{}, preParam *bkeygen.LocalPreParams, keygenCurrent *string, msgID string) TssKeyGen {
	return TssKeyGen{
		logger:          log.With().Str("module", "keyGen").Logger(),
		priKey:          privKey,
		preParams:       preParam,
		tssCommonStruct: common.NewTssCommon(localP2PID, broadcastChan, conf, msgID),
		stopChan:        stopChan,
		homeBase:        homeBase,
		syncMsg:         make(chan *p2p.Message),
		localParty:      nil,
		keygenCurrent:   keygenCurrent,
		commStopChan:    make(chan struct{}),
	}
}

func (tKeyGen *TssKeyGen) GetTssKeyGenChannels() (chan *p2p.Message, chan *p2p.Message) {
	return tKeyGen.tssCommonStruct.TssMsg, tKeyGen.syncMsg
}

func (tKeyGen *TssKeyGen) GetTssCommonStruct() *common.TssCommon {
	return tKeyGen.tssCommonStruct
}

func (tKeyGen *TssKeyGen) GenerateNewKey(keygenReq KeyGenReq) (*crypto.ECPoint, error) {
	tKeyGen.logger.Info().Msgf("----keygen parties are %v\n", keygenReq.Keys)
	pubKey, err := sdk.Bech32ifyAccPub(tKeyGen.priKey.PubKey())
	if err != nil {
		return nil, fmt.Errorf("fail to genearte the key: %w", err)
	}
	partiesID, localPartyID, err := common.GetParties(keygenReq.Keys, nil, pubKey)
	if err != nil {
		return nil, fmt.Errorf("fail to get keygen parties: %w", err)
	}
	keyGenLocalStateItem := common.KeygenLocalStateItem{
		ParticipantKeys: keygenReq.Keys,
		LocalPartyKey:   pubKey,
	}

	threshold, err := common.GetThreshold(len(partiesID))
	if err != nil {
		return nil, err
	}
	ctx := btss.NewPeerContext(partiesID)
	params := btss.NewParameters(ctx, localPartyID, len(partiesID), threshold)
	outCh := make(chan btss.Message, len(partiesID))
	endCh := make(chan bkeygen.LocalPartySaveData, len(partiesID))
	errChan := make(chan struct{})
	if tKeyGen.preParams == nil {
		tKeyGen.logger.Error().Err(err).Msg("error, empty pre-parameters")
		return nil, errors.New("error, empty pre-parameters")
	}
	keyGenParty := bkeygen.NewLocalParty(params, outCh, endCh, *tKeyGen.preParams)
	partyIDMap := common.SetupPartyIDMap(partiesID)
	err = common.SetupIDMaps(partyIDMap, tKeyGen.tssCommonStruct.PartyIDtoP2PID)
	if err != nil {
		tKeyGen.logger.Error().Msgf("error in creating mapping between partyID and P2P ID")
		return nil, err
	}

	tKeyGen.tssCommonStruct.SetPartyInfo(&common.PartyInfo{
		Party:      keyGenParty,
		PartyIDMap: partyIDMap,
	})
	tKeyGen.tssCommonStruct.P2PPeers = common.GetPeersID(tKeyGen.tssCommonStruct.PartyIDtoP2PID, tKeyGen.tssCommonStruct.GetLocalPeerID())
	// we set the coordinator of the keygen
	tKeyGen.tssCommonStruct.Coordinator, err = tKeyGen.tssCommonStruct.GetCoordinator("keygen")
	if err != nil {
		tKeyGen.logger.Error().Err(err).Msg("error in get the coordinator")
		return nil, err
	}
	standbyPeers, errNodeSync := tKeyGen.tssCommonStruct.NodeSync(tKeyGen.syncMsg, p2p.TSSKeyGenSync)
	if errNodeSync != nil {
		tKeyGen.logger.Error().Err(err).Msg("node sync error")
		tKeyGen.logger.Error().Err(err).Msgf("the nodes online are +%v", standbyPeers)
		standbyPeers = common.RemoveCoordinator(standbyPeers, tKeyGen.GetTssCommonStruct().Coordinator)
		_, blamePubKeys, err := tKeyGen.tssCommonStruct.GetBlamePubKeysLists(standbyPeers)
		if err != nil {
			tKeyGen.logger.Error().Err(err).Msgf("error in get blame node pubkey + %v", errNodeSync)
			return nil, errNodeSync
		}
		tKeyGen.tssCommonStruct.BlamePeers.SetBlame(common.BlameNodeSyncCheck, blamePubKeys)
		return nil, errNodeSync
	}
	// start keygen
	go func() {
		defer tKeyGen.logger.Info().Msg("keyGenParty started")
		if err := keyGenParty.Start(); nil != err {
			tKeyGen.logger.Error().Err(err).Msg("fail to start keygen party")
			close(errChan)
		}
	}()
	go tKeyGen.processInboundMessages()
	r, err := tKeyGen.processKeyGen(errChan, outCh, endCh, keyGenLocalStateItem)
	return r, err
}

func (tKeyGen *TssKeyGen) processInboundMessages() {
	tKeyGen.logger.Info().Msg("start processing inbound messages")
	defer tKeyGen.logger.Info().Msg("stop processing inbound messages")
	for {
		select {
		case <-tKeyGen.commStopChan:
			return
		case m, ok := <-tKeyGen.tssCommonStruct.TssMsg:
			if !ok {
				return
			}
			var wrappedMsg p2p.WrappedMessage
			if err := json.Unmarshal(m.Payload, &wrappedMsg); nil != err {
				tKeyGen.logger.Error().Err(err).Msg("fail to unmarshal wrapped message bytes")
				continue
			}

			if err := tKeyGen.tssCommonStruct.ProcessOneMessage(&wrappedMsg, m.PeerID.String()); err != nil {
				tKeyGen.logger.Error().Err(err).Msg("fail to process the received message")
			}
		}
	}
}


func (tKeyGen *TssKeyGen) processKeyGen(errChan chan struct{}, outCh <-chan btss.Message, endCh <-chan bkeygen.LocalPartySaveData, keyGenLocalStateItem common.KeygenLocalStateItem) (*crypto.ECPoint, error) {
	defer tKeyGen.logger.Info().Msg("finished keygen process")
	tKeyGen.logger.Info().Msg("start to read messages from local party")
	defer close(tKeyGen.commStopChan)
	tssConf := tKeyGen.tssCommonStruct.GetConf()
	for {
		select {
		case <-errChan: // when keyGenParty return
			tKeyGen.logger.Error().Msg("key gen failed")
			close(tKeyGen.tssCommonStruct.TssMsg)
			return nil, errors.New("error channel closed fail to start local party")

		case <-tKeyGen.stopChan: // when TSS processor receive signal to quit
			close(tKeyGen.tssCommonStruct.TssMsg)
			return nil, errors.New("received exit signal")

		case <-time.After(tssConf.KeyGenTimeout):
			// we bail out after KeyGenTimeoutSeconds

			tKeyGen.logger.Error().Msgf("fail to generate message with %s", tssConf.KeyGenTimeout.String())
			tssCommonStruct := tKeyGen.GetTssCommonStruct()
			localCachedItems := tssCommonStruct.TryGetAllLocalCached()
			blamePeers, err := tssCommonStruct.TssTimeoutBlame(localCachedItems)
			if err != nil {
				tKeyGen.logger.Error().Err(err).Msg("fail to get the blamed peers")
				tssCommonStruct.BlamePeers.SetBlame(common.BlameTssTimeout, nil)
				return nil, fmt.Errorf("fail to get the blamed peers %w", common.ErrTssTimeOut)
			}
			tssCommonStruct.BlamePeers.SetBlame(common.BlameTssTimeout, blamePeers)
			return nil, common.ErrTssTimeOut

		case msg := <-outCh:
			tKeyGen.logger.Debug().Msgf(">>>>>>>>>>msg: %s", msg.String())
			// for the sake of performance, we do not lock the status update
			// we report a rough status of current round
			*tKeyGen.keygenCurrent = msg.Type()
			err := tKeyGen.tssCommonStruct.ProcessOutCh(msg, p2p.TSSKeyGenMsg)
			if err != nil {
				return nil, err
			}

		case msg := <-endCh:
			tKeyGen.logger.Debug().Msgf("we have done the keygen %s", msg.ECDSAPub.Y().String())
			if err := tKeyGen.AddLocalPartySaveData(tKeyGen.homeBase, msg, keyGenLocalStateItem); nil != err {
				return nil, fmt.Errorf("fail to save key gen result to local store: %w", err)
			}
			return msg.ECDSAPub, nil
		}
	}
}
