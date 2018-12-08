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
	"math/rand"
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
	ValidatorAddress common.Address `json:"address"`
	SignerPubKey     hmType.PubKey  `json:"pubKey"`
	StartEpoch       uint64         `json:"startEpoch"`
	EndEpoch         uint64         `json:"endEpoch"`
	Amount           uint64         `json:"Amount"`
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

		// create new msg
		msg := staking.NewMsgValidatorJoin(m.ValidatorAddress, m.SignerPubKey, m.StartEpoch, m.EndEpoch, m.Amount, uint64(rand.Int()))

		txBytes, err := helper.CreateTxBytes(msg)
		if err != nil {
			RestLogger.Error("Unable to create txBytes", "ValidatorAddress", m.ValidatorAddress.String(), "ValidatorPubKey", m.SignerPubKey)
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

type removeValidator struct {
	ValidatorAddress common.Address `json:"address"`
}

func newValidatorExitHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m removeValidator

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

		msg := staking.NewMsgValidatorExit(m.ValidatorAddress, uint64(rand.Int()))

		txBytes, err := helper.CreateTxBytes(msg)
		if err != nil {
			RestLogger.Error("Unable to create txBytes", "validatorAddress", m.ValidatorAddress)
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
	ValidatorAddress common.Address `json:"address"`
	NewSignerPubKey  hmType.PubKey  `json:"pubKey"`
	NewAmount        uint64         `json:"amount"`
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

		// create msg validator update
		msg := staking.NewMsgValidatorUpdate(m.ValidatorAddress, m.NewSignerPubKey, m.NewAmount, uint64(rand.Int()))

		txBytes, err := helper.CreateTxBytes(msg)
		if err != nil {
			RestLogger.Error("Unable to create txBytes", "currentValidatorAddress", m.ValidatorAddress.Hex(), "newSignerPubKey", m.NewSignerPubKey)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		resp, err := helper.SendTendermintRequest(cliCtx, txBytes)
		if err != nil {
			RestLogger.Error("Error while sending request to Tendermint", "error", err, "ValidatorAddress", m.ValidatorAddress.Hex(), "newSignerPubKey", m.NewSignerPubKey, "txBytes", txBytes)
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
