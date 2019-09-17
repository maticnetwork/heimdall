package app

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/maticnetwork/heimdall/auth"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/bank"
	"github.com/maticnetwork/heimdall/bor"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/clerk"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/supply"
	"github.com/maticnetwork/heimdall/types"
)

// GenesisAccount genesis account
type GenesisAccount struct {
	Address       types.HeimdallAddress `json:"address"`
	Coins         types.Coins           `json:"coins"`
	Sequence      uint64                `json:"sequence_number"`
	AccountNumber uint64                `json:"account_number"`
}

// NewGenesisAccount creates new genesis account
func NewGenesisAccount(acc authTypes.Account) GenesisAccount {
	gacc := GenesisAccount{
		Address:       acc.GetAddress(),
		Coins:         acc.GetCoins(),
		AccountNumber: acc.GetAccountNumber(),
		Sequence:      acc.GetSequence(),
	}

	return gacc
}

func BaseToGenesisAcc(acc authTypes.BaseAccount) GenesisAccount {
	return GenesisAccount{
		Address:       acc.Address,
		Coins:         acc.Coins,
		Sequence:      acc.Sequence,
		AccountNumber: acc.AccountNumber,
	}
}

// GenesisState to Unmarshal
type GenesisState struct {
	Accounts []GenesisAccount  `json:"accounts"`
	GenTxs   []json.RawMessage `json:"gentxs"`

	AuthData   auth.GenesisState   `json:"auth"`
	BankData   bank.GenesisState   `json:"bank"`
	SupplyData supply.GenesisState `json:"supply"`
	// GovData  gov.GenesisState  `json:"gov"`

	BorData        bor.GenesisState        `json:"bor"`
	CheckpointData checkpoint.GenesisState `json:"checkpoint"`
	StakingData    staking.GenesisState    `json:"staking"`
	ClerkData      clerk.GenesisState      `json:"clerk"`
}

// NewGenesisState creates new genesis state
func NewGenesisState(
	accounts []GenesisAccount,

	authData auth.GenesisState,
	bankData bank.GenesisState,
	supplyData supply.GenesisState,
	// govData gov.GenesisState,

	borData bor.GenesisState,
	checkpointData checkpoint.GenesisState,
	stakingData staking.GenesisState,
	clerkData clerk.GenesisState,
) GenesisState {
	return GenesisState{
		Accounts: accounts,

		AuthData:   authData,
		BankData:   bankData,
		SupplyData: supplyData,
		// GovData:  govData,

		BorData:        borData,
		CheckpointData: checkpointData,
		StakingData:    stakingData,
		ClerkData:      clerkData,
	}
}

// Sanitize sorts accounts and coin sets.
func (gs GenesisState) Sanitize() {
	sort.Slice(gs.Accounts, func(i, j int) bool {
		return gs.Accounts[i].AccountNumber < gs.Accounts[j].AccountNumber
	})

	for _, acc := range gs.Accounts {
		acc.Coins = acc.Coins.Sort()
	}
}

// ValidateGenesisState ensures that the genesis state obeys the expected invariants
func ValidateGenesisState(genesisState GenesisState) error {
	if err := validateGenesisStateAccounts(genesisState.Accounts); err != nil {
		return err
	}

	// skip stakingData validation as genesis is created from txs
	if len(genesisState.GenTxs) > 0 {
		return nil
	}

	if err := auth.ValidateGenesis(genesisState.AuthData); err != nil {
		return err
	}
	if err := bank.ValidateGenesis(genesisState.BankData); err != nil {
		return err
	}
	if err := supply.ValidateGenesis(genesisState.SupplyData); err != nil {
		return err
	}
	if err := staking.ValidateGenesis(genesisState.StakingData); err != nil {
		return err
	}
	if err := bor.ValidateGenesis(genesisState.BorData); err != nil {
		return err
	}
	if err := checkpoint.ValidateGenesis(genesisState.CheckpointData); err != nil {
		return err
	}
	if err := clerk.ValidateGenesis(genesisState.ClerkData); err != nil {
		return err
	}

	return nil
}

// validateGenesisStateAccounts performs validation of genesis accounts. It
// ensures that there are no duplicate accounts in the genesis state and any
// provided vesting accounts are valid.
func validateGenesisStateAccounts(accs []GenesisAccount) error {
	addrMap := make(map[string]bool, len(accs))
	for _, acc := range accs {
		addrStr := acc.Address.String()

		// disallow any duplicate accounts
		if _, ok := addrMap[addrStr]; ok {
			return fmt.Errorf("duplicate account found in genesis state; address: %s", addrStr)
		}

		addrMap[addrStr] = true
	}

	return nil
}

// NewDefaultGenesisState generates the default state for gaia.
func NewDefaultGenesisState() GenesisState {
	return GenesisState{
		Accounts:       nil,
		AuthData:       auth.DefaultGenesisState(),
		BankData:       bank.DefaultGenesisState(),
		StakingData:    staking.DefaultGenesisState(),
		CheckpointData: checkpoint.DefaultGenesisState(),
		BorData:        bor.DefaultGenesisState(staking.DefaultGenesisState().CurrentValSet),
		ClerkData:      clerk.DefaultGenesisState(),
		GenTxs:         nil,
	}
}
