package rest

import (
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/checkpoint"
	restClient "github.com/maticnetwork/heimdall/client/rest"
	"github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/checkpoint/new",
		newCheckpointHandler(cdc, cliCtx),
	).Methods("POST")
	r.HandleFunc("/checkpoint/ack", newCheckpointACKHandler(cdc, cliCtx)).Methods("POST")
	r.HandleFunc("/checkpoint/no-ack", newCheckpointNoACKHandler(cdc, cliCtx)).Methods("POST")
}

type (
	// HeaderBlockReq struct for incoming checkpoint
	HeaderBlockReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		Proposer   types.HeimdallAddress `json:"proposer"`
		RootHash   common.Hash           `json:"rootHash"`
		StartBlock uint64                `json:"startBlock"`
		EndBlock   uint64                `json:"endBlock"`
	}

	// HeaderACKReq struct for sending ACK for a new headers
	// by providing the header index assigned my mainchain contract
	HeaderACKReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		Proposer    types.HeimdallAddress `json:"proposer"`
		HeaderBlock uint64                `json:"headerBlock"`
	}

	// HeaderNoACKReq struct for sending no-ack for a new headers
	HeaderNoACKReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		Proposer types.HeimdallAddress `json:"proposer"`
	}
)

func newCheckpointHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req HeaderBlockReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// draft a message and send response
		msg := checkpoint.NewMsgCheckpointBlock(
			req.Proposer,
			req.StartBlock,
			req.EndBlock,
			req.RootHash,
			uint64(time.Now().Unix()),
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func newCheckpointACKHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req HeaderACKReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// draft a message and send response
		msg := checkpoint.NewMsgCheckpointAck(req.Proposer, req.HeaderBlock)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func newCheckpointNoACKHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req HeaderNoACKReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// draft a message and send response
		msg := checkpoint.NewMsgCheckpointNoAck(
			req.Proposer,
			uint64(time.Now().Unix()),
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
