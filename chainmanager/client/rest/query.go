//nolint
package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"

	chainTypes "github.com/maticnetwork/heimdall/chainmanager/types"
)

//It represents the bank balance of particluar account
//swagger:response chainManagerParamsResponse
type chainManagerParamsResponse struct {
	//in:body
	Output chainManagerParams `json:"output"`
}

type chainManagerParams struct {
	Height string       `json:"height"`
	Result chainManager `json:"result"`
}

type chainManager struct {
	MainChainConfirmation int `json:"mainchain_tx_confirmations"`

	MaticChainConfirmation int `json:"maticchain_tx_confirmations"`

	ChainManager ContractAddresses `json:"chain_params"`
}

type ContractAddresses struct {
	BorChainId             string `json:"bor_chain_id"`
	MaticChainAddress      string `json:"matic_token_address"`
	StalkingManagerAddress string `json:"staking_manager_address"`
	SlashManagerAddress    string `json:"slash_manager_address"`
	RootChainAddress       string `json:"root_chain_address"`
	StalkignInfoAddress    string `json:"staking_info_address"`
	StateSenderAddress     string `json:"state_sender_address"`
	StateReceiverAddress   string `json:"state_receiver_address"`
	ValidatorSetAddress    string `json:"validator_set_address"`
}

// swagger:route GET /chainmanager/params chain-manager chainManagerParams
// It returns the chain-manager parameters
// responses:
//   200: chainManagerParamsResponse
// HTTP request handler to query the auth params values
func paramsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", chainTypes.QuerierRoute, chainTypes.QueryParams)

		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters chainManagerParams
type Height struct {

	//Block Height
	//in:query
	Height string `json:"height"`
}
