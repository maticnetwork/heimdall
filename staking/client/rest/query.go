// nolint
package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"

	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

// It represents the staking total power
//
//swagger:response stakingTotalPowerResponse
type stakingTotalPowerResponse struct {
	//in:body
	Output stakingTotalPowerStructure `json:"output"`
}

type stakingTotalPowerStructure struct {
	Height string            `json:"height"`
	Result stakingTotalPower `json:"result"`
}

type stakingTotalPower struct {
	Result int64 `json:"result"`
}

// It represents the signer by address or id
//
//swagger:response stakingValidatorResponse
type stakingValidatorResponse struct {
	//in:body
	Output stakingValidatorStructure `json:"output"`
}

type stakingValidatorStructure struct {
	Height string    `json:"height"`
	Result validator `json:"result"`
}

type validator struct {
	ID           int    `json:"ID"`
	StartEpoch   int    `json:"startEpoch"`
	EndEpoch     int    `json:"endEpoch"`
	Nonce        int    `json:"nonce"`
	Power        int    `json:"power"`
	PubKey       string `json:"pubKey"`
	Signer       string `json:"signer"`
	Last_Updated string `json:"last_updated"`
	Jailed       bool   `json:"jailed"`
	Accum        int    `json:"accum"`
}

// It represents the validor status
//
//swagger:response stakingValidatorStatusResponse
type stakingValidatorStatusResponse struct {
	//in:body
	Output stakingValidatorStatusStructure `json:"output"`
}

type stakingValidatorStatusStructure struct {
	Height string                 `json:"height"`
	Result stakingValidatorStatus `json:"result"`
}
type stakingValidatorStatus struct {
	Result bool `json:"result"`
}

// It represents the validator set
//
//swagger:response stakingValidatorSetResponse
type stakingValidatorSetResponse struct {
	//in:body
	Output stakingValidatorSetStructure `json:"output"`
}

type stakingValidatorSetStructure struct {
	Height string     `json:"height"`
	Result validators `json:"result"`
}

type validators struct {
	Validators []validator `json:"validators"`
}

//swagger:response stakingIsOldTxResponse
type stakingIsOldTxResponse struct {
	//in:body
	Output isOldTx `json:"output"`
}

type isOldTx struct {
	Height string `json:"height"`
	Result bool   `json:"result"`
}

// It represents the proposer based on time
//
//swagger:response stakingProposerByTimeResponse
type stakingProposerByTimeResponse struct {
	//in:body
	Output stakingProposerByTimeStructure `json:"output"`
}

type stakingProposerByTimeStructure struct {
	Height string      `json:"height"`
	Result []validator `json:"result"`
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/staking/totalpower",
		getTotalValidatorPower(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/staking/signer/{address}",
		validatorByAddressHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/staking/validator-status/{address}",
		validatorStatusByAddressHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/staking/validator/{id}",
		validatorByIDHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/staking/validator-set",
		validatorSetHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/staking/proposer/{times}",
		proposerHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/staking/milestoneProposer/{times}",
		milestoneProposerHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/staking/current-proposer",
		currentProposerHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/staking/proposer-bonus-percent",
		proposerBonusPercentHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/staking/isoldtx",
		StakingTxStatusHandlerFn(cliCtx),
	).Methods("GET")
}

// swagger:route GET /staking/totalpower staking stakingTotalPower
// It returns the total power of all the validators
// responses:
//
//	200: stakingTotalPowerResponse
//
// Returns total power of current validator set
func getTotalValidatorPower(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		RestLogger.Debug("Fetching total validator power")

		totalPowerBytes, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTotalValidatorPower), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if !hmRest.ReturnNotFoundIfNoContent(w, totalPowerBytes, "total power not found") {
			return
		}

		var totalPower uint64
		if err := jsoniter.ConfigFastest.Unmarshal(totalPowerBytes, &totalPower); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := jsoniter.ConfigFastest.Marshal(map[string]interface{}{"result": totalPower})
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

//swagger:parameters stakingSignerByAddress
type signerAddress struct {

	//Address of the signer
	//required:true
	//in:path
	Address string `json:"address"`
}

// swagger:route GET /staking/signer/{address} staking stakingSignerByAddress
// It returns the signer by address
// responses:
//
//	200: stakingValidatorResponse
//
// Returns validator information by signer address
func validatorByAddressHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		signerAddress := common.HexToAddress(vars["address"])

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQuerySignerParams(signerAddress.Bytes()))
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySigner), queryParams)
		if err != nil {
			RestLogger.Error("Error while fetching signer", "Error", err.Error())
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		// error if no checkpoint found
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No checkpoint found"); !ok {
			return
		}

		// return result
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters stakingValidatorStatus
type validatorAddress struct {

	//Address of the validator
	//required:true
	//in:path
	Address string `json:"address"`
}

