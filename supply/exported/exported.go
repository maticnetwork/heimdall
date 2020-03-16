package exported

import (
	"github.com/maticnetwork/heimdall/auth/exported"
)

// ModuleAccountI defines an account interface for modules that hold tokens in
// an escrow.
type ModuleAccountI interface {
	exported.Account

	GetName() string
	GetPermissions() []string
	HasPermission(string) bool
}
