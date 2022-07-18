//nolint
package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	restClient "github.com/maticnetwork/heimdall/client/rest"
	"github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

//It represents Propose Span msg.
//swagger:response stakingNewValidatorResponse
type stakingNewValidatorResponse struct {
	//in:body
	Output stakingNewValidatorOutput `json:"output"`
}

type stakingNewValidatorOutput struct {
	Type  string                   `json:"type"`
	Value stakingNewValidatorValue `json:"value"`
}

type stakingNewValidatorValue struct {
	Msg       stakingNewValidatorMsg `json:"msg"`
	Signature string                 `json:"signature"`
	Memo      string                 `json:"memo"`
}

type stakingNewValidatorMsg struct {
	Type  string                 `json:"type"`
	Value stakingNewValidatorVal `json:"value"`
}

type stakingNewValidatorVal struct {
	From            string `json:"from"`
	ID              uint64 `json:"ID"`
	ActivationEpoch uint64 `json:"activationEpoch"`
	Amount          string `json:"amount"`
	SignerPubKey    string `json:"pubKey"`
	TxHash          string `json:"tx_hash"`
	LogIndex        uint64 `json:"log_index"`
	BlockNumber     uint64 `json:"block_number" yaml:"block_number"`
	Nonce           uint64 `json:"nonce"`
}

//It represents validator stake update msg
//swagger:response stakingValidatorStakeUpdateResponse
type stakingValidatorStakeUpdateResponse struct {
	//in:body
	Output stakingValidatorStakeUpdateOutput `json:"output"`
}

type stakingValidatorStakeUpdateOutput struct {
	Type  string                           `json:"type"`
	Value stakingValidatorStakeUpdateValue `json:"value"`
}

type stakingValidatorStakeUpdateValue struct {
	Msg       stakingValidatorStakeUpdateMsg `json:"msg"`
	Signature string                         `json:"signature"`
	Memo      string                         `json:"memo"`
}

type stakingValidatorStakeUpdateMsg struct {
	Type  string                         `json:"type"`
	Value stakingValidatorStakeUpdateVal `json:"value"`
}

type stakingValidatorStakeUpdateVal struct {
	From        string `json:"from"`
	ID          uint64 `json:"ID"`
	Amount      string `json:"amount"`
	TxHash      string `json:"tx_hash"`
	LogIndex    uint64 `json:"log_index"`
	BlockNumber uint64 `json:"block_number"`
	Nonce       uint64 `json:"nonce"`
}

//It represents validator update msg.
//swagger:response stakingValidatorUpdateResponse
type stakingValidatorUpdateResponse struct {
	//in:body
	Output stakingValidatorUpdateOutput `json:"output"`
}

type stakingValidatorUpdateOutput struct {
	Type  string                      `json:"type"`
	Value stakingValidatorUpdateValue `json:"value"`
}

type stakingValidatorUpdateValue struct {
	Msg       stakingValidatorUpdateMsg `json:"msg"`
	Signature string                    `json:"signature"`
	Memo      string                    `json:"memo"`
}

type stakingValidatorUpdateMsg struct {
	Type  string                    `json:"type"`
	Value stakingValidatorUpdateVal `json:"value"`
}

type stakingValidatorUpdateVal struct {
	From            string `json:"from"`
	ID              uint64 `json:"id"`
	NewSignerPubKey string `json:"pub_key"`
	TxHash          string `json:"tx_hash"`
	LogIndex        uint64 `json:"log_index"`
	BlockNumber     uint64 `json:"block_number"`
	Nonce           uint64 `json:"nonce"`
}

//It represents validator deleting msg
//swagger:response stakingValidatorDeleteResponse
type stakingValidatorDeleteResponse struct {
	//in:body
	Output stakingValidatorDeleteOutput `json:"output"`
}

type stakingValidatorDeleteOutput struct {
	Type  string                      `json:"type"`
	Value stakingValidatorDeleteValue `json:"value"`
}

type stakingValidatorDeleteValue struct {
	Msg       stakingValidatorDeleteMsg `json:"msg"`
	Signature string                    `json:"signature"`
	Memo      string                    `json:"memo"`
}

