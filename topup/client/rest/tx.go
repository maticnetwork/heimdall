package rest

import (
	"math/big"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	restClient "github.com/maticnetwork/heimdall/client/rest"
	topupTypes "github.com/maticnetwork/heimdall/topup/types"
	"github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/topup/fee", TopupHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc("/topup/withdraw", WithdrawFeeHandlerFn(cliCtx)).Methods("POST")
}

//
// Topup req
//

// TopupReq defines the properties of a topup request's body.
type TopupReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	TxHash      string `json:"tx_hash" yaml:"tx_hash"`
	LogIndex    uint64 `json:"log_index" yaml:"log_index"`
	User        string `json:"user" yaml:"user"`
	Fee         string `json:"fee" yaml:"fee"`
	BlockNumber uint64 `json:"block_number" yaml:"block_number"`
}

// TopupHandlerFn - http request handler to topup coins to a address.
func TopupHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TopupReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// get from address
		fromAddr := types.HexToHeimdallAddress(req.BaseReq.From)

		// get signer
		user := types.HexToHeimdallAddress(req.User)

		// fee amount
		fee, ok := sdk.NewIntFromString(req.Fee)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid amount")
		}

		msg := topupTypes.NewMsgTopup(
			fromAddr,
			user,
			fee,
			types.HexToHeimdallHash(req.TxHash),
			req.LogIndex,
			req.BlockNumber,
		)

		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

//
// Withdraw Fee req
//

// WithdrawFeeReq defines the properties of a withdraw fee request's body.
type WithdrawFeeReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  string       `json:"amount" yaml:"amount"`
}

// WithdrawFeeHandlerFn - http request handler to withdraw fee coins from a address.
func WithdrawFeeHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req WithdrawFeeReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Bad request")
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Request validation failed")
			return
		}

		// get from address
		fromAddr := types.HexToHeimdallAddress(req.BaseReq.From)
		amountStr := "0"
		if req.Amount != "" {
			amountStr = req.Amount
		}
		amount, ok := big.NewInt(0).SetString(amountStr, 10)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Bad amount")
			return
		}

		// get msg
		msg := topupTypes.NewMsgWithdrawFee(
			fromAddr,
			sdk.NewIntFromBigInt(amount),
		)
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
