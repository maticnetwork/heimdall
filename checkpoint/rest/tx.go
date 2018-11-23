package rest

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/mux"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"

	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/types"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/checkpoint/new",
		newCheckpointHandler(cliCtx),
	).Methods("POST")
}

type HeaderBlock struct {
	Proposer   string `json:"proposer"`
	RootHash   string `json:"rootHash"`
	StartBlock uint64 `json:"startBlock"`
	EndBlock   uint64 `json:"endBlock"`
}

func newCheckpointHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m HeaderBlock

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

		txBytes, err := CreateTxBytes(m)
		if err != nil {
			RestLogger.Error("Unable to create txBytes", "endBlock", m.EndBlock, "startBlock", m.StartBlock, "rootHash", m.RootHash)
		}

		resp, err := SendTendermintRequest(cliCtx, txBytes)
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

func SendTendermintRequest(cliCtx context.CLIContext, txBytes []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	RestLogger.Info("Broadcasting tx bytes to Tendermint", "txBytes", hex.EncodeToString(txBytes))
	return cliCtx.BroadcastTx(txBytes)
}

func CreateTxBytes(m HeaderBlock) ([]byte, error) {
	msg := checkpoint.NewMsgCheckpointBlock(
		common.HexToAddress(m.Proposer),
		m.StartBlock,
		m.EndBlock,
		common.HexToHash(m.RootHash))

	tx := types.NewBaseTx(msg)

	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		RestLogger.Error("Error generating TX Bytes", "error", err)

		return []byte(""), err
	}
	return txBytes, nil
}
