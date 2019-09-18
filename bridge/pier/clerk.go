package pier

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/common"
	httpClient "github.com/tendermint/tendermint/rpc/client"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// storage key
var lastEventRecordKey = []byte("clerk-event-record-key")

const (
	// polling
	clerkPolling = 20 * time.Second
)

// ClerkService service spans
type ClerkService struct {
	// Base service
	common.BaseService

	// storage client
	storageClient *leveldb.DB

	// header listener subscription
	cancel context.CancelFunc

	// contract caller
	contractConnector helper.ContractCaller

	// cli context
	cliCtx cliContext.CLIContext

	// queue connector
	queueConnector *QueueConnector

	// http client to subscribe to
	httpClient *httpClient.HTTP
}

// NewClerkService returns new service object
func NewClerkService(cdc *codec.Codec, queueConnector *QueueConnector, httpClient *httpClient.HTTP) *ClerkService {
	// create logger
	logger := Logger.With("module", ClerkServiceStr)

	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastAsync
	cliCtx.TrustNode = true

	// creating clerk service
	clerkService := &ClerkService{
		storageClient:     getBridgeDBInstance(viper.GetString(BridgeDBFlag)),
		contractConnector: contractCaller,

		cliCtx:         cliCtx,
		queueConnector: queueConnector,
		httpClient:     httpClient,
	}

	clerkService.BaseService = *common.NewBaseService(logger, ClerkServiceStr, clerkService)
	return clerkService
}

// OnStart starts new block subscription
func (s *ClerkService) OnStart() error {
	s.BaseService.OnStart() // Always call the overridden method.

	// create cancellable context
	clerkCtx, cancel := context.WithCancel(context.Background())

	s.cancel = cancel

	// start polling for checkpoint in buffer
	go s.startPolling(clerkCtx, clerkPolling)

	// subscribed to new head
	s.Logger.Debug("Started Span service")
	return nil
}

// OnStop stops all necessary go routines
func (s *ClerkService) OnStop() {
	s.BaseService.OnStop()
	s.httpClient.Stop()

	// cancel ack process
	s.cancel()
	// close bridge db instance
	closeBridgeDBInstance()
}

// polls heimdall and checks if new span needs to be proposed
func (s *ClerkService) startPolling(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	// stop ticker when everything done
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go s.commit()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (s *ClerkService) commit() {
	// get current span number from bor chain
	currentStateCounter := s.contractConnector.CurrentStateCounter()
	if currentStateCounter == nil {
		currentStateCounter = big.NewInt(0)
	}

	// get current storage
	lastEventRecord, _ := s.fetchLastEventRecordID()

	// start from
	start := lastEventRecord

	// if diff >= 250, ignore
	if currentStateCounter.Uint64() > start && currentStateCounter.Uint64()-start >= 250 {
		start = currentStateCounter.Uint64() - 1
	}

	// create tag query
	var tags []string
	tags = append(tags, fmt.Sprintf("record-id>%v", start))
	tags = append(tags, "action='event-record'")

	s.Logger.Debug("[COMMIT RECORD] Querying heimdall event record txs",
		"start", start,
		"lastEventRecord", lastEventRecord,
		"currentStateCounter", currentStateCounter.Uint64(),
		"tags", strings.Join(tags, " AND "),
	)

	// search txs
	txs, err := helper.SearchTxs(s.cliCtx, s.cliCtx.Codec, tags, 1, 50) // first page, 50 limit
	if err != nil {
		s.Logger.Error("Error while searching txs", "error", err)
		return
	}

	s.Logger.Debug("[COMMIT RECORD] Found new state txs",
		"length", len(txs),
	)

	// loop through tx
	end := start
	for _, tx := range txs {
		txHash, err := hex.DecodeString(tx.TxHash)
		if err != nil {
			s.Logger.Error("Error while searching txs", "error", err)
		} else {
			s.broadcastToBor(tx.Height, txHash)
		}

		for _, tag := range tx.Tags {
			if tag.Key == clerkTypes.RecordID {
				if e, err := strconv.ParseUint(tag.Value, 10, 64); err == nil && e > end {
					end = e
				}
			}
		}
	}

	// save last record id
	if end != start {
		s.saveLastEventRecordID(end)
	}
}

// fetches last event record processed in DB
func (s *ClerkService) fetchLastEventRecordID() (uint64, error) {
	hasLastID, _ := s.storageClient.Has(lastEventRecordKey, nil)
	if hasLastID {
		lastLastIDBytes, err := s.storageClient.Get(lastEventRecordKey, nil)
		if err != nil {
			s.Logger.Info("Error while fetching last span bytes from storage", "error", err)
			return 0, err
		}

		s.Logger.Debug("Got last block from bridge storage", "lastSpan", string(lastLastIDBytes))
		result, err := strconv.ParseUint(string(lastLastIDBytes), 10, 64)
		if err != nil {
			return 0, nil
		}

		return result, nil
	}
	return 0, errors.New("No last id found")
}

func (s *ClerkService) saveLastEventRecordID(result uint64) {
	// set last block to storage
	s.storageClient.Put(lastEventRecordKey, []byte(strconv.FormatUint(result, 10)), nil)
}

// checks state counter
func (s *ClerkService) getStateSyncerCounter() (*hmTypes.Span, error) {
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

// isRecordProposer check if current user is proposer
func (s *ClerkService) isRecordProposer(lastSpan *hmTypes.Span) bool {
	// sort validator address
	selectedProducers := types.SortValidatorByAddress(lastSpan.SelectedProducers)

	// get last validator as proposer
	proposer := selectedProducers[len(selectedProducers)-1]

	// check proposer
	return bytes.Equal(proposer.Signer.Bytes(), helper.GetAddress())
}

// broadcastToBor broadcasts to bor
func (s *ClerkService) broadcastToBor(height int64, txHash []byte) error {
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
	encodedData := s.encodeCommitStateData(
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
	stateReceiverAddress := helper.GetStateReceiverAddress()
	msg := ethereum.CallMsg{
		To:   &stateReceiverAddress,
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

func (s *ClerkService) encodeCommitStateData(voteSignBytes []byte, sigs []byte, txData []byte, proof []byte) []byte {
	// state receiver ABI
	stateReceiverABI := s.contractConnector.StateReceiverABI

	// commit state
	data, err := stateReceiverABI.Pack("commitState", voteSignBytes, sigs, txData, proof)
	if err != nil {
		Logger.Error("Unable to pack tx for commit state", "error", err)
		return nil
	}

	// return data
	return data
}
