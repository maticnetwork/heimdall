package pier

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/common"
	httpClient "github.com/tendermint/tendermint/rpc/client"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/bor"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

const (
	lastSpanKey = "span-key" // storage key

	// polling
	spanPolling = 20 * time.Second
)

// SpanService service spans
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
	queueConnector *QueueConnector

	// http client to subscribe to
	httpClient *httpClient.HTTP
}

// NewSpanService returns new service object
func NewSpanService(cdc *codec.Codec, queueConnector *QueueConnector, httpClient *httpClient.HTTP) *SpanService {
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
	cliCtx.TrustNode = true

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
	go s.startPolling(spanCtx, spanPolling)

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
			if lastSpan, err := s.getLastSpan(); err == nil && lastSpan != nil {
				if s.isSpanProposer(lastSpan) {
					go s.propose(lastSpan)
				}
			}
			go s.commit()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

// propose producers for next span if needed
func (s *SpanService) propose(lastSpan *hmTypes.Span) {
	// call with last span on record + new span duration and see if it has been proposed
	currentBlock, err := s.getCurrentChildBlock()
	if err != nil {
		s.Logger.Error("Unable to fetch current block", "error", err)
		return
	}

	if lastSpan.StartBlock <= currentBlock && currentBlock <= lastSpan.EndBlock {
		// send propose span
		msg, err := s.fetchNextSpanDetails(lastSpan.ID+1, lastSpan.EndBlock+1)
		if err != nil {
			s.Logger.Error("Unable to fetch next span details", "error", err)
			return
		}

		// log new span
		s.Logger.Info("Proposing new span", "spanId", msg.ID, "startBlock", msg.StartBlock, "endBlock", msg.EndBlock)

		// broadcast to heimdall
		if err := s.queueConnector.BroadcastToHeimdall(msg); err != nil {
			s.Logger.Error("Error while broadcasting msg to heimdall", "error", err)
			return
		}
	}
}

func (s *SpanService) commit() {
	// get current span number from bor chain
	currentSpanNumber := s.contractConnector.CurrentSpanNumber()
	if currentSpanNumber == nil {
		currentSpanNumber = big.NewInt(0)
	}

	// create tag query
	var tags []string
	tags = append(tags, fmt.Sprintf("bor-sync-id>%v", currentSpanNumber))
	tags = append(tags, "action='propose-span'")

	s.Logger.Debug("[COMMIT SPAN] Querying heimdall span txs",
		"currentSpanNumber", currentSpanNumber,
		"tags", strings.Join(tags, " AND "),
	)

	// search txs
	txs, err := helper.SearchTxs(s.cliCtx, s.cliCtx.Codec, tags, 1, 20) // first page, 50 limit
	if err != nil {
		s.Logger.Error("Error while searching txs", "error", err)
		return
	}

	s.Logger.Debug("[COMMIT SPAN] Found new span txs",
		"length", len(txs),
	)

	// loop through tx
	for _, tx := range txs {
		txHash, err := hex.DecodeString(tx.TxHash)
		if err != nil {
			s.Logger.Error("Error while searching txs", "error", err)
		} else {
			s.broadcastToBor(tx.Height, txHash)
		}
	}
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
func (s *SpanService) getLastSpan() (*hmTypes.Span, error) {
	// fetch latest start block from heimdall via rest query
	result, err := FetchFromAPI(s.cliCtx, GetHeimdallServerEndpoint(LatestSpanURL))
	if err != nil {
		s.Logger.Error("Error while fetching latest span")
		return nil, err
	}

	var lastSpan hmTypes.Span
	err = json.Unmarshal(result.Result, &lastSpan)
	if err != nil {
		s.Logger.Error("Error unmarshalling", "error", err)
		return nil, err
	}

	return &lastSpan, nil
}

// getCurrentChildBlock gets the current child block
func (s *SpanService) getCurrentChildBlock() (uint64, error) {
	childBlock, err := s.contractConnector.GetMaticChainBlock(nil)
	if err != nil {
		return 0, err
	}
	return childBlock.Number.Uint64(), nil
}

// isSpanProposer check if current user is proposer
func (s *SpanService) isSpanProposer(lastSpan *hmTypes.Span) bool {
	// sort validator address
	selectedProducers := types.SortValidatorByAddress(lastSpan.SelectedProducers)

	// get last validator as proposer
	proposer := selectedProducers[len(selectedProducers)-1]

	// check proposer
	return bytes.Equal(proposer.Signer.Bytes(), helper.GetAddress())
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

	// fetch next span details
	result, err := FetchFromAPI(s.cliCtx, req.URL.String())
	if err != nil {
		s.Logger.Error("Error fetching proposers", "error", err)
		return
	}

	err = json.Unmarshal(result.Result, &msg)
	if err != nil {
		s.Logger.Error("Error unmarshalling propose tx msg ", "error", err)
		return
	}
	s.Logger.Info("==>Generated proposer span msg", "msg", msg)
	return msg, nil
}

// // SubscribeToTx subscribes to a broadcasted Tx and waits for its commitment to a block
// func (s *SpanService) SubscribeToTx(tx tmTypes.Tx, start, end uint64) error {
// 	data, err := WaitForOneEvent(tx, s.httpClient)
// 	if err != nil {
// 		s.Logger.Error("Unable to wait for tx", "error", err)
// 		return err
// 	}

// 	switch t := data.(type) {
// 	case tmTypes.EventDataTx:
// 		go s.DispatchProposal(t.Height, t.Tx.Hash(), tx)
// 	default:
// 		s.Logger.Info("No cases matched while trying to send propose new committee")
// 	}
// 	return nil
// }

// broadcastToBor broadcasts to bor
func (s *SpanService) broadcastToBor(height int64, txHash []byte) error {
	// extraData
	votes, sigs, chainID, err := FetchVotes(height, s.httpClient)
	if err != nil {
		s.Logger.Error("Error fetching votes", "height", height)
		return err
	}

	// proof
	tx, err := helper.QueryTxWithProof(s.cliCtx, txHash)
	if err != nil {
		return err
	}
	// fmt.Println("TxBytes: ", hex.EncodeToString(tx.Tx[4:]))
	// fmt.Println("Leaf: ", hex.EncodeToString(tx.Proof.Leaf()))
	// fmt.Println("Root: ", tx.Proof.RootHash.String())
	proofList := helper.GetMerkleProofList(&tx.Proof.Proof)
	proof := helper.AppendBytes(proofList...)

	// encode commit span
	encodedData := s.encodeCommitSpanData(
		helper.GetVoteBytes(votes, chainID),
		sigs,
		tx.Tx[authTypes.PulpHashLength:],
		proof,
	)

	// fmt.Println("data : ",
	// 	fmt.Sprintf(`"0x%s","0x%s","0x%s","0x%s"`,
	// 		hex.EncodeToString(helper.GetVoteBytes(votes, chainID)),
	// 		hex.EncodeToString(sigs),
	// 		hex.EncodeToString(tx.Tx[4:]),
	// 		hex.EncodeToString(proof),
	// 	))

	// get validator address
	validatorSetAddress := helper.GetValidatorSetAddress()
	msg := ethereum.CallMsg{
		To:   &validatorSetAddress,
		Data: encodedData,
	}

	// encode msg data
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// broadcast to bor queue
	if err := s.queueConnector.BroadcastToBor(data); err != nil {
		s.Logger.Error("Error while dispatching to bor queue", "error", err)
		return err
	}

	return nil
}

//
// ABI encoding
//

func (s *SpanService) encodeCommitSpanData(voteSignBytes []byte, sigs []byte, txData []byte, proof []byte) []byte {
	// validator set ABI
	validatorSetABI := s.contractConnector.ValidatorSetABI
	// commit span
	data, err := validatorSetABI.Pack("commitSpan", voteSignBytes, sigs, txData, proof)
	if err != nil {
		Logger.Error("Unable to pack tx for commit span", "error", err)
		return nil
	}

	// return data
	return data
}
