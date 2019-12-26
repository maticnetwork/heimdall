package types

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// HeimdallModuleBasic is the standard form for basic non-dependant elements of an application module.
type HeimdallModuleBasic interface {
	module.AppModuleBasic

	// verify genesis
	VerifyGenesis(map[string]json.RawMessage) error
}
