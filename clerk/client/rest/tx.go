//nolint
package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/bridge/setu/util"
	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	restClient "github.com/maticnetwork/heimdall/client/rest"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

//It represents New checkpoint msg.
//swagger:response clerkNewEventResponse
type clerkNewEventResponse struct {
	//in:body
	Output clerkNewEvent `json:"output"`
}

type clerkNewEvent struct {
	Type  string             `json:"type"`
	Value clerkNewEventValue `json:"value"`
}

type clerkNewEventValue struct {
	Msg       clerkNewEventMsg `json:"msg"`
	Signature string           `json:"signature"`
	Memo      string           `json:"memo"`
}

type clerkNewEventMsg struct {
	Type  string           `json:"type"`
	Value clerkNewEventVal `json:"value"`
}

type clerkNewEventVal struct {
	From            string `json:"from"`
	TxHash          string `json:"tx_hash"`
	LogIndex        string `json:"log_index"`
	BlockNumber     string `json:"block_number" yaml:"block_number"`
	ID              string `json:"id"`
	ContractAddress string `json:"contract_address" yaml:"contract_address"`
	BorChainID      string `json:"bor_chain_id"`
	Data            string `json:"data"`
}

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/clerk/records",
		newEventRecordHandler(cliCtx),
	).Methods("POST")
}

// AddRecordReq add validator request object
type AddRecordReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	TxHash          types.HeimdallHash `json:"tx_hash"`
	LogIndex        uint64             `json:"log_index"`
	BlockNumber     uint64             `json:"block_number" yaml:"block_number"`
	ID              uint64             `json:"id"`
	ContractAddress string             `json:"contract_address" yaml:"contract_address"`
	BorChainID      string             `json:"bor_chain_id"`
	Data            string             `json:"data"`
}

//swagger:parameters clerkNewEvent
type clerkNewEventParam struct {

	//Body
	//required:true
	//in:body
	Input clerkNewEventInput `json:"input"`
}

type clerkNewEventInput struct {
	BaseReq BaseReq `json:"base_req"`

	TxHash          string `json:"tx_hash"`
	LogIndex        string `json:"log_index"`
	BlockNumber     string `json:"block_number"`
	ID              string `json:"id"`
	ContractAddress string `json:"contract_address"`
	BorChainID      string `json:"bor_chain_id"`
	Data            string `json:"data"`
}

type BaseReq struct {

	//Address of the sender
	From string `json:"address"`

	//Chain ID of Heimdall
	ChainID string `json:"chain_id"`
}

// swagger:route POST /clerk/records  clerk clerkNewEvent
// It returns the prepared msg for new clerk event
// responses:
//   200: clerkNewEventResponse

func newEventRecordHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read req from request
		var req AddRecordReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// get ContractAddress
		contractAddress := types.HexToHeimdallAddress(req.ContractAddress)

		if util.GetBlockHeight(cliCtx) > helper.GetSpanOverrideHeight() && len(types.HexToHexBytes(req.Data)) > helper.MaxStateSyncSize {
			RestLogger.Info(`Data is too large to process, Resetting to ""`, "id", req.ID)
			req.Data = ""
		} else if len(types.HexToHexBytes(req.Data)) > helper.LegacyMaxStateSyncSize {
			RestLogger.Info(`Data is too large to process, Resetting to ""`, "id", req.ID)
			req.Data = ""
		}

		// create new msg
		msg := clerkTypes.NewMsgEventRecord(
			types.HexToHeimdallAddress(req.BaseReq.From),
			req.TxHash,
			req.LogIndex,
			req.BlockNumber,
			req.ID,
			contractAddress,
			types.HexToHexBytes(req.Data),
			req.BorChainID,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
