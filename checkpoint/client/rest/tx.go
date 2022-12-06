//nolint
package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	restClient "github.com/maticnetwork/heimdall/client/rest"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

//It represents New checkpoint msg.
//swagger:response checkpointNewResponse
type checkpointNewResponse struct {
	//in:body
	Output checkpointNew `json:"output"`
}

type checkpointNew struct {
	Type  string             `json:"type"`
	Value checkpointNewValue `json:"value"`
}

type checkpointNewValue struct {
	Msg       checkpointNewMsg `json:"msg"`
	Signature string           `json:"signature"`
	Memo      string           `json:"memo"`
}

type checkpointNewMsg struct {
	Type  string           `json:"type"`
	Value checkpointNewVal `json:"value"`
}

type checkpointNewVal struct {
	Proposer        string `json:"proposer"`
	StartBlock      string `json:"start_block"`
	EndBlock        string `json:"end_block"`
	RootHash        string `json:"root_hash"`
	AccountRootHash string `json:"account_root_hash"`
	BorChainId      string `json:"bor_chain_id"`
}

//It represents Propose Span msg.
//swagger:response checkpointAckResponse
type checkpointAckResponse struct {
	//in:body
	Output checkpointAck `json:"output"`
}

type checkpointAck struct {
	Type  string             `json:"type"`
	Value checkpointAckValue `json:"value"`
}

type checkpointAckValue struct {
	Msg       checkpointAckMsg `json:"msg"`
	Signature string           `json:"signature"`
	Memo      string           `json:"memo"`
}

type checkpointAckMsg struct {
	Type  string           `json:"type"`
	Value checkpointAckVal `json:"value"`
}

type checkpointAckVal struct {
	From       string `json:"from"`
	Number     string `json:"number"`
	StartBlock string `json:"start_block"`
	EndBlock   string `json:"end_block"`
	Proposer   string `json:"proposer"`
	RootHash   string `json:"root_Hash"`
	TxHash     string `json:"tx_hash"`
	LogIndex   string `json:"log_index"`
}

//It represents Propose Span msg.
//swagger:response checkpointNoAckResponse
type checkpointNoAckResponse struct {
	//in:body
	Output checkpointNoAck `json:"output"`
}

type checkpointNoAck struct {
	Type  string               `json:"type"`
	Value checkpointNoAckValue `json:"value"`
}

type checkpointNoAckValue struct {
	Msg       checkpointNoAckMsg `json:"msg"`
	Signature string             `json:"signature"`
	Memo      string             `json:"memo"`
}

type checkpointNoAckMsg struct {
	Type  string             `json:"type"`
	Value checkpointNoAckVal `json:"value"`
}

type checkpointNoAckVal struct {
	From string `json:"from"`
}

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/checkpoint/new",
		newCheckpointHandler(cliCtx),
	).Methods("POST")
	r.HandleFunc("/checkpoint/ack", newCheckpointACKHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/checkpoint/no-ack", newCheckpointNoACKHandler(cliCtx)).Methods("POST")
}

type (
	// HeaderBlockReq struct for incoming checkpoint
	HeaderBlockReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		Proposer        hmTypes.HeimdallAddress `json:"proposer"`
		RootHash        hmTypes.HeimdallHash    `json:"root_Hash"`
		AccountRootHash hmTypes.HeimdallHash    `json:"account_root_hash"`
		StartBlock      uint64                  `json:"start_block"`
		EndBlock        uint64                  `json:"end_block"`
		BorChainID      string                  `json:"bor_chain_id"`
	}

	// HeaderACKReq struct for sending ACK for a new headers
	// by providing the header index assigned my mainchain contract
	HeaderACKReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		From        hmTypes.HeimdallAddress `json:"from"`
		HeaderBlock uint64                  `json:"header_block"`
		StartBlock  uint64                  `json:"start_block"`
		EndBlock    uint64                  `json:"end_block"`
		Proposer    hmTypes.HeimdallAddress `json:"proposer"`
		RootHash    hmTypes.HeimdallHash    `json:"root_Hash"`
		TxHash      hmTypes.HeimdallHash    `json:"tx_hash"`
		LogIndex    uint64                  `json:"log_index"`
	}

	// HeaderNoACKReq struct for sending no-ack for a new headers
	HeaderNoACKReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		Proposer hmTypes.HeimdallAddress `json:"proposer"`
	}
)

//swagger:parameters checkpointNew
type checkpointNewParam struct {

	//Body
	//required:true
	//in:body
	Input checkpointNewInput `json:"input"`
}

type checkpointNewInput struct {
	BaseReq         BaseReq `json:"base_req"`
	Proposer        string  `json:"proposer"`
	RootHash        string  `json:"root_Hash"`
	AccountRootHash string  `json:"account_root_hash"`
	StartBlock      string  `json:"start_block"`
	EndBlock        string  `json:"end_block"`
	BorChainID      string  `json:"bor_chain_id"`
}

type BaseReq struct {

	//Address of the sender
	From string `json:"address"`

	//Chain ID of Heimdall
	ChainID string `json:"chain_id"`
}

// swagger:route POST /checkpoint/new checkpoint checkpointNew
// It returns the prepared msg for new checkpoint
// responses:
//   200: checkpointNewResponse

func newCheckpointHandler(cliCtx context.CLIContext) http.HandlerFunc {
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
		msg := types.NewMsgCheckpointBlock(
			req.Proposer,
			req.StartBlock,
			req.EndBlock,
			req.RootHash,
			req.AccountRootHash,
			req.BorChainID,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

//swagger:parameters checkpointAck
type checkpointAckParam struct {

	//Body
	//required:true
	//in:body
	Input checkpointAckInput `json:"input"`
}

type checkpointAckInput struct {

	//required:true
	BaseReq BaseReq `json:"base_req"`

	From        string `json:"from"`
	HeaderBlock string `json:"header_block"`
	StartBlock  string `json:"start_block"`
	EndBlock    string `json:"end_block"`
	Proposer    string `json:"proposer"`
	RootHash    string `json:"root_Hash"`
	TxHash      string `json:"tx_hash"`
	LogIndex    string `json:"log_index"`
}

// swagger:route POST /checkpoint/ack checkpoint checkpointAck
// It returns the prepared msg for ack checkpoint
// responses:
//   200: checkpointAckResponse

func newCheckpointACKHandler(cliCtx context.CLIContext) http.HandlerFunc {
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
		msg := types.NewMsgCheckpointAck(
			req.From,
			req.HeaderBlock,
			req.Proposer,
			req.StartBlock,
			req.EndBlock,
			req.RootHash,
			req.TxHash,
			req.LogIndex,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

//swagger:parameters checkpointNoAck
type checkpointNoAckParam struct {

	//Body
	//required:true
	//in:body
	Input checkpointAckInput `json:"input"`
}

type checkpointNoAckInput struct {

	//required:true
	BaseReq BaseReq `json:"base_req"`

	Proposer string `json:"proposer"`
}

// swagger:route POST /checkpoint/no-ack checkpoint checkpointNoAck
// It returns the prepared msg for no-ack checkpoint
// responses:
//   200: checkpointNoAckResponse

func newCheckpointNoACKHandler(cliCtx context.CLIContext) http.HandlerFunc {
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
		msg := types.NewMsgCheckpointNoAck(
			req.Proposer,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
