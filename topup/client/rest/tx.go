//nolint
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

//It represents topup fee msg.
//swagger:response topupFeeResponse
type topupFeeResponse struct {
	//in:body
	Output topupFeeOutput `json:"output"`
}

type topupFeeOutput struct {
	Type  string        `json:"type"`
	Value topupFeeValue `json:"value"`
}

type topupFeeValue struct {
	Msg       topupFeeMsg `json:"msg"`
	Signature string      `json:"signature"`
	Memo      string      `json:"memo"`
}

type topupFeeMsg struct {
	Type  string      `json:"type"`
	Value topupFeeVal `json:"value"`
}

type topupFeeVal struct {
	FromAddress string `json:"from_address"`
	TxHash      string `json:"tx_hash"`
	LogIndex    uint64 `json:"log_index"`
	User        string `json:"user"`
	Fee         string `json:"fee"`
	BlockNumber uint64 `json:"block_number"`
}

//It represents topup withdraw msg.
//swagger:response topupWithdrawResponse
type topupWithdrawResponse struct {
	//in:body
	Output topupWithdrawOutput `json:"output"`
}

type topupWithdrawOutput struct {
	Type  string             `json:"type"`
	Value topupWithdrawValue `json:"value"`
}

type topupWithdrawValue struct {
	Msg       topupWithdrawMsg `json:"msg"`
	Signature string           `json:"signature"`
	Memo      string           `json:"memo"`
}

type topupWithdrawMsg struct {
	Type  string           `json:"type"`
	Value topupWithdrawVal `json:"value"`
}

type topupWithdrawVal struct {
	FromAddress string `json:"from_address"`
	Amount      string `json:"amount"`
}

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

//swagger:parameters topupFee
type topupFeeParam struct {

	//Body
	//required:true
	//in:body
	Input topupFeeInput `json:"input"`
}

type topupFeeInput struct {
	BaseReq     BaseReq `json:"base_req"`
	TxHash      string  `json:"tx_hash"`
	LogIndex    uint64  `json:"log_index"`
	User        string  `json:"user"`
	Fee         string  `json:"fee"`
	BlockNumber uint64  `json:"block_number"`
}

type BaseReq struct {

	//Address of the sender
	//required:true
	//in:body
	From string `json:"address"`

	//Chain ID of Heimdall
	//required:true
	//in:body
	ChainID string `json:"chain_id"`
}

// swagger:route POST /topup/fee topup topupFee
// It returns the prepared msg for topup fee
// responses:
//   200: topupFeeResponse
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

//swagger:parameters topupWithdraw
type topupWithdrawParam struct {

	//Body
	//required:true
	//in:body
	Input topupWithdrawInput `json:"input"`
}

type topupWithdrawInput struct {
	BaseReq BaseReq `json:"base_req"`
	Amount  string  `json:"amount"`
}

// swagger:route POST /topup/withdraw topup topupWithdraw
// It returns the prepared msg for topup withdraw
// responses:
//   200: topupWithdrawResponse
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
