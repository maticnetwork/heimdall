package app

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/maticnetwork/heimdall/x/gov/types"
)

func MacPerms() map[string][]string {
	return map[string][]string{
		authtypes.FeeCollectorName: nil,
		govtypes.ModuleName:        {authtypes.Burner},
	}
}

func MacAddrs() map[string]bool {
	perms := MacPerms()
	addrs := make(map[string]bool, len(perms))
	for k := range perms {
		addrs[authtypes.NewModuleAddress(k).String()] = true
	}
	return addrs
}
