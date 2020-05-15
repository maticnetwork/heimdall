package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"

	restClient "github.com/maticnetwork/heimdall/client/rest"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/slashing/validators/{validatorAddr}/unjail",
		newUnjailRequestHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		"/slashing/tick",
		newTickRequestHandlerFn(cliCtx),
	).Methods("POST")

	r.HandleFunc(
		"/slashing/tick-ack",
		newTickAckHandler(cliCtx),
	).Methods("POST")
}

// Unjail TX body
type UnjailReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	ID          uint64 `json:"ID"`
	TxHash      string `json:"tx_hash"`
	LogIndex    uint64 `json:"log_index"`
	BlockNumber uint64 `json:"block_number" yaml:"block_number"`
}

type TickReq struct {
	BaseReq           rest.BaseReq `json:"base_req"`
	ID                uint64       `json:"ID"`
	Proposer          string       `json:"proposer"`
	SlashingInfoBytes string       `json:"slashing_info_bytes"`
}

type TickAckReq struct {
	BaseReq     rest.BaseReq `json:"base_req"`
	ID          uint64       `json:"ID"`
	Amount      uint64       `json:"amount"`
	TxHash      string       `json:"tx_hash"`
	LogIndex    uint64       `json:"log_index"`
	BlockNumber uint64       `json:"block_number" yaml:"block_number"`
}

func newUnjailRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read req from Request
		var req UnjailReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgUnjail(
			hmTypes.HexToHeimdallAddress(req.BaseReq.From),
			req.ID,
			hmTypes.HexToHeimdallHash(req.TxHash),
			req.LogIndex,
			req.BlockNumber,
		)
		err := msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func newTickRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// read req from Request
		var req TickReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgTick(
			req.ID,
			hmTypes.HexToHeimdallAddress(req.Proposer),
			hmTypes.HexToHexBytes(req.SlashingInfoBytes),
		)

		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})

	}
}

func newTickAckHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req TickAckReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		msg := types.NewMsgTickAck(
			hmTypes.HexToHeimdallAddress(req.BaseReq.From),
			req.ID,
			req.Amount,
			hmTypes.HexToHeimdallHash(req.TxHash),
			req.LogIndex,
			req.BlockNumber,
		)

		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
