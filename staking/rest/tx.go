package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	hmType "github.com/maticnetwork/heimdall/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/staking/validators",
		newValidatorJoinHandler(cdc, cliCtx),
	).Methods("POST")
	r.HandleFunc("/staking/validators", newValidatorUpdateHandler(cdc, cliCtx)).Methods("PUT")
	r.HandleFunc("/staking/validators", newValidatorExitHandler(cdc, cliCtx)).Methods("DELETE")
}

type (
	removeValidator struct {
		ID     uint64 `json:"ID"`
		TxHash string `json:"tx_hash"`
	}
	addValidator struct {
		ID           uint64        `json:"ID"`
		SignerPubKey hmType.PubKey `json:"pubKey"`
		TxHash       string        `json:"tx_hash"`
	}
	updateValidator struct {
		ID              uint64        `json:"ID"`
		NewSignerPubKey hmType.PubKey `json:"pubKey"`
		NewAmount       json.Number   `json:"amount"`
		TxHash          string        `json:"tx_hash"`
	}
)

func newValidatorJoinHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m addValidator

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(body, &m)
		if err != nil {
			RestLogger.Error("Error unmarshalling json while adding validator", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create new msg
		msg := staking.NewMsgValidatorJoin(m.ID, m.SignerPubKey, common.HexToHash(m.TxHash))

		txBytes, err := helper.CreateTxBytes(msg)
		if err != nil {
			RestLogger.Error("Unable to create txBytes", "ValidatorID", m.ID, "ValidatorPubKey", m.SignerPubKey)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := helper.SendTendermintRequest(cliCtx, txBytes, helper.BroadcastAsync)
		if err != nil {
			RestLogger.Error("Error while sending request to Tendermint", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		result, err := json.Marshal(&resp)
		if err != nil {
			RestLogger.Error("Error while marshalling tendermint response", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}
		rest.PostProcessResponse(w, cdc, result, cliCtx.Indent)
	}
}

func newValidatorExitHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m removeValidator

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		err = json.Unmarshal(body, &m)
		if err != nil {
			RestLogger.Error("Error unmarshalling json epoch checkpoint", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		msg := staking.NewMsgValidatorExit(m.ID, common.HexToHash(m.TxHash))

		txBytes, err := helper.CreateTxBytes(msg)
		if err != nil {
			RestLogger.Error("Unable to create txBytes", "validatorID", m.ID)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		resp, err := helper.SendTendermintRequest(cliCtx, txBytes, helper.BroadcastAsync)
		if err != nil {
			RestLogger.Error("Error while sending request to Tendermint", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		result, err := json.Marshal(&resp)
		if err != nil {
			RestLogger.Error("Error while marshalling tendermint response", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}
		rest.PostProcessResponse(w, cdc, result, cliCtx.Indent)
	}
}

func newValidatorUpdateHandler(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m updateValidator

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		err = json.Unmarshal(body, &m)
		if err != nil {
			RestLogger.Error("Error unmarshalling json epoch checkpoint", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		// create msg validator update
		msg := staking.NewMsgValidatorUpdate(m.ID, m.NewSignerPubKey, m.NewAmount, common.HexToHash(m.TxHash))

		txBytes, err := helper.CreateTxBytes(msg)
		if err != nil {
			RestLogger.Error("Unable to create txBytes", "validatorID", m.ID, "newSignerPubKey", m.NewSignerPubKey)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		resp, err := helper.SendTendermintRequest(cliCtx, txBytes, helper.BroadcastAsync)
		if err != nil {
			RestLogger.Error("Error while sending request to Tendermint", "error", err, "validatorID", m.ID, "newSignerPubKey", m.NewSignerPubKey, "txBytes", txBytes)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		result, err := json.Marshal(&resp)
		if err != nil {
			RestLogger.Error("Error while marshalling tendermint response", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}
		rest.PostProcessResponse(w, cdc, result, cliCtx.Indent)
	}
}
