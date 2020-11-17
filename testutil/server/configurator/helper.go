package configurator

import (
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MakeTestAddresses makes n test addresses for using in configurator.Fixture tests.
func MakeTestAddresses(count int) []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, count)
	for i := 0; i < count; i++ {
		_, _, addrs[i] = testdata.KeyTestPubAddr()
	}
	return addrs
}
