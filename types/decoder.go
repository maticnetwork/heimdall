package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RLPTxDecoder decodes the txBytes to a BaseTx
func RLPTxDecoder(pulp *Pulp) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, sdk.Error) {
		msg, err := pulp.DecodeBytes(txBytes)
		if err != nil {
			return nil, sdk.ErrTxDecode(err.Error())
		}

		return &BaseTx{
			Msg: msg.(sdk.Msg),
		}, nil
	}
}
