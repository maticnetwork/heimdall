package util

import (
	"time"

	"github.com/patrickmn/go-cache"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	chainmanagerTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
)

type ParamsContext struct {
	cliCtx      cliContext.CLIContext
	paramsCache *cache.Cache
	key         string
}

type Params struct {
	ChainmanagerParams *chainmanagerTypes.Params
	CheckpointParams   *checkpointTypes.Params
}

// NewParamsContext creates new params context
func NewParamsContext(cliCtx cliContext.CLIContext) *ParamsContext {

	paramsContext := ParamsContext{
		key:         "params",
		cliCtx:      cliCtx,
		paramsCache: cache.New(1*time.Hour, 1*time.Hour),
	}

	return &paramsContext
}

// GetParamsContext gets ParamsContext
func (paramsContext *ParamsContext) GetParams() (params Params, err error) {
	var found bool

	data, found := paramsContext.paramsCache.Get(paramsContext.key)
	if found {
		params = data.(Params)
	} else {
		// Fetch params and add to cache
		params, err := fetchLatestParams(paramsContext.cliCtx)
		if err != nil {
			paramsContext.paramsCache.Set(paramsContext.key, params, 1*time.Hour)
		}
	}
	return
}

func fetchLatestParams(cliContext cliContext.CLIContext) (params Params, err error) {
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
