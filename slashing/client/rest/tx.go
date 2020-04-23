package rest

// import (
// 	"bytes"
// 	"net/http"

// 	"github.com/cosmos/cosmos-sdk/client/context"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/gorilla/mux"

// 	"github.com/maticnetwork/heimdall/staking/types"
// 	"github.com/maticnetwork/heimdall/types/rest"
// )

// func registerTxHandlers(ctx context.CLIContext, m codec.Marshaler, txg tx.Generator, r *mux.Router) {
// 	r.HandleFunc("/slashing/validators/{validatorAddr}/unjail", NewUnjailRequestHandlerFn(ctx, m, txg)).Methods("POST")
// }

// // Unjail TX body
// type UnjailReq struct {
// 	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
// }

// // NewUnjailRequestHandlerFn returns an HTTP REST handler for creating a MsgUnjail
// // transaction.
// func NewUnjailRequestHandlerFn(ctx context.CLIContext, m codec.Marshaler, txg tx.Generator) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx = ctx.WithMarshaler(m)
// 		vars := mux.Vars(r)
// 		bech32Validator := vars["validatorAddr"]

// 		var req UnjailReq
// 		if !rest.ReadRESTReq(w, r, ctx.Marshaler, &req) {
// 			return
// 		}

// 		req.BaseReq = req.BaseReq.Sanitize()
// 		if !req.BaseReq.ValidateBasic(w) {
// 			return
// 		}

// 		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		valAddr, err := sdk.ValAddressFromBech32(bech32Validator)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		if !bytes.Equal(fromAddr, valAddr) {
// 			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own validator address")
// 			return
// 		}

// 		msg := types.NewMsgUnjail(valAddr)
// 		err = msg.ValidateBasic()
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}
// 		tx.WriteGeneratedTxResponse(ctx, w, txg, req.BaseReq, msg)
// 	}
// }

// func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
// 	r.HandleFunc(
// 		"/slashing/validators/{validatorAddr}/unjail",
// 		unjailRequestHandlerFn(cliCtx),
// 	).Methods("POST")
// }

// func unjailRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		vars := mux.Vars(r)

// 		bech32validator := vars["validatorAddr"]

// 		var req UnjailReq
// 		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
// 			return
// 		}

// 		req.BaseReq = req.BaseReq.Sanitize()
// 		if !req.BaseReq.ValidateBasic(w) {
// 			return
// 		}

// 		valAddr, err := sdk.ValAddressFromBech32(bech32validator)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
// 			return
// 		}

// 		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		if !bytes.Equal(fromAddr, valAddr) {
// 			rest.WriteErrorResponse(w, http.StatusUnauthorized, "must use own validator address")
// 			return
// 		}

// 		msg := types.NewMsgUnjail(valAddr)
// 		err = msg.ValidateBasic()
// 		if err != nil {
// 			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
// 			return
// 		}

// 		authclient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
// 	}
// }
