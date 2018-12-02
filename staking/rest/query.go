package rest

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	hmcommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/types"
	"net/http"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Get all delegations from a delegator
	r.HandleFunc(
		"/staking/validator/{address}",
		ValidatorByAddressHandlerFn(cdc, cliCtx),
	).Methods("GET")

}

// Returns validator information by address
func ValidatorByAddressHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		validatorAddress := common.HexToAddress(vars["address"])

		res, err := cliCtx.QueryStore(hmcommon.GetValidatorKey(validatorAddress.Bytes()), "staker")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// the query will return empty if there is no data
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var _validator types.Validator
		err = json.Unmarshal(res, &_validator)
		if err != nil {
			RestLogger.Info("Unable to marshall validator")
		}

		result, err := json.Marshal(&_validator)
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)

	}
}
