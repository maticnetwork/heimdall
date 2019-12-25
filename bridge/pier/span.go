package pier

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/common"
	httpClient "github.com/tendermint/tendermint/rpc/client"

	"github.com/maticnetwork/heimdall/bor"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
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
			s.checkAndPropose()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

// checkAndPropose will check if current user is span proposer and proposes the span
func (s *SpanService) checkAndPropose() {
	lastSpan, err := s.getLastSpan()
	if err == nil && lastSpan != nil {
		nextSpanMsg, err := s.fetchNextSpanDetails(lastSpan.ID+1, lastSpan.EndBlock+1)

		// check if current user is among next span producers
		if err == nil && s.isSpanProposer(nextSpanMsg.SelectedProducers) {
			go s.propose(lastSpan, nextSpanMsg)
		}
	}
}

// propose producers for next span if needed
func (s *SpanService) propose(lastSpan *types.Span, nextSpanMsg *types.Span) {
	// call with last span on record + new span duration and see if it has been proposed
	currentBlock, err := s.getCurrentChildBlock()
	if err != nil {
		s.Logger.Error("Unable to fetch current block", "error", err)
		return
	}

	if lastSpan.StartBlock <= currentBlock && currentBlock <= lastSpan.EndBlock {
		// log new span
		s.Logger.Info("Proposing new span", "spanId", nextSpanMsg.ID, "startBlock", nextSpanMsg.StartBlock, "endBlock", nextSpanMsg.EndBlock)

		// broadcast to heimdall
		msg := bor.MsgProposeSpan{
			ID:         nextSpanMsg.ID,
			Proposer:   types.BytesToHeimdallAddress(helper.GetAddress()),
			StartBlock: nextSpanMsg.StartBlock,
			EndBlock:   nextSpanMsg.EndBlock,
			ChainID:    nextSpanMsg.ChainID,
		}
		if err := s.queueConnector.BroadcastToHeimdall(msg); err != nil {
			s.Logger.Error("Error while broadcasting msg to heimdall", "error", err)
			return
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
func (s *SpanService) getLastSpan() (*types.Span, error) {
	// fetch latest start block from heimdall via rest query
	result, err := FetchFromAPI(s.cliCtx, GetHeimdallServerEndpoint(LatestSpanURL))
	if err != nil {
		s.Logger.Error("Error while fetching latest span")
		return nil, err
	}

	var lastSpan types.Span
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

// isSpanProposer checks if current user is span proposer
func (s *SpanService) isSpanProposer(nextSpanProducers []types.Validator) bool {
	// anyone among next span producers can become next span proposer
	for _, val := range nextSpanProducers {
		if bytes.Equal(val.Signer.Bytes(), helper.GetAddress()) {
			return true
		}
	}
	return false
}

func (s *SpanService) fetchNextSpanDetails(id uint64, start uint64) (*types.Span, error) {
	req, err := http.NewRequest("GET", GetHeimdallServerEndpoint(NextSpanInfoURL), nil)
	if err != nil {
		s.Logger.Error("Error creating a new request", "error", err)
		return nil, err
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
		return nil, err
	}

	var msg types.Span
	if err = json.Unmarshal(result.Result, &msg); err != nil {
		s.Logger.Error("Error unmarshalling propose tx msg ", "error", err)
		return nil, err
	}

	s.Logger.Debug("Generated proposer span msg", "msg", msg)
	return &msg, nil
}
