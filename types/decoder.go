package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// RLPTxDecoder decodes the txBytes to a BaseTx
func RLPTxDecoder() sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, sdk.Error) {
		var tx = BaseTx{}
		err := rlp.DecodeBytes(txBytes, &tx)
		if err != nil {
			return nil, sdk.ErrTxDecode(err.Error())
		}
		return tx, nil

	}
}
