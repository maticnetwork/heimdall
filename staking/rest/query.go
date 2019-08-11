package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Get all delegations from a delegator
	r.HandleFunc(
		"/staking/signer/{address}",
		validatorByAddressHandlerFn(cdc, cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/staking/validator/{id}",
		validatorByIDHandlerFn(cdc, cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/staking/validator-set",
		validatorSetHandlerFn(cdc, cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/staking/proposer/{times}",
		proposerHandlerFn(cdc, cliCtx),
	).Methods("GET")
}

// Returns validator information by signer address
func validatorByAddressHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		signerAddress := common.HexToAddress(vars["address"])

		res, err := cliCtx.QueryStore(staking.GetValidatorKey(signerAddress.Bytes()), "staking")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// the query will return empty if there is no data
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNoContent, err.Error())
			return
		}

		var _validator types.Validator
		err = cdc.UnmarshalBinaryBare(res, &_validator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := json.Marshal(_validator)
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

// Returns validator information by val ID
func validatorByIDHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// get id
		id, ok := rest.ParseUint64OrReturnBadRequest(w, vars["id"])
		if !ok {
			return
		}

		signerAddr, err := cliCtx.QueryStore(staking.GetValidatorMapKey(hmTypes.NewValidatorID(id).Bytes()), "staking")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryStore(staking.GetValidatorKey(signerAddr), "staking")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// the query will return empty if there is no data
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNoContent, err.Error())
			return
		}

		var _validator types.Validator
		err = cdc.UnmarshalBinaryBare(res, &_validator)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := json.Marshal(_validator)
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

// get current validator set
func validatorSetHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryStore(staking.CurrentValidatorSetKey, "staking")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		// the query will return empty if there is no data
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNoContent, err.Error())
			return
		}
		var _validatorSet hmTypes.ValidatorSet
		cdc.UnmarshalBinaryBare(res, &_validatorSet)

		// todo format validator set to remove pubkey like we did for validator
		result, err := json.Marshal(&_validatorSet)
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

// get proposer for current validator set
func proposerHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		times, err := strconv.Atoi(vars["times"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		RestLogger.Debug("Calculating proposers", "Count", times)

		res, err := cliCtx.QueryStore(staking.CurrentValidatorSetKey, "staking")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// the query will return empty if there is no data
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNoContent, err.Error())
			return
		}

		var _validatorSet hmTypes.ValidatorSet
		err = cdc.UnmarshalBinaryBare(res, &_validatorSet)
		if err != nil {
			RestLogger.Error("Error while marshalling validator set", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if times > len(_validatorSet.Validators) {
			times = len(_validatorSet.Validators)
		}

		// proposers
		var proposers []hmTypes.Validator

		for index := 0; index < times; index++ {
			RestLogger.Info("Getting proposer for current validator set", "Index", index, "TotalProposers", times)
			proposers = append(proposers, *(_validatorSet.GetProposer()))
			_validatorSet.IncrementAccum(1)
		}

		// TODO format validator set to remove pubkey like we did for validator
		result, err := json.Marshal(&proposers)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}
