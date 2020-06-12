package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	restClient "github.com/maticnetwork/heimdall/client/rest"
	"github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/staking/validators",
		newValidatorJoinHandler(cliCtx),
	).Methods("POST")
	r.HandleFunc("/staking/validators/stake", newValidatorStakeUpdateHandler(cliCtx)).Methods("PUT")
	r.HandleFunc("/staking/validators", newValidatorUpdateHandler(cliCtx)).Methods("PUT")
	r.HandleFunc("/staking/validators", newValidatorExitHandler(cliCtx)).Methods("DELETE")
}

type (
	// AddValidatorReq add validator request object
	AddValidatorReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		ID              uint64         `json:"ID"`
		ActivationEpoch uint64         `json:"activationEpoch"`
		Amount          string         `json:"amount"`
		SignerPubKey    hmTypes.PubKey `json:"pubKey"`
		TxHash          string         `json:"tx_hash"`
		LogIndex        uint64         `json:"log_index"`
		BlockNumber     uint64         `json:"block_number" yaml:"block_number"`
		Nonce           uint64         `json:"nonce"`
	}

	// UpdateSignerReq update validator signer request object
	UpdateSignerReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		ID              uint64         `json:"ID"`
		NewSignerPubKey hmTypes.PubKey `json:"pubKey"`
		TxHash          string         `json:"tx_hash"`
		LogIndex        uint64         `json:"log_index"`
		BlockNumber     uint64         `json:"block_number" yaml:"block_number"`
		Nonce           uint64         `json:"nonce"`
	}

	// UpdateValidatorStakeReq update validator stake request object
	UpdateValidatorStakeReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		ID          uint64 `json:"ID"`
		Amount      string `json:"amount"`
		TxHash      string `json:"tx_hash"`
		LogIndex    uint64 `json:"log_index"`
		BlockNumber uint64 `json:"block_number" yaml:"block_number"`
		Nonce       uint64 `json:"nonce"`
	}

	// RemoveValidatorReq remove validator request object
	RemoveValidatorReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		ID                uint64 `json:"ID"`
		DeactivationEpoch uint64 `json:"deactivationEpoch"`
		TxHash            string `json:"tx_hash"`
		LogIndex          uint64 `json:"log_index"`
		BlockNumber       uint64 `json:"block_number" yaml:"block_number"`
		Nonce             uint64 `json:"nonce"`
	}
)

func newValidatorJoinHandler(cliCtx context.CLIContext) http.HandlerFunc {
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

		amount, ok := sdk.NewIntFromString(req.Amount)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid amount")
		}

		// create new msg
		msg := types.NewMsgValidatorJoin(
			hmTypes.HexToHeimdallAddress(req.BaseReq.From),
			req.ID,
			req.ActivationEpoch,
			amount,
			req.SignerPubKey,
			hmTypes.HexToHeimdallHash(req.TxHash),
			req.LogIndex,
			req.BlockNumber,
			req.Nonce,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func newValidatorExitHandler(cliCtx context.CLIContext) http.HandlerFunc {
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
		msg := types.NewMsgValidatorExit(
			hmTypes.HexToHeimdallAddress(req.BaseReq.From),
			req.ID,
			req.DeactivationEpoch,
			hmTypes.HexToHeimdallHash(req.TxHash),
			req.LogIndex,
			req.BlockNumber,
			req.Nonce,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func newValidatorUpdateHandler(cliCtx context.CLIContext) http.HandlerFunc {
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
		msg := types.NewMsgSignerUpdate(
			hmTypes.HexToHeimdallAddress(req.BaseReq.From),
			req.ID,
			req.NewSignerPubKey,
			hmTypes.HexToHeimdallHash(req.TxHash),
			req.LogIndex,
			req.BlockNumber,
			req.Nonce,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func newValidatorStakeUpdateHandler(cliCtx context.CLIContext) http.HandlerFunc {
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

		amount, ok := sdk.NewIntFromString(req.Amount)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid amount")
		}

		// create msg validator update
		msg := types.NewMsgStakeUpdate(
			hmTypes.HexToHeimdallAddress(req.BaseReq.From),
			req.ID,
			amount,
			hmTypes.HexToHeimdallHash(req.TxHash),
			req.LogIndex,
			req.BlockNumber,
			req.Nonce,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
