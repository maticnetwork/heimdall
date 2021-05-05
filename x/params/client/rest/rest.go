package rest

// import (
// 	"fmt"
// 	"net/http"

// 	"github.com/cosmos/cosmos-sdk/client"
// 	"github.com/cosmos/cosmos-sdk/types/rest"
// 	"github.com/maticnetwork/heimdall/helper"
// 	govrest "github.com/maticnetwork/heimdall/x/gov/client/rest"
// 	"github.com/maticnetwork/heimdall/x/params/types/proposal"
// )

// // ProposalRESTHandler returns a ProposalRESTHandler that exposes the param
// // change REST handler with a given sub-route.
// func ProposalRESTHandler(clientCtx client.Context) govrest.ProposalRESTHandler {
// 	return govrest.ProposalRESTHandler{
// 		SubRoute: "param_change",
// 		Handler:  postProposalHandlerFn(clientCtx),
// 	}
// }

// func postProposalHandlerFn(clientCtx client.Context) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var req paramscutils.ParamChangeProposalReq
// 		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
// 			return
// 		}

// 		req.BaseReq = req.BaseReq.Sanitize()
// 		if !req.BaseReq.ValidateBasic(w) {
// 			return
// 		}

// 		content := proposal.NewParameterChangeProposal(req.Title, req.Description, req.Changes.ToParamChanges())

// 		validatorID, err := cmd.Flags().GetInt(req.FlagValidatorID)
// 		if err != nil {
// 			return err
// 		}
// 		if validatorID == 0 {
// 			return fmt.Errorf("Valid validator ID required")
// 		}

// 		msg, err := govtypes.NewMsgSubmitProposal(content, req.Deposit, req.Proposer)
// 		if rest.CheckBadRequestError(w, err) {
// 			return
// 		}
// 		if rest.CheckBadRequestError(w, msg.ValidateBasic()) {
// 			return
// 		}

// 		return helper.GenerateOrBroadcastTxCli(clientCtx, cmd.Flags(), msg)
// 	}
// }
