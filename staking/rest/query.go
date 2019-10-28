package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/checkpoint"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/staking"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
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
		"/staking/validator-account/{id}",
		validatorAccountByIDHandlerFn(cdc, cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/staking/proposer/{times}",
		proposerHandlerFn(cdc, cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/staking/current-proposer",
		currentProposerHandlerFn(cdc, cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/staking/slash-validator",
		slashValidatorHandlerFn(cdc, cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/staking/initial-account-root",
		initialAccountRootHandlerFn(cdc, cliCtx),
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
			rest.WriteErrorResponse(w, http.StatusNoContent, errors.New("no content found for requested key").Error())
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
			rest.WriteErrorResponse(w, http.StatusNoContent, errors.New("no content found for requested key").Error())
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
			rest.WriteErrorResponse(w, http.StatusNoContent, errors.New("no content found for requested key").Error())
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

// get validatoraccount by valID
func validatorAccountByIDHandlerFn(
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

		res, err := cliCtx.QueryStore(staking.GetValidatorAccountMapKey(hmTypes.NewValidatorID(id).Bytes()), "staking")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNoContent, errors.New("no content found for requested key").Error())
		}

		var _valAccount hmTypes.ValidatorAccount

		cdc.UnmarshalBinaryBare(res, &_valAccount)

		result, err := json.Marshal(&_valAccount)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
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
			rest.WriteErrorResponse(w, http.StatusNoContent, errors.New("no content found for requested key").Error())
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
			_validatorSet.IncrementProposerPriority(1)
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

// currentProposerHandlerFn get proposer for current validator set
func currentProposerHandlerFn(
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
			rest.WriteErrorResponse(w, http.StatusNoContent, errors.New("no content found for requested key").Error())
			return
		}

		var _validatorSet hmTypes.ValidatorSet
		err = cdc.UnmarshalBinaryBare(res, &_validatorSet)
		if err != nil {
			RestLogger.Error("Error while marshalling validator set", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := json.Marshal(&_validatorSet.Proposer)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

// initialAccountRootHandlerFn returns genesis accountroothash
func initialAccountRootHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", checkpointTypes.QuerierRoute, checkpoint.QueryInitialAccountRoot), nil)
		RestLogger.Debug("initial accountRootHash querier response", "res", res)

		if err != nil {
			RestLogger.Error("Error while calculating Initial AccountRoot  ", "Error", err.Error())
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(res) == 0 {
			RestLogger.Error("AccountRootHash not found ", "Error", err.Error())
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("AccountRootHash not found").Error())
			return
		}

		var accountRootHash = hmTypes.BytesToHeimdallHash(res)
		RestLogger.Debug("Fetched initial accountRootHash ", "AccountRootHash", accountRootHash)

		result, err := json.Marshal(&accountRootHash)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, result)
	}
}

func slashValidatorHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := r.URL.Query()
		valID, ok := rest.ParseUint64OrReturnBadRequest(w, params.Get("val_id"))
		if !ok {
			return
		}

		slashAmount, ok := big.NewInt(0).SetString(params.Get("slash_amount"), 10)
		if !ok {
			return
		}

		RestLogger.Info("Slashing validator - ", "valID", valID, "slashAmount", slashAmount)

		slashingParams := stakingTypes.NewValidatorSlashParams(types.ValidatorID(valID), slashAmount)

		bz, err := cliCtx.Codec.MarshalJSON(slashingParams)
		if err != nil {
			RestLogger.Info("Error marshallingSlashing Params - ", "err", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", stakingTypes.QuerierRoute, staking.QuerySlashValidator), bz)
		if err != nil {
			RestLogger.Info("Error query data - ", "err", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		RestLogger.Info("slashedd validator successfully ", "res", res)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
