package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	restClient "github.com/maticnetwork/heimdall/client/rest"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types"
	hmType "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/staking/validators",
		newValidatorJoinHandler(cdc, cliCtx),
	).Methods("POST")
	r.HandleFunc("/staking/validators/stake", newValidatorStakeUpdateHandler(cdc, cliCtx)).Methods("PUT")
	r.HandleFunc("/staking/validators", newValidatorUpdateHandler(cdc, cliCtx)).Methods("PUT")
	r.HandleFunc("/staking/validators", newValidatorExitHandler(cdc, cliCtx)).Methods("DELETE")
}

type (
	// AddValidatorReq add validator request object
	AddValidatorReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		ID           uint64        `json:"ID"`
		SignerPubKey hmType.PubKey `json:"pubKey"`
		TxHash       string        `json:"tx_hash"`
		LogIndex     uint64        `json:"log_index"`
	}

	// UpdateSignerReq update validator signer request object
	UpdateSignerReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		ID              uint64        `json:"ID"`
		NewSignerPubKey hmType.PubKey `json:"pubKey"`
		TxHash          string        `json:"tx_hash"`
		LogIndex        uint64        `json:"log_index"`
	}

	// UpdateValidatorStakeReq update validator stake request object
	UpdateValidatorStakeReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		ID       uint64 `json:"ID"`
		TxHash   string `json:"tx_hash"`
		LogIndex uint64 `json:"log_index"`
	}

	// RemoveValidatorReq remove validator request object
	RemoveValidatorReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		ID       uint64 `json:"ID"`
		TxHash   string `json:"tx_hash"`
		LogIndex uint64 `json:"log_index"`
	}
)

func newValidatorJoinHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read req from request
		var req AddValidatorReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// create new msg
		msg := staking.NewMsgValidatorJoin(
			types.HexToHeimdallAddress(req.BaseReq.From),
			req.ID,
			req.SignerPubKey,
			types.HexToHeimdallHash(req.TxHash),
			req.LogIndex,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func newValidatorExitHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read req from request
		var req RemoveValidatorReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// draft new msg
		msg := staking.NewMsgValidatorExit(
			types.HexToHeimdallAddress(req.BaseReq.From),
			req.ID,
			types.HexToHeimdallHash(req.TxHash),
			req.LogIndex,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func newValidatorUpdateHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read req from request
		var req UpdateSignerReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// create msg validator update
		msg := staking.NewMsgSignerUpdate(
			types.HexToHeimdallAddress(req.BaseReq.From),
			req.ID,
			req.NewSignerPubKey,
			types.HexToHeimdallHash(req.TxHash),
			req.LogIndex,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func newValidatorStakeUpdateHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read req from request
		var req UpdateValidatorStakeReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// create msg validator update
		msg := staking.NewMsgStakeUpdate(
			types.HexToHeimdallAddress(req.BaseReq.From),
			req.ID,
			types.HexToHeimdallHash(req.TxHash),
			req.LogIndex,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
