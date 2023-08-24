//nolint
package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"

	"github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

//swagger:response slashingSigningInfoByIdResponse
type slashingSigningInfoByIdResponse struct {
	//in:body
	Output slashingSigningInfoByIdStructure `json:"output"`
}

type slashingSigningInfoByIdStructure struct {
	Height string              `json:"height"`
	Result slashingSigningInfo `json:"result"`
}

//swagger:response slashingInfosResponse
type slashingInfosResponse struct {
	//in:body
	Output slashingInfosStructure `json:"output"`
}

type slashingInfosStructure struct {
	Height string                `json:"height"`
	Result []slashingSigningInfo `json:"result"`
}

type slashingSigningInfo struct {
	ValID       int64 `json:"valID"`
	StartHeight int64 `json:"startHeight"`
	IndexOffset int64 `json:"indexOffset"`
}

//swagger:response slashingLatestInfoByIdResponse
type slashingLatestInfoByIdResponse struct {
	//in:body
	Output slashingLatestInfoByIdStructure `json:"output"`
}

type slashingLatestInfoByIdStructure struct {
	Height string                `json:"height"`
	Result ValidatorSlashingInfo `json:"result"`
}

//swagger:response slashingLatestInfosResponse
type slashingLatestInfosResponse struct {
	//in:body
	Output slashingLatestInfosStructure `json:"output"`
}

type slashingLatestInfosStructure struct {
	Height string                  `json:"height"`
	Result []ValidatorSlashingInfo `json:"result"`
}

type ValidatorSlashingInfo struct {
	ID            int64 `json:"ID"`
	SlashedAmount int64 `json:"SlashedAmount"`
	IsJailed      bool  `json:"IsJailed"`
}

//It represents the slashing parameters
//swagger:response slashingParametersResponse
type slashingParametersResponse struct {
	//in:body
	Output slashingParametersStructure `json:"output"`
}

type slashingParametersStructure struct {
	Height string `json:"height"`
	Result params `json:"result"`
}

type params struct {
	SignedBlockWindow       int64  `json:"signed_blocks_window"`
	MinSignedPerWindow      string `json:"min_signed_per_window"`
	DowntimeJailDuration    int64  `json:"downtime_jail_duration"`
	SlashFunctionDoubleSign string `json:"slash_fraction_double_sign"`
	SlashFractionDowntime   string `json:"slash_fraction_downtime"`
	SlashFractionLimit      string `json:"slash_fraction_limit"`
	JailFractionLimit       string `json:"jail_fraction_limit"`
	MaxEvidenceAge          string `json:"max_evidence_age"`
	EnableSlashing          bool   `json:"enable_slashing"`
}

//It represents the slashing count
//swagger:response slashingCountResponse
type slashingCountResponse struct {
	//in:body
	Output slashingCountStructure `json:"output"`
}

type slashingCountStructure struct {
	Height string `json:"height"`
	Result int64  `json:"result"`
}

//It represents the slashing count
//swagger:response slashingInfosBytesResponse
type slashingInfosBytesResponse struct {
	//in:body
	Output slashingInfosBytesStructure `json:"output"`
}

type slashingInfosBytesStructure struct {
	Height string `json:"height"`
	Result string `json:"result"`
}

//swagger:response slashingIsOldTxResponse
type slashingIsOldTxResponse struct {
	//in:body
	Output isOldTx `json:"output"`
}

type isOldTx struct {
	Height string `json:"height"`
	Result bool   `json:"result"`
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/slashing/validators/{id}/signing_info",
		signingInfoHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/slashing/signing_infos",
		signingInfoHandlerListFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/slashing/validators/{id}/latest_slash_info",
		latestSlashInfoHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/slashing/latest_slash_infos",
		latestSlashInfoHandlerListFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/slashing/tick_slash_infos",
		tickSlashInfoHandlerListFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/slashing/latest_slash_info_bytes",
		latestSlashInfoBytesHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/slashing/parameters",
		queryParamsHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/slashing/isoldtx",
		SlashingTxStatusHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/slashing/tick-count",
		tickCountHandlerFn(cliCtx),
	).Methods("GET")
}

//swagger:parameters slashingSigningInfoById
type validatorID struct {

	//ID of the validator
	//required:true
	//in:path
	Id int64 `json:"id"`
}

