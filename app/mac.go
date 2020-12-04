package app

import (
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func MacPerms() map[string][]string {
	return map[string][]string{
		authtypes.FeeCollectorName: nil,
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
