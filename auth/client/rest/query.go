//nolint
package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/types"
	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

//It represents the auth params
//swagger:response authParamsResponse
type authParamsResponse struct {
	//in:body
	Output authParamsStructure `json:"output"`
}

type authParamsStructure struct {
	Height string     `json:"height"`
	Result authParams `json:"result"`
}

type authParams struct {
	MaxMemoCharacters      string `json:"max_memo_characters"`
	TxSigLimit             int64  `json:"tx_sig_limit"`
	TxSizeCostPerByte      int64  `json:"tx_size_cost_per_byte"`
	SigVerifyCostEd2219    int64  `json:"sig_verify_cost_ed25519"`
	SigVerifyCostSecp256k1 int64  `json:"sig_verify_cost_secp256k1"`
	MaxTxGas               int64  `json:"max_tx_gas"`
	TxFees                 int64  `json:"tx_fees"`
}

//swagger:response authAccountSequenceResponse
type authAccountSequenceResponse struct {
	//in:body
	Output authAccountSequenceStructure `json:"output"`
}

type authAccountSequenceStructure struct {
	Height string              `json:"height"`
	Result authAccountSequence `json:"result"`
}

type authAccountSequence struct {
	Address       string `json:"address"`
	AccountNumber string `json:"account_number"`
	Sequence      string `json:"sequence"`
}

//swagger:response authAccountResponse
type authAccountResponse struct {
	//in:body
	Output authAccountStructure `json:"output"`
}

type authAccountStructure struct {
	Height string      `json:"height"`
	Result authAccount `json:"result"`
}

type authAccount struct {
	Type  string `json:"type"`
	Value value  `json:"value"`
}

type value struct {
	Address       string    `json:"address"`
	Coins         []coin    `json:"coins"`
	PublicKey     publicKey `json:"public_key"`
	AccountNumber string    `json:"account_number"`
	Sequence      string    `json:"sequence"`
}

type coin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type publicKey struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

//swagger:parameters authAccount authAccountSequence
type address struct {

	//Account Address
	//in:path
	//required:true
	Address string `json:"address"`
}

// swagger:route GET /auth/accounts/{address} auth authAccount
// It returns the account details.
// responses:
//   200: authAccountResponse
// QueryAccountRequestHandlerFn query account REST Handler
func QueryAccountRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)

		// key
		key := types.HexToHeimdallAddress(vars["address"])
		if key.Empty() {
			hmRest.WriteErrorResponse(w, http.StatusNotFound, errors.New("Invalid address").Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// account getter
		accGetter := authTypes.NewAccountRetriever(cliCtx)

		account, height, err := accGetter.GetAccountWithHeight(key)
		if err != nil {
			if err := accGetter.EnsureExists(key); err != nil {
				cliCtx = cliCtx.WithHeight(height)
				hmRest.PostProcessResponse(w, cliCtx, authTypes.BaseAccount{})

				return
			}

			hmRest.WriteErrorResponse(w, http.StatusNotFound, err.Error())

			return
		}

		cliCtx = cliCtx.WithHeight(height)
		hmRest.PostProcessResponse(w, cliCtx, account)
	}
}

// swagger:route GET /auth/accounts/{address}/sequence auth authAccountSequence
// It returns the account sequence
// responses:
//   200: authAccountSequenceResponse
// QueryAccountSequenceRequestHandlerFn query account sequence REST Handler
func QueryAccountSequenceRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)

		// key
		key := types.HexToHeimdallAddress(vars["address"])
		if key.Empty() {
			hmRest.WriteErrorResponse(w, http.StatusNotFound, errors.New("Invalid address").Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// account getter
		accGetter := authTypes.NewAccountRetriever(cliCtx)

		account, height, err := accGetter.GetAccountWithHeight(key)
		if err != nil {
			if err := accGetter.EnsureExists(key); err != nil {
				cliCtx = cliCtx.WithHeight(height)
				hmRest.PostProcessResponse(w, cliCtx, authTypes.BaseAccount{})

				return
			}

			hmRest.WriteErrorResponse(w, http.StatusNotFound, err.Error())

			return
		}

		// get result
		result := authTypes.LightBaseAccount{
			Address:       account.GetAddress(),
			Sequence:      account.GetSequence(),
			AccountNumber: account.GetAccountNumber(),
		}

		cliCtx = cliCtx.WithHeight(height)
		hmRest.PostProcessResponse(w, cliCtx, result)
	}
}

// swagger:route GET /auth/params auth authParams
// It returns the auth parameters.
// responses:
//   200: authParamsResponse
// HTTP request handler to query the auth params values
func paramsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", authTypes.QuerierRoute, authTypes.QueryParams)

		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