type stakingValidatorDeleteMsg struct {
	Type  string                    `json:"type"`
	Value stakingValidatorDeleteVal `json:"value"`
}

type stakingValidatorDeleteVal struct {
	From              string `json:"from"`
	ID                uint64 `json:"id"`
	DeactivationEpoch uint64 `json:"deactivationEpoch"`
	TxHash            string `json:"tx_hash"`
	LogIndex          uint64 `json:"log_index"`
	BlockNumber       uint64 `json:"block_number"`
	Nonce             uint64 `json:"nonce"`
}

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/staking/validators",
		newValidatorJoinHandler(cliCtx),
	).Methods("POST")
	r.HandleFunc("/staking/validators/stake", newValidatorStakeUpdateHandler(cliCtx)).Methods("PUT")
	r.HandleFunc("/staking/validators", newValidatorUpdateHandler(cliCtx)).Methods("PUT")
	r.HandleFunc("/staking/validators", newValidatorExitHandler(cliCtx)).Methods("DELETE")
}

type (
	// AddValidatorReq add validator request object
	AddValidatorReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		ID              uint64         `json:"ID"`
		ActivationEpoch uint64         `json:"activationEpoch"`
		Amount          string         `json:"amount"`
		SignerPubKey    hmTypes.PubKey `json:"pubKey"`
		TxHash          string         `json:"tx_hash"`
		LogIndex        uint64         `json:"log_index"`
		BlockNumber     uint64         `json:"block_number" yaml:"block_number"`
		Nonce           uint64         `json:"nonce"`
	}

	// UpdateSignerReq update validator signer request object
	UpdateSignerReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		ID              uint64         `json:"ID"`
		NewSignerPubKey hmTypes.PubKey `json:"pubKey"`
		TxHash          string         `json:"tx_hash"`
		LogIndex        uint64         `json:"log_index"`
		BlockNumber     uint64         `json:"block_number" yaml:"block_number"`
		Nonce           uint64         `json:"nonce"`
	}

	// UpdateValidatorStakeReq update validator stake request object
	UpdateValidatorStakeReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		ID          uint64 `json:"ID"`
		Amount      string `json:"amount"`
		TxHash      string `json:"tx_hash"`
		LogIndex    uint64 `json:"log_index"`
		BlockNumber uint64 `json:"block_number" yaml:"block_number"`
		Nonce       uint64 `json:"nonce"`
	}

	// RemoveValidatorReq remove validator request object
	RemoveValidatorReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		ID                uint64 `json:"ID"`
		DeactivationEpoch uint64 `json:"deactivationEpoch"`
		TxHash            string `json:"tx_hash"`
		LogIndex          uint64 `json:"log_index"`
		BlockNumber       uint64 `json:"block_number" yaml:"block_number"`
		Nonce             uint64 `json:"nonce"`
	}
)

//swagger:parameters stakingNewValidator
type stakingNewValidatorParam struct {

	//Body
	//required:true
	//in:body
	Input stakingNewValidatorInput `json:"input"`
}

