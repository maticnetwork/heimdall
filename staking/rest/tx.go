package rest

import (
	"encoding/json"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
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
	}

	// UpdateValidatorReq update validator request object
	UpdateValidatorReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		ID              uint64        `json:"ID"`
		NewSignerPubKey hmType.PubKey `json:"pubKey"`
		NewAmount       json.Number   `json:"amount"`
		TxHash          string        `json:"tx_hash"`
	}

	// RemoveValidatorReq remove validator request object
	RemoveValidatorReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		ID     uint64 `json:"ID"`
		TxHash string `json:"tx_hash"`
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
			common.HexToHash(req.TxHash),
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
			common.HexToHash(req.TxHash),
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func newValidatorUpdateHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read req from request
		var req UpdateValidatorReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// create msg validator update
		msg := staking.NewMsgValidatorUpdate(
			types.HexToHeimdallAddress(req.BaseReq.From),
			req.ID,
			req.NewSignerPubKey,
			common.HexToHash(req.TxHash),
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
