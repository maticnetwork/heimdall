package pier

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/common"
	httpClient "github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/bor"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

const (
	lastSpanKey = "span-key" // storage key
)

type SpanService struct {
	// Base service
	common.BaseService

	// storage client
	storageClient *leveldb.DB

	// Rootchain instance
	validatorSet *rootchain.Rootchain

	// header listener subscription
	cancelSpanService context.CancelFunc

	// contract caller
	contractConnector helper.ContractCaller

	// cli context
	cliCtx cliContext.CLIContext

	// queue connector
	queueConnector QueueConnector

	// http client to subscribe to
	httpClient *httpClient.HTTP
}

// NewSpanService returns new service object
func NewSpanService(cdc *codec.Codec, queueConnector QueueConnector, httpClient *httpClient.HTTP) *SpanService {
	// create logger
	logger := Logger.With("module", SpanServiceStr)

	// root chain instance
	rootchainInstance, err := helper.GetRootChainInstance()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}
	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastAsync

	// creating checkpointer object
	spanService := &SpanService{
		storageClient:     getBridgeDBInstance(viper.GetString(BridgeDBFlag)),
		validatorSet:      rootchainInstance,
		contractConnector: contractCaller,

		cliCtx:         cliCtx,
		queueConnector: queueConnector,
		httpClient:     httpClient,
	}

	spanService.BaseService = *common.NewBaseService(logger, SpanServiceStr, spanService)
	return spanService
}

// OnStart starts new block subscription
func (s *SpanService) OnStart() error {
	s.BaseService.OnStart() // Always call the overridden method.

	// create cancellable context
	spanCtx, cancelSpanService := context.WithCancel(context.Background())

	s.cancelSpanService = cancelSpanService

	// start polling for checkpoint in buffer
	go s.startPolling(spanCtx, 10*time.Second)

	// subscribed to new head
	s.Logger.Debug("Started Span service")
	return nil
}

// OnStop stops all necessary go routines
func (s *SpanService) OnStop() {
	s.BaseService.OnStop()
	s.httpClient.Stop()

	// cancel ack process
	s.cancelSpanService()
	// close bridge db instance
	closeBridgeDBInstance()
}

