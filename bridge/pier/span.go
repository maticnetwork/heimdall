package pier

import (
	"context"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/maticnetwork/heimdall/bor"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/common"
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

	cliCtx cliContext.CLIContext
}

// NewAckService returns new service object
func NewSpanService() *SpanService {
	// create logger
	logger := Logger.With("module", SpanServiceStr)

	// root chain instance
	rootchainInstance, err := helper.GetRootChainInstance()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	cliCtx := cliContext.NewCLIContext()
	cliCtx.BroadcastMode = client.BroadcastAsync

	// creating checkpointer object
	spanService := &SpanService{
		storageClient:     getBridgeDBInstance(viper.GetString(BridgeDBFlag)),
		rootChainInstance: rootchainInstance,
		cliCtx:            cliCtx,
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
	go s.startPolling(spanCtx, time.Duration(bor.DefaultSpanDuration/2))

	// subscribed to new head
	s.Logger.Debug("Started Span service")
	return nil
}

// OnStop stops all necessary go routines
func (s *SpanService) OnStop() {
	s.BaseService.OnStop()

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
			go s.propose(ctx)
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

// propose producers for next span if needed
func (s *SpanService) propose(ctx context.Context) {
	// call with last span on record + new span duration and see if it has been proposed
	lastProcessedSpan, err := s.fetchLastSpan()
	if err != nil {
		s.Logger.Error("Unable to fetch last processed span from storage", "error", err)
	}

	// if no, send propose span

	// if yes, check if propose span has been reflected in bor

	// if no send propose span to bor

	// if yes bbye
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

func (s *SpanService) checkSpanStatus() {

}
