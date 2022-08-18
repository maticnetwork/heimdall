//nolint
package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"

	"github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

//swagger:response topupIsOldTxResponse
type topupIsOldTxResponse struct {
	//in:body
	Output isOldTx `json:"output"`
}

type isOldTx struct {
	Height string `json:"height"`
	Result bool   `json:"result"`
}

//It represents the dividend account information
//swagger:response topupDividendAccountResponse
type topupDividendAccountResponse struct {
	//in:body
	Output topupDividendAccountStructure `json:"output"`
}

type topupDividendAccountStructure struct {
	Height string          `json:"height"`
	Result DividendAccount `json:"result"`
}

type DividendAccount struct {
	User      string `json:"user"`
	FeeAmount string `json:"feeAmount"`
}

//swagger:response topupDividendAccountRootResponse
type topupDividendAccountRootResponse struct {
	//in:body
	Output topupDividendAccountRootStructure `json:"output"`
}

type topupDividendAccountRootStructure struct {
	Height string `json:"height"`
	Result []byte `json:"result"`
}

//swagger:response topupDividendAccountProofVerifyResponse
type topupDividendAccountProofVerifyResponse struct {
	//in:body
	Output topupDividendAccountProofVerifyStructure `json:"output"`
}

type topupDividendAccountProofVerifyStructure struct {
	Height string `json:"height"`
	Result bool   `json:"result"`
}

//swagger:response topupDividendAccountProofResponse
type topupDividendAccountProofResponse struct {
	//in:body
	Output topupDividendAccountProofStructure `json:"output"`
}

type topupDividendAccountProofStructure struct {
	Height string               `json:"height"`
	Result DividendAccountProof `json:"result"`
}

type DividendAccountProof struct {
	User  string `json:"user"`
	Proof string `json:"accountProof"`
	Index uint64 `json:"index"`
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/topup/isoldtx",
		TopupTxStatusHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/topup/dividend-account/{address}",
		dividendAccountByAddressHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/topup/dividend-account-root",
		dividendAccountRootHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/topup/account-proof/{address}/verify",
		VerifyAccountProofHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(
		"/topup/account-proof/{address}",
		dividendAccountProofHandlerFn(cliCtx),
	).Methods("GET")
}

//swagger:parameters topupIsOldTx
type topupTxParams struct {
	//Log Index of the transaction
	//required:true
	//in:query
	LogIndex int64 `json:"logindex"`

	//Hash of the transaction
	//required:true
	//in:query
	Txhash string `json:"txhash"`
}

// swagger:route GET /topup/isoldtx topup topupIsOldTx
// It returns whether the transaction is old
// responses:
//   200: topupIsOldTxResponse
// Returns topup tx status information
func TopupTxStatusHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := r.URL.Query()

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get logIndex
		logindex, ok := rest.ParseUint64OrReturnBadRequest(w, vars.Get("logindex"))
		if !ok {
			return
		}

		txHash := vars.Get("txhash")
		if txHash == "" {
			return
		}

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQuerySequenceParams(txHash, logindex))
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		seqNo, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySequence), queryParams)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// error if no tx status found
		if ok := hmRest.ReturnNotFoundIfNoContent(w, seqNo, "No sequence found"); !ok {
			return
		}

		res := true

		// return result
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters topupDividendAccountProofByAddress topupDividendAccountProofVerify topupDividendAccountByAddress
type topupAddress struct {
	//Address
	//required:true
	//in:path
	Address string `json:"address"`
}

// swagger:route GET /topup/dividend-account/{address} topup topupDividendAccountByAddress
// It returns the Dividend Account information by User Address
// responses:
//   200: topupDividendAccountResponse
// Returns Dividend Account information by User Address
func dividendAccountByAddressHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get address
		userAddress := hmTypes.HexToHeimdallAddress(vars["address"])

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryDividendAccountParams(userAddress))
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDividendAccount), queryParams)
		if err != nil {
			RestLogger.Error("Error while fetching Dividend account", "Error", err.Error())
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		// error if no dividend account found
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No Dividend Account found"); !ok {
			return
		}

		// return result
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// swagger:route GET /topup/dividend-account-root topup topupDividendAccountRoot
// It returns the genesis account roothash
// responses:
//   200: topupDividendAccountRootResponse
// dividendAccountRootHandlerFn returns genesis accountroothash
func dividendAccountRootHandlerFn(
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDividendAccountRoot), nil)
		if err != nil {
			RestLogger.Error("Error while calculating dividend AccountRoot  ", "Error", err.Error())
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		// error if AccountRoot not found
		if !hmRest.ReturnNotFoundIfNoContent(w, res, "Dividend AccountRoot not found") {
			RestLogger.Error("AccountRoot not found ", "Error", err.Error())

			return
		}

		var accountRootHash = hmTypes.BytesToHeimdallHash(res)

		RestLogger.Debug("Fetched Dividend accountRootHash ", "AccountRootHash", accountRootHash)

		result, err := jsoniter.ConfigFastest.Marshal(&accountRootHash)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		// return result

		rest.PostProcessResponse(w, cliCtx, result)
	}
}

// swagger:route GET /topup/account-proof/{address} topup topupDividendAccountProofByAddress
// It returns the account proof by User Address
// responses:
//   200: topupDividendAccountProofResponse
// Returns Merkle path for dividendAccountID using dividend Account Tree
func dividendAccountProofHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get id
		userAddress := hmTypes.HexToHeimdallAddress(vars["address"])

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryAccountProofParams(userAddress))
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAccountProof), queryParams)
		if err != nil {
			RestLogger.Error("Error while fetching merkle proof", "Error", err.Error())
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		// error if account proof  not found
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No proof found"); !ok {
			return
		}

		// return result
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// swagger:route GET /topup/account-proof/{address}/verify topup topupDividendAccountProofVerify
// It returns true if given Merkle path for dividendAccountID is valid
// responses:
//   200: topupDividendAccountProofVerifyResponse
// VerifyAccountProofHandlerFn - Returns true if given Merkle path for dividendAccountID is valid
func VerifyAccountProofHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		params := r.URL.Query()
		userAddress := hmTypes.HexToHeimdallAddress(vars["address"])
		accountProof := params.Get("proof")

		RestLogger.Info("Verify Account Proof", "userAddress", userAddress, "accountProof", accountProof)

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryVerifyAccountProofParams(userAddress, accountProof))
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryVerifyAccountProof), queryParams)
		if err != nil {
			RestLogger.Error("Error while verifying merkle proof", "Error", err.Error())
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		var accountProofStatus bool
		if err = jsoniter.ConfigFastest.Unmarshal(res, &accountProofStatus); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err = jsoniter.ConfigFastest.Marshal(map[string]interface{}{"result": accountProofStatus})
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// return result
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters topupDividendAccountProofVerify topupDividendAccountProofByAddress topupDividendAccountRoot topupDividendAccountByAddress topupIsOldTx
type Height struct {

	//Block Height
	//in:query
	Height string `json:"height"`
}
