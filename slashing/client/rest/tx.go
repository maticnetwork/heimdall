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

//It represents unjail msg.
//swagger:response slashingUnjailResponse
type slashingUnjailResponse struct {
	//in:body
	Output slashingUnjailOutput `json:"output"`
}

type slashingUnjailOutput struct {
	Type  string              `json:"type"`
	Value slashingUnjailValue `json:"value"`
}

type slashingUnjailValue struct {
	Msg       slashingUnjailMsg `json:"msg"`
	Signature string            `json:"signature"`
	Memo      string            `json:"memo"`
}

type slashingUnjailMsg struct {
	Type  string            `json:"type"`
	Value slashingUnjailVal `json:"value"`
}

type slashingUnjailVal struct {
	From        string `json:"from"`
	ID          uint64 `json:"id"`
	TxHash      string `json:"tx_hash"`
	LogIndex    uint64 `json:"log_index"`
	BlockNumber uint64 `json:"block_number"`
}

//It represents Propose Span msg.
//swagger:response slashingNewTickResponse
type slashingNewTickResponse struct {
	//in:body
	Output slashingNewTickOutput `json:"output"`
}

type slashingNewTickOutput struct {
	Type  string               `json:"type"`
	Value slashingNewTickValue `json:"value"`
}

type slashingNewTickValue struct {
	Msg       slashingNewTickMsg `json:"msg"`
	Signature string             `json:"signature"`
	Memo      string             `json:"memo"`
}

type slashingNewTickMsg struct {
	Type  string             `json:"type"`
	Value slashingNewTickVal `json:"value"`
}

type slashingNewTickVal struct {
	ID                uint64 `json:"id"`
	Proposer          string `json:"proposer"`
	SlashingInfoBytes string `json:"slashinginfobytes"`
}

//It represents Propose Span msg.
//swagger:response slashingTickAckResponse
type slashingTickAckResponse struct {
	//in:body
	Output slashingTickAckOutput `json:"output"`
}

type slashingTickAckOutput struct {
	Type  string               `json:"type"`
	Value slashingTickAckValue `json:"value"`
}

type slashingTickAckValue struct {
	Msg       slashingTickAckMsg `json:"msg"`
	Signature string             `json:"signature"`
	Memo      string             `json:"memo"`
}

type slashingTickAckMsg struct {
	Type  string             `json:"type"`
	Value slashingTickAckVal `json:"value"`
}

type slashingTickAckVal struct {
	From        string `json:"from"`
	ID          uint64 `json:"tick_id"`
	Amount      uint64 `json:"slashed_amount"`
	TxHash      string `json:"tx_hash"`
	LogIndex    uint64 `json:"log_index"`
	BlockNumber uint64 `json:"block_number"`
}

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

//swagger:parameters slashingUnjail
type slashingUnjailParam struct {

	//Body
	//required:true
	//in:body
	Input slashingUnjailInput `json:"input"`

	//Validator Address
	//required:true
	//in:path
	ValidatorAddr string `json:"validatorAddr"`
}

type slashingUnjailInput struct {
	BaseReq     BaseReq `json:"base_req"`
	ID          uint64  `json:"ID"`
	TxHash      string  `json:"tx_hash"`
	LogIndex    uint64  `json:"log_index"`
	BlockNumber uint64  `json:"block_number" yaml:"block_number"`
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

// swagger:route POST /slashing/validators/{validatorAddr}/unjail slashing slashingUnjail
// It returns the prepared msg for unjail
// responses:
//   200: slashingUnjailResponse
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

//swagger:parameters slashingNewTick
type slashingNewTickParam struct {

	//Body
	//required:true
	//in:body
	Input slashingNewTickInput `json:"input"`
}

type slashingNewTickInput struct {
	BaseReq           BaseReq `json:"base_req"`
	ID                uint64  `json:"ID"`
	Proposer          string  `json:"proposer"`
	SlashingInfoBytes string  `json:"slashing_info_bytes"`
}

// swagger:route POST /slashing/tick slashing slashingNewTick
// It returns the prepared msg for new tick
// responses:
//   200: slashingNewTickResponse
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

//swagger:parameters slashingTickAck
type slashingTickAckParam struct {

	//Body
	//required:true
	//in:body
	Input slashingTickAckInput `json:"input"`
}

type slashingTickAckInput struct {
	BaseReq     BaseReq `json:"base_req"`
	ID          uint64  `json:"ID"`
	Amount      uint64  `json:"amount"`
	TxHash      string  `json:"tx_hash"`
	LogIndex    uint64  `json:"log_index"`
	BlockNumber uint64  `json:"block_number" yaml:"block_number"`
}

// swagger:route POST slashing/tick-ack slashing slashingTickAck
// It returns the prepared msg for tick-ack
// responses:
//   200: stakingTickAckResponse
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
