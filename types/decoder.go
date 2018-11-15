package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/maticnetwork/heimdall/checkpoint"
)

// RLPTxDecoder decodes the txBytes to a BaseTx
func RLPTxDecoder() sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, sdk.Error) {
		var tx = checkpoint.BaseTx{}
		err := rlp.DecodeBytes(txBytes, &tx)
		if err != nil {
			//todo create own error
			return nil, sdk.ErrTxDecode(err.Error())
		}

		return tx, nil

	}
}
