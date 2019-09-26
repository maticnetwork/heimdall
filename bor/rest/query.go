package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/bor"
	borTypes "github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/checkpoint"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Get span details from start block
	r.HandleFunc("/bor/span/{id}", getSpanHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/bor/latest-span", getLatestSpanHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/bor/span-proposer", getSpanProposersHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/bor/prepare-next-span", prepareNextSpanHandlerFn(cdc, cliCtx)).Methods("GET")
}

func getSpanHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// get to address
		spanID, ok := rest.ParseUint64OrReturnBadRequest(w, vars["id"])
		if !ok {
			return
		}

		res, err := cliCtx.QueryStore(bor.GetSpanKey(spanID), "bor")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// the query will return empty if there is no data in buffer
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, errors.New("No content found for requested span").Error())
			return
		}

		var span types.Span
		err = cdc.UnmarshalBinaryBare(res, &span)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := json.Marshal(&span)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

func getLatestSpanHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//
		// Get latest span start block
		//

		res, err := cliCtx.QueryStore(bor.LastSpanIDKey, "bor")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := rest.ReturnNotFoundIfNoContent(w, res); !ok {
			return
		}

		lastestSpanStart, ok := rest.ParseUint64OrReturnBadRequest(w, string(res))
		if !ok {
			return
		}

		//
		// Get latest span
		//

		res, err = cliCtx.QueryStore(bor.GetSpanKey(lastestSpanStart), "bor")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := rest.ReturnNotFoundIfNoContent(w, res); !ok {
			return
		}

		var span types.Span
		err = cdc.UnmarshalBinaryBare(res, &span)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := json.Marshal(&span)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// return result
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

func getSpanProposersHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO add span proposer logic selection here
		// for now its set as last producer sorted by address in current span

		//
		// Get latest span start block
		//

		res, err := cliCtx.QueryStore(bor.LastSpanIDKey, "bor")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := rest.ReturnNotFoundIfNoContent(w, res); !ok {
			return
		}

		lastestSpanID, ok := rest.ParseUint64OrReturnBadRequest(w, string(res))
		if !ok {
			return
		}

		//
		// Get latest span
		//

		res, err = cliCtx.QueryStore(bor.GetSpanKey(lastestSpanID), "bor")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := rest.ReturnNotFoundIfNoContent(w, res); !ok {
			return
		}

		var span types.Span
		err = cdc.UnmarshalBinaryBare(res, &span)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// selected producers
		selectedProducers := types.SortValidatorByAddress(span.SelectedProducers)

		//
		// Get span proposers
		//

		result, err := json.Marshal(&selectedProducers)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

func prepareNextSpanHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		spanID, ok := rest.ParseUint64OrReturnBadRequest(w, params.Get("span_id"))
		if !ok {
			return
		}
		startBlock, ok := rest.ParseUint64OrReturnBadRequest(w, params.Get("start_block"))
		if !ok {
			return
		}
		chainID := params.Get("chain_id")
		proposer := params.Get("proposer")

		//
		// Get span duration
		//

		// fetch duration
		res, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", borTypes.QuerierRoute, bor.QueryParams, bor.ParamSpan), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("Span duration not found ").Error())
			return
		}
		var spanDuration uint64
		if err := cliCtx.Codec.UnmarshalJSON(res, &spanDuration); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		//
		// Get ack count
		//

		// fetch ack count
		res, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", checkpointTypes.QuerierRoute, checkpoint.QueryAckCount), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("Ack not found").Error())
			return
		}

		var ackCount uint64
		if err := cliCtx.Codec.UnmarshalJSON(res, &ackCount); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		//
		// Get producer count
		//

		// fetch producer count
		res, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", borTypes.QuerierRoute, bor.QueryParams, bor.ParamProducerCount), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("Producer count not found").Error())
			return
		}

		var producerCount uint64
		if err := cliCtx.Codec.UnmarshalJSON(res, &producerCount); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		//
		// Validators
		//

		res, err = cliCtx.QueryStore(staking.CurrentValidatorSetKey, "staking")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		// the query will return empty if there is no data
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNoContent, errors.New("no content found for requested key").Error())
			return
		}
		var _validatorSet types.ValidatorSet
		err = cdc.UnmarshalBinaryBare(res, &_validatorSet)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNoContent, errors.New("unable to unmarshall binary bare").Error())
			return
		}
		var currentValidators []types.Validator

		for _, val := range _validatorSet.Validators {
			if val.IsCurrentValidator(uint64(ackCount)) {
				// append if validator is current valdiator
				currentValidators = append(currentValidators, *val)
			}
		}

		contractCaller, err := helper.NewContractCaller()
		//
		// Get last eth header processed
		//

		// fetch last eth header
		res, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", borTypes.QuerierRoute, bor.QueryParams, bor.ParamLastEthBlock), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("Producer count not found").Error())
			return
		}

		var lastEthHeader *big.Int
		if err := cliCtx.Codec.UnmarshalJSON(res, &lastEthHeader); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if lastEthHeader == nil {
			lastEthHeader = big.NewInt(0)
		}
		// fetch block header
		blockHeader, err := contractCaller.GetMainChainBlock(lastEthHeader.Add(lastEthHeader, big.NewInt(1)))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("error fetching block from mainchain").Error())
			return
		}

		selectedProducerIndices, err := bor.SelectNextProducers(blockHeader.Hash(), currentValidators, producerCount)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("error selecting producers from validator set").Error())
			return
		}

		IDToPower := make(map[uint64]uint64)
		for _, ID := range selectedProducerIndices {
			IDToPower[ID] = IDToPower[ID] + 1
		}

		var selectedProducers []types.Validator
		for _, val := range currentValidators {
			if IDToPower[val.ID.Uint64()] > 0 {
				val.Power = IDToPower[val.ID.Uint64()]
				selectedProducers = append(selectedProducers, val)
			}
		}

		// draft a propose span message
		msg := bor.NewMsgProposeSpan(
			spanID,
			types.HexToHeimdallAddress(proposer),
			startBlock,
			startBlock+spanDuration-1,
			types.ValToMinVal(currentValidators),
			types.ValToMinVal(selectedProducers),
			chainID,
		)
		result, err := json.Marshal(&msg)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, result)
	}
}
