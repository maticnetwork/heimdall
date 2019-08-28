package rest

import (
	"log"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

//-----------------------------------------------------------------------------
// Building / Sending utilities

// WriteGenerateStdTxResponse writes response for the generate only mode.
func WriteGenerateStdTxResponse(
	w http.ResponseWriter,
	cliCtx context.CLIContext,
	br rest.BaseReq,
	msgs []sdk.Msg,
) {

	gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, br.GasAdjustment, client.DefaultGasAdjustment)
	if !ok {
		return
	}

	simAndExec, gas, err := client.ParseGas(br.Gas)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	txBldr := authTypes.NewTxBuilder(
		utils.GetTxEncoder(cliCtx.Codec), br.AccountNumber, br.Sequence, gas, gasAdj,
		br.Simulate, br.ChainID, br.Memo, br.Fees, br.GasPrices,
	)

	if br.Simulate || simAndExec {
		if gasAdj < 0 {
			rest.WriteErrorResponse(w, http.StatusBadRequest, client.ErrInvalidGasAdjustment.Error())
			return
		}

		// txBldr, err = utils.EnrichWithGas(txBldr, cliCtx, msgs)
		// if err != nil {
		// 	rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		// 	return
		// }

		if br.Simulate {
			rest.WriteSimulationResponse(w, cliCtx.Codec, txBldr.Gas())
			return
		}
	}

	stdMsg, err := txBldr.BuildSignMsg(msgs)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	output, err := cliCtx.Codec.MarshalJSON(authTypes.NewStdTx(stdMsg.Msg, nil, stdMsg.Memo))
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(output); err != nil {
		log.Printf("could not write response: %v", err)
	}
	return
}
