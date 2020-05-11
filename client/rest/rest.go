package rest

import (
	"log"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/client/utils"
	"github.com/maticnetwork/heimdall/helper"
	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

//-----------------------------------------------------------------------------
// Building / Sending utilities

// WriteGenerateStdTxResponse writes response for the generate only mode.
func WriteGenerateStdTxResponse(
	w http.ResponseWriter,
	cliCtx context.CLIContext,
	br hmRest.BaseReq,
	msgs []sdk.Msg,
) {

	gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, br.GasAdjustment, client.DefaultGasAdjustment)
	if !ok {
		return
	}

	simAndExec, gas, err := client.ParseGas(br.Gas)
	if err != nil {
		hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	txBldr := authTypes.NewTxBuilder(
		helper.GetTxEncoder(cliCtx.Codec), br.AccountNumber, br.Sequence, gas, gasAdj,
		br.Simulate, br.ChainID, br.Memo, br.Fees, br.GasPrices,
	)

	if br.Simulate || simAndExec {
		if gasAdj < 0 {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, utils.ErrInvalidGasAdjustment.Error())
			return
		}

		// txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, msgs)
		// if err != nil {
		// 	hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		// 	return
		// }

		if br.Simulate {
			hmRest.WriteSimulationResponse(w, cliCtx.Codec, txBldr.Gas())
			return
		}
	}

	stdMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	output, err := cliCtx.Codec.MarshalJSON(authTypes.NewStdTx(stdMsg.Msg, nil, stdMsg.Memo))
	if err != nil {
		hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(output); err != nil {
		log.Printf("could not write response: %v", err)
	}
}
