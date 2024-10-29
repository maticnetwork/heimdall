package module

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/maticnetwork/heimdall/types"
)

// HeimdallModuleBasic is the standard form for basic non-dependant elements of an application module.
type HeimdallModuleBasic interface {
	module.AppModuleBasic

	// verify genesis
	VerifyGenesis(map[string]json.RawMessage) error
}

// SideModule is the standard form for side tx elements of an application module
type SideModule interface {
	NewSideTxHandler() types.SideTxHandler
	NewPostTxHandler() types.PostTxHandler
}

type ModuleGenesisData struct {
	// Path specifies the JSON path where the data should be appended.
	// For example, "moduleA.data" refers to the "data" array within "moduleA".
	Path string

	// Data is the JSON data chunk to be appended.
	Data json.RawMessage

	// NextKey is the last key used to append data.
	NextKey []byte
}

// StreamedGenesisExporter defines an interface for modules to export their genesis data incrementally.
type StreamedGenesisExporter interface {
	// ExportPartialGenesis returns the partial genesis state of the module,
	// the rest of the genesis data will be exported in subsequent calls to NextGenesisData.
	ExportPartialGenesis(ctx sdk.Context) (json.RawMessage, error)
	// NextGenesisData returns the next chunk of genesis data.
	// Returns nil NextKey when no more data is available.
	NextGenesisData(ctx sdk.Context, nextKey []byte, max int) (*ModuleGenesisData, error)
}
