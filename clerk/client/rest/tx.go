package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	restClient "github.com/maticnetwork/heimdall/client/rest"
	"github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/clerk/records",
		newEventRecordHandler(cdc, cliCtx),
	).Methods("POST")
}

// AddRecordReq add validator request object
type AddRecordReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	TxHash   types.HeimdallHash `json:"tx_hash"`
	LogIndex uint64             `json:"log_index"`
	ID       uint64             `json:"id"`
}

func newEventRecordHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read req from request
		var req AddRecordReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// create new msg
		msg := clerkTypes.NewMsgEventRecord(
			types.HexToHeimdallAddress(req.BaseReq.From),
			req.TxHash,
			req.LogIndex,
			req.ID,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