// swagger:route GET /slashing/validators/{id}/signing_info slashing slashingSigningInfoById
// It returns the signing infos of the validator based on Id
// responses:
//   200: slashingSigningInfoByIdResponse
// http request handler to query signing info
func signingInfoHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
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

		params := types.NewQuerySigningInfoParams(hmTypes.ValidatorID(id))

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySigningInfo)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// swagger:route GET /slashing/signing_infos slashing slashingInfos
// It returns the signing infos.
// responses:
//   200: slashingInfosResponse
// http request handler to query signing info
func signingInfoHandlerListFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQuerySigningInfosParams(page, limit)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySigningInfos)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters slashingLatestInfoById
type ID struct {

	//ID of the validator
	//required:true
	//in:path
	Id int64 `json:"id"`
}

// swagger:route GET /slashing/validators/{id}/latest_slash_info slashing slashingLatestInfoById
// It returns the latest signing infos of the validator based on Id
// responses:
//   200: slashingLatestInfoByIdResponse
// http request handler to query slashing info
func latestSlashInfoHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
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

		params := types.NewQuerySlashingInfoParams(hmTypes.ValidatorID(id))

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySlashingInfo)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// swagger:route GET /slashing/latest_slash_infos slashing slashingLatestInfos
// It returns the latest signing infos
// responses:
//   200: slashingLatestInfosResponse
// http request handler to query signing info
func latestSlashInfoHandlerListFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQuerySlashingInfosParams(page, limit)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySlashingInfos)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// swagger:route GET /slashing/latest_slash_info_bytes slashing slashingLatestSlashInfoBytes
// It returns the latest signing info byte
// responses:
//		200:slashingInfosBytesResponse
// http request handler to query signing info
func latestSlashInfoBytesHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySlashingInfoBytes), nil)
		RestLogger.Debug("slashInfoBytes querier response", "res", res)

		if err != nil {
			RestLogger.Error("Error while calculating slashInfoBytes  ", "Error", err.Error())
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// error if no slashInfoBytes found
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "SlashInfoBytes not found"); !ok {
			RestLogger.Error("SlashInfoBytes not found ", "Error", err.Error())
			return
		}

		var slashInfoBytes = hmTypes.BytesToHexBytes(res)
		RestLogger.Debug("Fetched slashInfoBytes ", "SlashInfoBytes", slashInfoBytes.String())

		result, err := jsoniter.ConfigFastest.Marshal(&slashInfoBytes)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// return result
		rest.PostProcessResponse(w, cliCtx, result)

	}
}

// swagger:route GET /slashing/parameters slashing slashingParameters
// It returns the slashing parameters
// responses:
//   200: slashingParametersResponse
func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/parameters", types.QuerierRoute)

		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters slashingTickInfos
type slashingTickInfosParams struct {

	//Page number
	//required:true
	//in:query
	Page int64 `json:"page"`

	//Limit per page
	//required:true
	//in:query
	Limit int64 `json:"limit"`
}

// swagger:route GET /slashing/tick_slash_infos slashing slashingTickInfos
// It returns the tick slash infos
// responses:
//   200: slashingLatestInfosResponse
// http request handler to query tick slashing info
func tickSlashInfoHandlerListFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryTickSlashingInfosParams(page, limit)
		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTickSlashingInfos)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters slashingIsOldTx
type slashingTxParams struct {
	//Log Index of the transaction
	//required:true
	//in:query
	LogIndex int64 `json:"logindex"`

	//Hash of the transaction
	//required:true
	//in:query
	Txhash string `json:"txhash"`
}

// swagger:route GET /slashing/isoldtx slashing slashingIsOldTx
// It returns whether the transaction is old
// responses:
//   200: slashingIsOldTxResponse
// Returns slashing tx status information
func SlashingTxStatusHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
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
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQuerySlashingSequenceParams(txHash, logindex))
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		seqNo, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySlashingSequence), queryParams)
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

// swagger:route GET /slashing/tick-count slashing slashingTickCount
// It returns the slashing tick count
// responses:
//   200: slashingCountResponse
func tickCountHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		RestLogger.Debug("Fetching number of ticks from state")
		tickCountBytes, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTickCount), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, tickCountBytes, "No tick count found"); !ok {
			return
		}

		var tickCount uint64
		if err := jsoniter.ConfigFastest.Unmarshal(tickCountBytes, &tickCount); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := jsoniter.ConfigFastest.Marshal(&tickCount)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

//swagger:parameters slashingTickCount slashingIsOldTx slashingTickInfos slashingParameters slashingLatestSlashInfoBytes slashingLatestInfos slashingLatestInfoById slashingSigningInfoById slashingInfos
type Height struct {

	//Block Height
	//in:query
	Height string `json:"height"`
}
