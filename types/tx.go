package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TxSignature represents tx signature
type TxSignature []byte

// assertion
var _ sdk.Tx = BaseTx{}

// BaseTx represents base tx tendermint needs
type BaseTx struct {
	Msg        sdk.Msg
	Signatures TxSignature
}

// NewBaseTx drafts BaseTx with messages
func NewBaseTx(msg sdk.Msg, sig TxSignature) BaseTx {
	return BaseTx{
		Msg: msg,
	}
}

// GetMsgs returns array of messages
func (tx BaseTx) GetMsgs() []sdk.Msg {
	return []sdk.Msg{tx.Msg}
}

// ValidateBasic does a simple and lightweight validation check that doesn't
// require access to any other information.
func (tx BaseTx) ValidateBasic() sdk.Error {
	return nil
}
