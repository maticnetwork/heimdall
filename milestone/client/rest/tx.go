//nolint
package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	restClient "github.com/maticnetwork/heimdall/client/rest"
	"github.com/maticnetwork/heimdall/milestone/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

//It represents New milestone msg.
//swagger:response milestoneNewResponse
type milestoneNewResponse struct {
	//in:body
	Output milestoneNew `json:"output"`
}

type milestoneNew struct {
	Type  string            `json:"type"`
	Value milestoneNewValue `json:"value"`
}

type milestoneNewValue struct {
	Msg       milestoneNewMsg `json:"msg"`
	Signature string          `json:"signature"`
	Memo      string          `json:"memo"`
}

type milestoneNewMsg struct {
	Type  string          `json:"type"`
	Value milestoneNewVal `json:"value"`
}

type milestoneNewVal struct {
	Proposer   string `json:"proposer"`
	StartBlock string `json:"start_block"`
	EndBlock   string `json:"end_block"`
	RootHash   string `json:"root_hash"`
	BorChainId string `json:"bor_chain_id"`
}

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/milestone",
		newMilestoneHandler(cliCtx),
	).Methods("POST")

}

type (
	// HeaderBlockReq struct for incoming checkpoint
	HeaderBlockReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		Proposer   hmTypes.HeimdallAddress `json:"proposer"`
		RootHash   hmTypes.HeimdallHash    `json:"root_Hash"`
		StartBlock uint64                  `json:"start_block"`
		EndBlock   uint64                  `json:"end_block"`
		BorChainID string                  `json:"bor_chain_id"`
	}
)

//swagger:parameters milestoneNew
type milestoneNewParam struct {

	//Body
	//required:true
	//in:body
	Input milestoneNewInput `json:"input"`
}

type milestoneNewInput struct {
	BaseReq    BaseReq `json:"base_req"`
	Proposer   string  `json:"proposer"`
	RootHash   string  `json:"root_Hash"`
	StartBlock string  `json:"start_block"`
	EndBlock   string  `json:"end_block"`
	BorChainID string  `json:"bor_chain_id"`
}

type BaseReq struct {

	//Address of the sender
	From string `json:"address"`

	//Chain ID of Heimdall
	ChainID string `json:"chain_id"`
}

// swagger:route POST /milestone/new milestone milestoneNew
// It returns the prepared msg for new milestone
// responses:
//   200: milestoneNewResponse

func newMilestoneHandler(cliCtx context.CLIContext) http.HandlerFunc {
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
		msg := types.NewMsgMilestoneBlock(
			req.Proposer,
			req.StartBlock,
			req.EndBlock,
			req.RootHash,
			req.BorChainID,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
