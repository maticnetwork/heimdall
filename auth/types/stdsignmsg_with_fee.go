package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//__________________________________________________________

// StdSignDoc is replay-prevention structure.
// It includes the result of msg.GetSignBytes(),
// as well as the ChainID (prevent cross chain replay)
// and the Sequence numbers for each signature (prevent
// inchain replay and enforce tx ordering per account).
type StdSignDocWithFee struct {
	ChainID       string          `json:"chain_id" yaml:"chain_id"`
	AccountNumber uint64          `json:"account_number" yaml:"account_number"`
	Fee           json.RawMessage `json:"fee" yaml:"fee"`
	Sequence      uint64          `json:"sequence" yaml:"sequence"`
	Msg           json.RawMessage `json:"msg" yaml:"msg"`
	Memo          string          `json:"memo" yaml:"memo"`
}

// StdSignBytes returns the bytes to sign for a transaction.
func StdSignBytesWithFee(chainID string, accnum uint64, sequence uint64, fee StdFee, msg sdk.Msg, memo string) []byte {
	msgsBytes := json.RawMessage(msg.GetSignBytes())
	bz, err := ModuleCdc.MarshalJSON(StdSignDocWithFee{
		AccountNumber: accnum,
		ChainID:       chainID,
		Memo:          memo,
		Fee:           json.RawMessage(fee.Bytes()),
		Msg:           msgsBytes,
		Sequence:      sequence,
	})
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

// StdSignMsg is a convenience structure for passing along
// a Msg with the other requirements for a StdSignDoc before
// it is signed. For use in the CLI.
type StdSignMsgWithFee struct {
	ChainID       string  `json:"chain_id" yaml:"chain_id"`
	AccountNumber uint64  `json:"account_number" yaml:"account_number"`
	Sequence      uint64  `json:"sequence" yaml:"sequence"`
	Fee           StdFee  `json:"fee" yaml:"fee"`
	Msg           sdk.Msg `json:"msg" yaml:"msg"`
	Memo          string  `json:"memo" yaml:"memo"`
}

// Bytes returns message bytes
func (msg StdSignMsgWithFee) Bytes() []byte {
	return StdSignBytesWithFee(msg.ChainID, msg.AccountNumber, msg.Sequence, msg.Fee, msg.Msg, msg.Memo)
}
