package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommon "github.com/maticnetwork/heimdall/types/common"
)

var cdc = codec.NewLegacyAmino()

//
// Validator Join
//

var _ sdk.Msg = &MsgValidatorJoin{}

// NewMsgValidatorJoin creates new validator-join
func NewMsgValidatorJoin(
	from hmCommon.HeimdallAddress,
	id uint64,
	activationEpoch uint64,
	amount sdk.Int,
	pubkey hmCommon.PubKey,
	txhash hmCommon.HeimdallHash,
	logIndex uint64,
	blockNumber uint64,
	nonce uint64,
) MsgValidatorJoin {

	return MsgValidatorJoin{
		From:            from,
		ID:              hmTypes.NewValidatorID(id),
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
	return []sdk.AccAddress{hmCommon.HeimdallAddressToAccAddress(msg.From)}
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
func (msg MsgValidatorJoin) GetTxHash() hmCommon.HeimdallHash {
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
func NewMsgStakeUpdate(from hmCommon.HeimdallAddress, id uint64, newAmount sdk.Int, txhash hmCommon.HeimdallHash, logIndex uint64, blockNumber uint64, nonce uint64) MsgStakeUpdate {
	return MsgStakeUpdate{
		From:        from,
		ID:          hmTypes.NewValidatorID(id),
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
	return []sdk.AccAddress{hmCommon.HeimdallAddressToAccAddress(msg.From)}
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
func (msg MsgStakeUpdate) GetTxHash() hmCommon.HeimdallHash {
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
	from hmCommon.HeimdallAddress,
	id uint64,
	pubKey hmCommon.PubKey,
	txhash hmCommon.HeimdallHash,
	logIndex uint64,
	blockNumber uint64,
	nonce uint64,
) MsgSignerUpdate {
	return MsgSignerUpdate{
		From:            from,
		ID:              hmTypes.NewValidatorID(id),
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
	return []sdk.AccAddress{hmCommon.HeimdallAddressToAccAddress(msg.From)}
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
func (msg MsgSignerUpdate) GetTxHash() hmCommon.HeimdallHash {
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

func NewMsgValidatorExit(from hmCommon.HeimdallAddress, id uint64, deactivationEpoch uint64, txhash hmCommon.HeimdallHash, logIndex uint64, blockNumber uint64, nonce uint64) MsgValidatorExit {
	return MsgValidatorExit{
		From:              from,
		ID:                hmTypes.NewValidatorID(id),
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
	return []sdk.AccAddress{hmCommon.HeimdallAddressToAccAddress(msg.From)}
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
func (msg MsgValidatorExit) GetTxHash() hmCommon.HeimdallHash {
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