// swagger:route GET /staking/validator-status/{address} staking stakingValidatorStatus
// It returns the status of the validator
// responses:
//
//	200: stakingValidatorStatusResponse
func validatorStatusByAddressHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		signerAddress := common.HexToAddress(vars["address"])

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQuerySignerParams(signerAddress.Bytes()))
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		statusBytes, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryValidatorStatus), queryParams)
		if err != nil {
			RestLogger.Error("Error while fetching validator status", "Error", err.Error())
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		// error if no checkpoint found
		if !hmRest.ReturnNotFoundIfNoContent(w, statusBytes, "No validator found") {
			return
		}

		var status bool
		if err = jsoniter.ConfigFastest.Unmarshal(statusBytes, &status); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := jsoniter.ConfigFastest.Marshal(map[string]interface{}{"result": status})
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// return result
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters stakingValidatorById
type validatorID struct {

	//ID of the validator
	//required:true
	//in:path
	Id int64 `json:"id"`
}

// swagger:route GET /staking/validator/{id} staking stakingValidatorById
// It returns the staking validator information by id
// responses:
//
//	200: stakingValidatorResponse
//
// Returns validator information by val ID
func validatorByIDHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get id
		id, ok := rest.ParseUint64OrReturnBadRequest(w, vars["id"])
		if !ok {
			return
		}

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryValidatorParams(hmTypes.ValidatorID(id)))
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryValidator), queryParams)
		if err != nil {
			RestLogger.Error("Error while fetching validator", "Error", err.Error())
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		// error if no checkpoint found
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No checkpoint found"); !ok {
			return
		}

		// return result
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// swagger:route GET /staking/validator-set staking stakingValidatorSet
// It returns the current validator set
// responses:
//
//	200: stakingValidatorSetResponse
//
// get current validator set
func validatorSetHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentValidatorSet), nil)
		if err != nil {
			RestLogger.Error("Error while fetching current validator set ", "Error", err.Error())
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		// error if no checkpoint found
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No checkpoint found"); !ok {
			return
		}

		// return result
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters stakingProposerByTime
type Times struct {

	//time
	//required:true
	//in:path
	Times int64 `json:"times"`
}

// swagger:route GET /staking/proposer/{times} staking stakingProposerByTime
// It returns proposer for current validator set by time
// responses:
//
//	200: stakingProposerByTimeResponse
//
// get proposer for current validator set
func proposerHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get proposer times
		times, ok := rest.ParseUint64OrReturnBadRequest(w, vars["times"])
		if !ok {
			return
		}

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryProposerParams(times))
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryProposer), queryParams)
		if err != nil {
			RestLogger.Error("Error while fetching proposers ", "Error", err.Error())
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		// error if no checkpoint found
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No proposer found"); !ok {
			return
		}

		// return result
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func milestoneProposerHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get proposer times
		times, ok := rest.ParseUint64OrReturnBadRequest(w, vars["times"])
		if !ok {
			return
		}

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryProposerParams(times))
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMilestoneProposer), queryParams)
		if err != nil {
			RestLogger.Error("Error while fetching milestoneproposers ", "Error", err.Error())
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		// error if no checkpoint found
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No milestone proposer found"); !ok {
			return
		}

		// return result
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// swagger:route GET /staking/current-proposer staking stakingCurrentProposer
// It returns proposer for current validator set
// responses:
//
//	200: stakingValidatorResponse
//
// currentProposerHandlerFn get proposer for current validator set
func currentProposerHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentProposer), nil)
		if err != nil {
			RestLogger.Error("Error while fetching current proposer ", "Error", err.Error())
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		// error if no checkpoint found
		if !hmRest.ReturnNotFoundIfNoContent(w, res, "No proposer found") {
			return
		}

		// return result
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// Returns proposer Bonus Percent information
func proposerBonusPercentHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// fetch state reocrd
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryProposerBonusPercent), nil)
		if err != nil {
			RestLogger.Error("Error while fetching proposer bonus percentage", "Error", err.Error())
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())

			return
		}

		// error if no checkpoint found
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "Proposer bonus percentage not found"); !ok {
			RestLogger.Error("Proposer bonus percentage not found ", "Error", err.Error())
			return
		}

		var _proposerBonusPercent int64
		if err := jsoniter.ConfigFastest.Unmarshal(res, &_proposerBonusPercent); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := jsoniter.ConfigFastest.Marshal(_proposerBonusPercent)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		rest.PostProcessResponse(w, cliCtx, result)
	}
}

//swagger:parameters stakingIsOldTx
type stalkingTxParams struct {

	//Log Index of the transaction
	//required:true
	//in:query
	LogIndex int64 `json:"logindex"`

	//Hash of the transaction
	//required:true
	//in:query
	Txhash string `json:"txhash"`
}

// swagger:route GET /staking/isoldtx staking stakingIsOldTx
// It returns status of the transaction
// responses:
//
//	200: stakingIsOldTxResponse
//
// Returns staking tx status information
func StakingTxStatusHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := r.URL.Query()

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get logIndex
		logindex, ok := rest.ParseUint64OrReturnBadRequest(w, vars.Get("logindex"))
		if !ok {
			return
		}

		txHash := vars.Get("txhash")

		if txHash == "" {
			return
		}

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryStakingSequenceParams(txHash, logindex))
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		seqNo, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryStakingSequence), queryParams)
		if err != nil {
			RestLogger.Error("Error while fetching staking sequence", "Error", err.Error())
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())

			return
		}

		// error if no tx status found
		if ok := hmRest.ReturnNotFoundIfNoContent(w, seqNo, "No sequence found"); !ok {
			return
		}

		res := true

		// return result
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters stakingIsOldTx stakingProposerByTime stakingCurrentProposer stakingValidatorSet stakingValidatorById stakingValidatorStatus stakingSignerByAddress stakingTotalPower
type Height struct {

	//Block Height
	//in:query
	Height string `json:"height"`
}
