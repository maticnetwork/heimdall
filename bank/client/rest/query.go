package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/bank"
	"github.com/maticnetwork/heimdall/bank/types"
	bankTypes "github.com/maticnetwork/heimdall/bank/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

// QueryBalancesRequestHandlerFn query accountREST Handler
func QueryBalancesRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		addr := hmTypes.HexToHeimdallAddress(vars["address"])

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryBalanceParams(addr)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", bankTypes.QuerierRoute, bank.QueryBalance), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// the query will return empty if there is no data for this account
		if len(res) == 0 {
			rest.PostProcessResponse(w, cliCtx, sdk.Coins{})
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