// polls heimdall and checks if new span needs to be proposed
func (s *SpanService) startPolling(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	// stop ticker when everything done
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if lastSpan, err := s.getLastSpan(); err == nil {
				if s.isSpanProposer(lastSpan) {
					go s.propose(lastSpan)
				}
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

// propose producers for next span if needed
func (s *SpanService) propose(lastSpan hmTypes.Span) {
	// call with last span on record + new span duration and see if it has been proposed
	currentBlock, err := s.GetCurrentChildBlock()
	if err != nil {
		s.Logger.Error("Unable to fetch current block")
		return
	}

	s.Logger.Debug("Fetched current child block", "CurrentChildBlock", currentBlock)
	if currentBlock >= lastSpan.StartBlock && currentBlock >= lastSpan.EndBlock {
		s.Logger.Info("Need to propose committee for next span")

		// send propose span
		s.ProposeNewSpan(lastSpan.ID+1, lastSpan.EndBlock+1)
	}

	// TODO
	// query validator set contract and check latest state
	// if its behind push onchain
}

// fetches last span processed in DB
func (s *SpanService) fetchLastSpan() (int, error) {
	hasLastSpan, err := s.storageClient.Has([]byte(lastSpanKey), nil)
	if hasLastSpan {
		lastSpanBytes, err := s.storageClient.Get([]byte(lastSpanKey), nil)
		if err != nil {
			s.Logger.Info("Error while fetching last span bytes from storage", "error", err)
			return 0, err
		}

		s.Logger.Debug("Got last block from bridge storage", "lastSpan", string(lastSpanBytes))
		if result, err := strconv.Atoi(string(lastSpanBytes)); err != nil {
			return 0, nil
		} else {
			return result, nil
		}
	}
	return 0, err
}

// checks span status
func (s *SpanService) getLastSpan() (spanStart hmTypes.Span, err error) {
	// fetch latest start block from heimdall via rest query
	result, err := FetchFromAPI(s.cliCtx, GetHeimdallServerEndpoint(LatestSpanURL))
	if err != nil {
		s.Logger.Error("Error while fetching latest span")
		return
	}

	var lastSpan hmTypes.Span
	err = json.Unmarshal(result.Result, &lastSpan)
	if err != nil {
		s.Logger.Error("Error unmarshalling", "error", err)
		return lastSpan, err
	}

	return lastSpan, nil
}

// GetCurrentChildBlock gets the
func (s *SpanService) GetCurrentChildBlock() (uint64, error) {
	childBlock, err := s.contractConnector.GetMaticChainBlock(nil)
	if err != nil {
		return 0, err
	}
	return childBlock.Number.Uint64(), nil
}

func (s *SpanService) isSpanProposer(lastSpan hmTypes.Span) bool {
	// sort validator address
	selectedProducers := types.SortValidatorByAddress(lastSpan.SelectedProducers)

	// get last validator as proposer
	proposer := selectedProducers[len(selectedProducers)-1]

	s.Logger.Debug("Fetched proposer for span", "proposer", proposer.Signer.String())
	if bytes.Equal(proposer.Signer.Bytes(), helper.GetAddress()) {
		return true
	}
	return false
}

// ProposeNewSpan proposes new span by sending transaction to heimdall
func (s *SpanService) ProposeNewSpan(id uint64, start uint64) {
	msg, err := s.fetchNextSpanDetails(id, start)
	if err != nil {
		s.Logger.Error("Unable to fetch next span details", "error", err)
		return
	}

	s.Logger.Info("Fetched information for next span", "NewSpan", msg)

	// tx builder
	txBldr := authTypes.NewTxBuilderFromCLI().
		WithTxEncoder(helper.GetTxEncoder()).
		WithChainID(helper.GetGenesisDoc().ChainID)

	txBytes, err := helper.GetSignedTxBytes(s.cliCtx, txBldr, []sdk.Msg{msg})
	if err != nil {
		s.Logger.Error("Error creating tx bytes", "error", err)
		return
	}

	resp, err := helper.BroadcastTxBytes(s.cliCtx, txBytes, client.BroadcastSync)
	if err != nil {
		s.Logger.Error("Unable to send propose span to heimdall", "Error", err, "StartBlock", msg.StartBlock, "EndBlock", msg.EndBlock, "ChainID", msg.ChainID)
		return
	}

	// subscribe to tx
	go s.SubscribeToTx(txBytes, msg.StartBlock, msg.EndBlock)
	// send to bor

	s.Logger.Info("Transaction sent to heimdall", "TxHash", resp.TxHash)
}

func (s *SpanService) fetchNextSpanDetails(id uint64, start uint64) (msg bor.MsgProposeSpan, err error) {
	req, err := http.NewRequest("GET", GetHeimdallServerEndpoint(NextSpanInfoURL), nil)
	if err != nil {
		s.Logger.Error("Error creating a new request", "error", err)
		return
	}

	q := req.URL.Query()
	q.Add("span_id", strconv.FormatUint(id, 10))
	q.Add("start_block", strconv.FormatUint(start, 10))
	q.Add("chain_id", viper.GetString("bor-chain-id"))
	q.Add("proposer", helper.GetFromAddress(s.cliCtx).String())
	req.URL.RawQuery = q.Encode()

	// log url
	s.Logger.Debug("Sending request", "url", req.URL.String())

	result, err := FetchFromAPI(s.cliCtx, req.URL.String())
	if err != nil {
		Logger.Error("Error fetching proposers", "error", err)
		return
	}

	err = json.Unmarshal(result.Result, &msg)
	if err != nil {
		Logger.Error("Error unmarshalling propose tx msg ", "error", err)
		return
	}
	return msg, nil
}

// SubscribeToTx subscribes to a broadcasted Tx and waits for its commitment to a block
func (s *SpanService) SubscribeToTx(tx tmTypes.Tx, start, end uint64) error {
	data, err := WaitForOneEvent(tx, s.httpClient)
	if err != nil {
		s.Logger.Error("Unable to wait for tx", "error", err)
		return err
	}

	switch t := data.(type) {
	case tmTypes.EventDataTx:
		go s.DispatchProposal(t.Height, t.Tx.Hash(), tx)
	default:
		s.Logger.Info("No cases matched while trying to send propose new committee")
	}
	return nil
}

// DispatchProposal dispatches proposal
func (s *SpanService) DispatchProposal(height int64, txHash []byte, txBytes tmTypes.Tx) {
	// extraData
	votes, sigs, chainID, err := fetchVotes(height, s.httpClient)
	if err != nil {
		s.Logger.Error("Error fetching votes", "height", height)
		return
	}

	// proof
	tx, err := helper.QueryTxWithProof(s.cliCtx, txHash)
	fmt.Println("TxBytes: ", hex.EncodeToString(tx.Tx[4:]))
	fmt.Println("Leaf: ", hex.EncodeToString(tx.Proof.Leaf()))
	fmt.Println("Root: ", tx.Proof.RootHash.String())
	proofList := helper.GetMerkleProofList(&tx.Proof.Proof)

	var result []string
	for _, e := range proofList {
		result = append(result, hex.EncodeToString(e))
	}
	fmt.Println("Votes: ", hex.EncodeToString(helper.GetVoteBytes(votes, chainID)))
	fmt.Println("Sigs: ", hex.EncodeToString(sigs))
	fmt.Println("chainID", chainID)

	fmt.Println("data : ",
		fmt.Sprintf(`"0x%s","0x%s","0x%s","0x%s"`,
			hex.EncodeToString(helper.GetVoteBytes(votes, chainID)),
			hex.EncodeToString(sigs),
			hex.EncodeToString(tx.Tx[4:]),
			strings.Join(result, ""),
		))

	// print proof
	fmt.Println("Proof: ", strings.Join(result, ""))
	s.Logger.Info("txBytes comparison", "Param", hex.EncodeToString(txBytes), "ReceivedTx", hex.EncodeToString(tx.Tx), "trimmed", hex.EncodeToString(tx.Tx[4:]))
	s.contractConnector.CommitSpan(helper.GetVoteBytes(votes, chainID), sigs, tx.Tx[4:], []byte(strings.Join(result, "")))
}
