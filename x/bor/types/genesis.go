package types

import (
	"encoding/json"

	hmTypes "github.com/maticnetwork/heimdall/types"
	chainManagerTypes "github.com/maticnetwork/heimdall/x/chainmanager/types"
	"github.com/maticnetwork/heimdall/x/gov/types"
)

// NewGenesisState creates a new genesis state.
func NewGenesisState(params Params, spans []*hmTypes.Span) *GenesisState {
	return &GenesisState{
		Params: &params,
		Spans:  spans,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(), nil)
}

// ValidateGenesis performs basic validation of bor genesis data returning an
// error for any failed validation criteria.
func (gs GenesisState) ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	return nil
}

// genFirstSpan generates default first validator producer set
func genFirstSpan(valSet hmTypes.ValidatorSet, chainId string) []*hmTypes.Span {
	var firstSpan []*hmTypes.Span
	var selectedProducers []hmTypes.Validator
	validators := valSet.GetValidatorsSe()
	if len(valSet.Validators) > int(DefaultProducerCount) {
		selectedProducers = append(selectedProducers, validators[0:DefaultProducerCount]...)
	} else {
		selectedProducers = append(selectedProducers, validators...)
	}

	newSpan := hmTypes.NewSpan(0, 0, 0+DefaultFirstSpanDuration-1, valSet, selectedProducers, chainId)
	firstSpan = append(firstSpan, &newSpan)
	return firstSpan
}

// GetGenesisStateFromAppState returns staking GenesisState given raw application genesis state
func GetGenesisStateFromAppState(appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		types.ModuleCdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}

// SetGenesisStateToAppState sets state into app state
func SetGenesisStateToAppState(appState map[string]json.RawMessage, currentValSet hmTypes.ValidatorSet) (map[string]json.RawMessage, error) {
	// set state to bor state
	borState := GetGenesisStateFromAppState(appState)
	chainState := chainManagerTypes.GetGenesisStateFromAppState(appState)
	borState.Spans = genFirstSpan(currentValSet, chainState.Params.ChainParams.BorChainID)

	appState[ModuleName] = types.ModuleCdc.MustMarshalJSON(&borState)
	return appState, nil
}
