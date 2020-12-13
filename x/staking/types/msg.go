package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var cdc = codec.NewLegacyAmino()

//
// Validator Join
//

var _ sdk.Msg = &MsgValidatorJoin{}

// NewMsgValidatorJoin creates new validator-join
func NewMsgValidatorJoin(
	from string,
	id uint64,
	activationEpoch uint64,
	amount uint64,
	pubkey string,
	txhash string,
	logIndex uint64,
	blockNumber uint64,
	nonce uint64,
) MsgValidatorJoin {

	return MsgValidatorJoin{
		From:            from,
		ID:              id,
		ActivationEpoch: activationEpoch,
		Amount:          amount,
		SignerPubKey:    pubkey,
		TxHash:          txhash,
		LogIndex:        logIndex,
		BlockNumber:     blockNumber,
		Nonce:           nonce,
	}
}

func (msg MsgValidatorJoin) Type() string {
	return "validator-join"
}

func (msg MsgValidatorJoin) Route() string {
	return RouterKey
}

func (msg MsgValidatorJoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{[]byte(msg.From)}
}

func (msg MsgValidatorJoin) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgValidatorJoin) ValidateBasic() error {
	// if msg.ID == 0 {
	// 	return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	// }

	// if bytes.Equal(msg.SignerPubKey.Bytes(), helper.ZeroPubKey.Bytes()) {
	// 	return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid pub key %v", msg.SignerPubKey.String())
	// }

	// if msg.From.Empty() {
	// 	return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	// }

	return nil
}

// GetTxHash Returns tx hash
func (msg MsgValidatorJoin) GetTxHash() string {
	return msg.TxHash
}

// GetLogIndex Returns log index
func (msg MsgValidatorJoin) GetLogIndex() uint64 {
	return msg.LogIndex
}

// GetSideSignBytes returns side sign bytes
func (msg MsgValidatorJoin) GetSideSignBytes() []byte {
	return nil
}

// GetNonce Returns nonce
func (msg MsgValidatorJoin) GetNonce() uint64 {
	return msg.Nonce
}

//
// Stake update
//

var _ sdk.Msg = &MsgStakeUpdate{}

// NewMsgStakeUpdate represents stake update
func NewMsgStakeUpdate(
	from string,
	id uint64,
	newAmount uint64,
	txhash string,
	logIndex uint64,
	blockNumber uint64,
	nonce uint64,
) MsgStakeUpdate {
	return MsgStakeUpdate{
		From:        from,
		ID:          id,
		NewAmount:   newAmount,
		TxHash:      txhash,
		LogIndex:    logIndex,
		BlockNumber: blockNumber,
		Nonce:       nonce,
	}
}

func (msg MsgStakeUpdate) Type() string {
	return "validator-stake-update"
}

func (msg MsgStakeUpdate) Route() string {
	return RouterKey
}

func (msg MsgStakeUpdate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{[]byte(msg.From)}
}

func (msg MsgStakeUpdate) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgStakeUpdate) ValidateBasic() error {
	// if msg.ID == 0 {
	// 	return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	// }

	// if msg.From.Empty() {
	// 	return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	// }

	return nil
}

// GetTxHash Returns tx hash
func (msg MsgStakeUpdate) GetTxHash() string {
	return msg.TxHash
}

// GetLogIndex Returns log index
func (msg MsgStakeUpdate) GetLogIndex() uint64 {
	return msg.LogIndex
}

// GetSideSignBytes returns side sign bytes
func (msg MsgStakeUpdate) GetSideSignBytes() []byte {
	return nil
}

// GetNonce Returns nonce
func (msg MsgStakeUpdate) GetNonce() uint64 {
	return msg.Nonce
}

//
// validator update
//
var _ sdk.Msg = &MsgSignerUpdate{}

func NewMsgSignerUpdate(
	from string,
	id uint64,
	pubKey string,
	txhash string,
	logIndex uint64,
	blockNumber uint64,
	nonce uint64,
) MsgSignerUpdate {
	return MsgSignerUpdate{
		From:            from,
		ID:              id,
		NewSignerPubKey: pubKey,
		TxHash:          txhash,
		LogIndex:        logIndex,
		BlockNumber:     blockNumber,
		Nonce:           nonce,
	}
}

func (msg MsgSignerUpdate) Type() string {
	return "signer-update"
}

func (msg MsgSignerUpdate) Route() string {
	return RouterKey
}

func (msg MsgSignerUpdate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{[]byte(msg.From)}
}

func (msg MsgSignerUpdate) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgSignerUpdate) ValidateBasic() error {
	// if msg.ID == 0 {
	// 	return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	// }

	// if msg.From.Empty() {
	// 	return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	// }

	// if bytes.Equal(msg.NewSignerPubKey.Bytes(), helper.ZeroPubKey.Bytes()) {
	// 	return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid pub key %v", msg.NewSignerPubKey.String())
	// }

	return nil
}

// GetTxHash Returns tx hash
func (msg MsgSignerUpdate) GetTxHash() string {
	return msg.TxHash
}

// GetLogIndex Returns log index
func (msg MsgSignerUpdate) GetLogIndex() uint64 {
	return msg.LogIndex
}

// GetSideSignBytes returns side sign bytes
func (msg MsgSignerUpdate) GetSideSignBytes() []byte {
	return nil
}

// GetNonce Returns nonce
func (msg MsgSignerUpdate) GetNonce() uint64 {
	return msg.Nonce
}

//
// validator exit
//

var _ sdk.Msg = &MsgValidatorExit{}

func NewMsgValidatorExit(
	from string,
	id uint64,
	deactivationEpoch uint64,
	txhash string,
	logIndex uint64,
	blockNumber uint64,
	nonce uint64,
) MsgValidatorExit {
	return MsgValidatorExit{
		From:              from,
		ID:                id,
		DeactivationEpoch: deactivationEpoch,
		TxHash:            txhash,
		LogIndex:          logIndex,
		BlockNumber:       blockNumber,
		Nonce:             nonce,
	}
}

func (msg MsgValidatorExit) Type() string {
	return "validator-exit"
}

func (msg MsgValidatorExit) Route() string {
	return RouterKey
}

func (msg MsgValidatorExit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{[]byte(msg.From)}
}

func (msg MsgValidatorExit) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgValidatorExit) ValidateBasic() error {
	// if msg.ID == 0 {
	// 	return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	// }

	// if msg.From.Empty() {
	// 	return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	// }

	return nil
}

// GetTxHash Returns tx hash
func (msg MsgValidatorExit) GetTxHash() string {
	return msg.TxHash
}

// GetLogIndex Returns log index
func (msg MsgValidatorExit) GetLogIndex() uint64 {
	return msg.LogIndex
}

// GetSideSignBytes returns side sign bytes
func (msg MsgValidatorExit) GetSideSignBytes() []byte {
	return nil
}

// GetNonce Returns nonce
func (msg MsgValidatorExit) GetNonce() uint64 {
	return msg.Nonce
}
