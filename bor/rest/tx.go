package rest

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/bor"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/bor/proposeSpan",
		postProposeSpanHandlerFn(cdc, cliCtx),
	).Methods("POST")
}

type (
	// ProposeSpan struct for proposing new span
	ProposeSpan struct {
		StartBlock uint64 `json:"startBlock"`
		ChainID    string `json:"chainID"`
	}
)

func postProposeSpanHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m ProposeSpan
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(body, &m)
		if err != nil {
			RestLogger.Error("Error unmarshalling propose span ", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// fetch duration
		res, err := cliCtx.QueryStore(bor.SpanDurationKey, "bor")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNoContent, errors.New("span duration not found ").Error())
			return
		}

		duration, err := strconv.ParseInt(string(res), 10, 64)
		if err != nil {
			RestLogger.Error("Unable to parse int", "response", res, "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err = cliCtx.QueryStore(staking.ACKCountKey, "staking")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		// The query will return empty if there is no data
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNoContent, err.Error())
			return
		}
		ackCount, err := strconv.ParseInt(string(res), 10, 64)
		if err != nil {
			RestLogger.Error("Unable to parse int", "Response", res, "Error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err = cliCtx.QueryStore(staking.CurrentValidatorSetKey, "staking")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		// the query will return empty if there is no data
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNoContent, err.Error())
			return
		}
		var _validatorSet types.ValidatorSet
		cdc.UnmarshalBinaryBare(res, &_validatorSet)
		var validators []types.MinimalVal

		for _, val := range _validatorSet.Validators {
			if val.IsCurrentValidator(uint64(ackCount)) {
				// append if validator is current valdiator
				validators = append(validators, (*val).MinimalVal())
			}
		}

		msg := bor.NewMsgProposeSpan(
			m.StartBlock,
			m.StartBlock+uint64(duration),
			validators,
			validators,
			m.ChainID,
			uint64(time.Now().Unix()),
		)

		resp, err := helper.BroadcastMsgs(cliCtx, []sdk.Msg{msg})
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
