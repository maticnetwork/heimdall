package auth

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	gasWantedPerTx sdk.Gas = 900000
	gasUsedPerTx   sdk.Gas = gasWantedPerTx - 60000 // TODO use proposer amount per tx
)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler() sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, simulate bool,
	) (newCtx sdk.Context, res sdk.Result, abort bool) {
		// AnteHandlers must have their own defer/recover in order for the BaseApp
		// to know how much gas was used! This is because the GasMeter is created in
		// the AnteHandler, but if it panics the context won't be set properly in
		// runTx's recover call.
		defer func() {
			if r := recover(); r != nil {
				switch rType := r.(type) {
				case sdk.ErrorOutOfGas:
					log := fmt.Sprintf("out of gas in location: %v", rType.Descriptor)
					res = sdk.ErrOutOfGas(log).Result()
					res.GasWanted = gasWantedPerTx
					res.GasUsed = newCtx.GasMeter().GasConsumed()
					abort = true
				default:
					panic(r)
				}
			}
		}()

		// get new gas meter context
		newCtx = getGasMeterContext(simulate, ctx)
		newCtx.GasMeter().ConsumeGas(gasUsedPerTx, "fee")

		return newCtx, sdk.Result{GasWanted: gasWantedPerTx, GasUsed: gasUsedPerTx}, false
	}
}

// getGasMeterContext returns a new context with a gas meter set from a given context.
func getGasMeterContext(simulate bool, ctx sdk.Context) sdk.Context {
	// In various cases such as simulation and during the genesis block, we do not
	// meter any gas utilization.
	if simulate || ctx.BlockHeight() == 0 {
		return ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	}

	return ctx.WithGasMeter(sdk.NewGasMeter(gasWantedPerTx))
}
