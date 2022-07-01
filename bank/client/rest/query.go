package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	bankTypes "github.com/maticnetwork/heimdall/bank/types"
	"github.com/maticnetwork/heimdall/types"
)

//It represents the bank balance of particluar account
//swagger:response bankBalanceByAddressResponse
type bankBalanceByAddressResponse struct {
	//in:body
	Output bankBalanceByAddress `json:"output"`
}

type bankBalanceByAddress struct {
	Height string        `json:"height"`
	Result []bankBalance `json:"result"`
}

type bankBalance struct {

	//Denomination of the token
	Denom string `json:"denom"`
	//Amount of token in the bank
	Amount string `json:"amount"`
}

//swagger:parameters bankBalanceByAddress
type borSpanListParam struct {

	//Address of the account
	//required:true
	//in:path
	Address string `json:"address"`

	//Address of the account
	//in:query
	Height string `json:"height"`
}

// swagger:route GET /bank/balances/{address} bank bankBalanceByAddress
// It returns the matic balance of particular address
// responses:
//   200: bankBalanceByAddressResponse
// QueryBalancesRequestHandlerFn query accountREST Handler
func QueryBalancesRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		addr := types.HexToHeimdallAddress(vars["address"])

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := bankTypes.NewQueryBalanceParams(addr)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", bankTypes.QuerierRoute, bankTypes.QueryBalance), bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		// the query will return empty if there is no data for this account
		if len(res) == 0 {
			rest.PostProcessResponse(w, cliCtx, sdk.Coins{})
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