type stakingNewValidatorInput struct {
	BaseReq         BaseReq `json:"base_req"`
	ID              uint64  `json:"ID"`
	ActivationEpoch uint64  `json:"activationEpoch"`
	Amount          string  `json:"amount"`
	SignerPubKey    string  `json:"pubKey"`
	TxHash          string  `json:"tx_hash"`
	LogIndex        uint64  `json:"log_index"`
	BlockNumber     uint64  `json:"block_number"`
	Nonce           uint64  `json:"nonce"`
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

// swagger:route POST /staking/validators staking stakingNewValidator
// It returns the prepared msg for new validator joining
// responses:
//   200: stakingNewValidatorResponse
func newValidatorJoinHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read req from request
		var req AddValidatorReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		amount, ok := sdk.NewIntFromString(req.Amount)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid amount")
		}

		// create new msg
		msg := types.NewMsgValidatorJoin(
			hmTypes.HexToHeimdallAddress(req.BaseReq.From),
			req.ID,
			req.ActivationEpoch,
			amount,
			req.SignerPubKey,
			hmTypes.HexToHeimdallHash(req.TxHash),
			req.LogIndex,
			req.BlockNumber,
			req.Nonce,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

//swagger:parameters stakingValidatorDelete
type stakingValidatorDelete struct {

	//Body
	//required:true
	//in:body
	Input stakingValidatorDeleteInput `json:"input"`
}

type stakingValidatorDeleteInput struct {
	BaseReq           BaseReq `json:"base_req"`
	ID                uint64  `json:"ID"`
	DeactivationEpoch uint64  `json:"deactivationEpoch"`
	TxHash            string  `json:"tx_hash"`
	LogIndex          uint64  `json:"log_index"`
	BlockNumber       uint64  `json:"block_number" yaml:"block_number"`
	Nonce             uint64  `json:"nonce"`
}

// swagger:route DELETE /staking/validators staking stakingValidatorDelete
// It returns the prepared msg for deleting the Validator
// responses:
//   200: stakingValidatorDeleteResponse
func newValidatorExitHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read req from request
		var req RemoveValidatorReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// draft new msg
		msg := types.NewMsgValidatorExit(
			hmTypes.HexToHeimdallAddress(req.BaseReq.From),
			req.ID,
			req.DeactivationEpoch,
			hmTypes.HexToHeimdallHash(req.TxHash),
			req.LogIndex,
			req.BlockNumber,
			req.Nonce,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

//swagger:parameters stakingValidatorUpdate
type stakingValidatorUpdate struct {

	//Body
	//required:true
	//in:body
	Input stakingValidatorUpdateInput `json:"input"`
}

type stakingValidatorUpdateInput struct {
	BaseReq         BaseReq `json:"base_req"`
	ID              uint64  `json:"ID"`
	NewSignerPubKey string  `json:"pubKey"`
	TxHash          string  `json:"tx_hash"`
	LogIndex        uint64  `json:"log_index"`
	BlockNumber     uint64  `json:"block_number" yaml:"block_number"`
	Nonce           uint64  `json:"nonce"`
}

// swagger:route PUT /staking/validators staking stakingValidatorUpdate
// It returns the prepared msg for updating the validator
// responses:
//   200: stakingValidatorUpdateResponse
func newValidatorUpdateHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read req from request
		var req UpdateSignerReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// create msg validator update
		msg := types.NewMsgSignerUpdate(
			hmTypes.HexToHeimdallAddress(req.BaseReq.From),
			req.ID,
			req.NewSignerPubKey,
			hmTypes.HexToHeimdallHash(req.TxHash),
			req.LogIndex,
			req.BlockNumber,
			req.Nonce,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

//swagger:parameters stakingValidatorStakeUpdate
type stakingValidatorStakeUpdate struct {

	//Body
	//required:true
	//in:body
	Input stakingValidatorStakeUpdateInput `json:"input"`
}

type stakingValidatorStakeUpdateInput struct {
	BaseReq     BaseReq `json:"base_req"`
	ID          uint64  `json:"ID"`
	Amount      string  `json:"amount"`
	TxHash      string  `json:"tx_hash"`
	LogIndex    uint64  `json:"log_index"`
	BlockNumber uint64  `json:"block_number" yaml:"block_number"`
	Nonce       uint64  `json:"nonce"`
}

// swagger:route PUT /staking/validators/stake  staking stakingValidatorStakeUpdate
// It returns the prepared msg for updating the validator's stake
// responses:
//   200: stakingValidatorStakeUpdateResponse

func newValidatorStakeUpdateHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// read req from request
		var req UpdateValidatorStakeReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		amount, ok := sdk.NewIntFromString(req.Amount)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "invalid amount")
		}

		// create msg validator update
		msg := types.NewMsgStakeUpdate(
			hmTypes.HexToHeimdallAddress(req.BaseReq.From),
			req.ID,
			amount,
			hmTypes.HexToHeimdallHash(req.TxHash),
			req.LogIndex,
			req.BlockNumber,
			req.Nonce,
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
