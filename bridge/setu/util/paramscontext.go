package util

import (
	"time"

	chainmanagerTypes "github.com/maticnetwork/heimdall/x/chainmanager/types"

	checkpointtypes "github.com/maticnetwork/heimdall/x/checkpoint/types"

	"github.com/patrickmn/go-cache"

	"github.com/cosmos/cosmos-sdk/client"
)

type ParamsContext struct {
	cliCtx      client.Context
	paramsCache *cache.Cache
	key         string
}

type Params struct {
	ChainmanagerParams *chainmanagerTypes.Params
	CheckpointParams   *checkpointtypes.Params
}

// NewParamsContext creates new params context
func NewParamsContext(cliCtx client.Context) *ParamsContext {

	paramsContext := ParamsContext{
		key:         "params",
		cliCtx:      cliCtx,
		paramsCache: cache.New(1*time.Hour, 1*time.Hour),
	}

	return &paramsContext
}

// GetParams updates cache if required and returns params
func (paramsContext *ParamsContext) GetParams() (params Params, err error) {
	var found bool
	data, found := paramsContext.paramsCache.Get(paramsContext.key)
	if found {
		params = data.(Params)
	} else {
		// Fetch params and add to cache
		params, err = fetchLatestParams(paramsContext.cliCtx)
		if err == nil {
			paramsContext.paramsCache.Set(paramsContext.key, params, 1*time.Hour)
		}
	}
	return
}

func fetchLatestParams(cliContext client.Context) (params Params, err error) {
	chainmanagerParams, err := GetChainmanagerParams(cliContext)
	if err != nil {
		return
	}

	checkpointParams, err := GetCheckpointParams(cliContext)
	if err != nil {
		return
	}

	params = Params{
		ChainmanagerParams: chainmanagerParams,
		CheckpointParams:   checkpointParams,
	}
	return
}
