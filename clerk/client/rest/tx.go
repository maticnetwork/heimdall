package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	restClient "github.com/maticnetwork/heimdall/client/rest"
	"github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

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

		if len(types.HexToHexBytes(req.Data)) > 10 {
			fmt.Println("Data is too large", req.ID)
			req.Data = "0x0"
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
