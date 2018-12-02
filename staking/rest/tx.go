package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	hmType "github.com/maticnetwork/heimdall/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/staking/validators",
		newValidatorJoinHandler(cliCtx),
	).Methods("POST")
	r.HandleFunc("/staking/validators", newValidatorUpdateHandler(cliCtx)).Methods("PUT")
	r.HandleFunc("/staking/validators", newValidatorExitHandler(cliCtx)).Methods("DELETE")
}

type addValidator struct {
	ValidatorAddress string `json:"address"`
	SignerPubKey     string `json:"pubKey"`
}

func newValidatorJoinHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m addValidator

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = json.Unmarshal(body, &m)
		if err != nil {
			RestLogger.Error("Error unmarshalling json while adding validator", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		validatorAddr := common.HexToAddress(m.ValidatorAddress)
		crytoPubKey := helper.BytesToPubkey(common.Hex2Bytes(m.SignerPubKey))

		msg := staking.NewMsgValidatorJoin(validatorAddr, common.Hex2Bytes(m.SignerPubKey))

		txBytes, err := helper.CreateTxBytes(msg)
		if err != nil {
			RestLogger.Error("Unable to create txBytes", "ValidatorAddress", validatorAddr.String(), "ValidatorPubKey", helper.BytesToPubkey(crytoPubKey.Bytes()))
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		resp, err := helper.SendTendermintRequest(cliCtx, txBytes)
		if err != nil {
			RestLogger.Error("Error while sending request to Tendermint", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		result, err := json.Marshal(&resp)
		if err != nil {
			RestLogger.Error("Error while marshalling tendermint response", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.Write(result)
	}
}

type remoteValidator struct {
	SignerAddress common.Address `json:"address"`
}

func newValidatorExitHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m remoteValidator

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = json.Unmarshal(body, &m)
		if err != nil {
			RestLogger.Error("Error unmarshalling json epoch checkpoint", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		msg := staking.NewMsgValidatorExit(m.SignerAddress)

		txBytes, err := helper.CreateTxBytes(msg)
		if err != nil {
			RestLogger.Error("Unable to create txBytes", "validatorAddress", m.SignerAddress)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		resp, err := helper.SendTendermintRequest(cliCtx, txBytes)
		if err != nil {
			RestLogger.Error("Error while sending request to Tendermint", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		result, err := json.Marshal(&resp)
		if err != nil {
			RestLogger.Error("Error while marshalling tendermint response", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(result)
	}
}

type updateValidator struct {
	ValidatorAddress   common.Address `json:"address"`
	NewValidatorPubKey hmType.PubKey  `json:"newPubKey"`
}

func newValidatorUpdateHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m updateValidator

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = json.Unmarshal(body, &m)
		if err != nil {
			RestLogger.Error("Error unmarshalling json epoch checkpoint", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		msg := staking.NewMsgValidatorUpdate(m.ValidatorAddress, m.NewValidatorPubKey[:])

		txBytes, err := helper.CreateTxBytes(msg)
		if err != nil {
			RestLogger.Error("Unable to create txBytes", "currentValidatorAddress", m.ValidatorAddress.Hex(), "newValidatorPubKey", m.NewValidatorPubKey)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		resp, err := helper.SendTendermintRequest(cliCtx, txBytes)
		if err != nil {
			RestLogger.Error("Error while sending request to Tendermint", "error", err, "ValidatorAddress", m.ValidatorAddress.Hex(), "newValidatorPubKey", m.NewValidatorPubKey, "txBytes", txBytes)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		result, err := json.Marshal(&resp)
		if err != nil {
			RestLogger.Error("Error while marshalling tendermint response", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(result)
	}
}
