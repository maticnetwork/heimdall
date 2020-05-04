package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/heimdall/bor/types"
	restClient "github.com/maticnetwork/heimdall/client/rest"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/bor/propose-span",
		postProposeSpanHandlerFn(cliCtx),
	).Methods("POST")
}

// ProposeSpanReq struct for proposing new span
type ProposeSpanReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	ID         uint64 `json:"span_id"`
	StartBlock uint64 `json:"start_block"`
	BorChainID string `json:"bor_chain_id"`
}

func postProposeSpanHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// read req from request
		var req ProposeSpanReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		//
		// Get span duration
		//

		// fetch duration
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryParams, types.ParamSpan), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("Span duration not found ").Error())
			return
		}

		var spanDuration uint64
		if err := json.Unmarshal(res, &spanDuration); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// fetch seed
		res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryNextSpanSeed), nil)
		if err != nil {
			RestLogger.Error("Error while fetching next span seed  ", "Error", err.Error())
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var seed common.Hash
		if err := json.Unmarshal(res, &seed); err != nil {
			return
		}

		// draft a propose span message
		msg := types.NewMsgProposeSpan(
			req.ID,
			hmTypes.HexToHeimdallAddress(req.BaseReq.From),
			req.StartBlock,
			req.StartBlock+spanDuration-1,
			req.BorChainID,
			seed,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
