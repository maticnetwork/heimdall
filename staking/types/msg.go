package types

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var cdc = codec.New()

//
// Validator Join
//

var _ sdk.Msg = &MsgValidatorJoin{}

type MsgValidatorJoin struct {
	From         hmTypes.HeimdallAddress `json:"from"`
	ID           hmTypes.ValidatorID     `json:"id"`
	SignerPubKey hmTypes.PubKey          `json:"pub_key"`
	TxHash       hmTypes.HeimdallHash    `json:"tx_hash"`
	LogIndex     uint64                  `json:"log_index"`
}

// NewMsgValidatorJoin creates new validator-join
func NewMsgValidatorJoin(
	from hmTypes.HeimdallAddress,
	id uint64,
	pubkey hmTypes.PubKey,
	txhash hmTypes.HeimdallHash,
	logIndex uint64,
) MsgValidatorJoin {

	return MsgValidatorJoin{
		From:         from,
		ID:           hmTypes.NewValidatorID(id),
		SignerPubKey: pubkey,
		TxHash:       txhash,
		LogIndex:     logIndex,
	}
}

func (msg MsgValidatorJoin) Type() string {
	return "validator-join"
}

func (msg MsgValidatorJoin) Route() string {
	return RouterKey
}

func (msg MsgValidatorJoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.HeimdallAddressToAccAddress(msg.From)}
}

func (msg MsgValidatorJoin) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgValidatorJoin) ValidateBasic() sdk.Error {
	if msg.ID <= 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	}

	if bytes.Equal(msg.SignerPubKey.Bytes(), helper.ZeroPubKey.Bytes()) {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid pub key %v", msg.SignerPubKey.String())
	}

	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	}

	return nil
}

//
// Stake update
//

//
// validator exit
//

var _ sdk.Msg = &MsgStakeUpdate{}

// MsgStakeUpdate represents stake update
type MsgStakeUpdate struct {
	From     hmTypes.HeimdallAddress `json:"from"`
	ID       hmTypes.ValidatorID     `json:"id"`
	TxHash   hmTypes.HeimdallHash    `json:"tx_hash"`
	LogIndex uint64                  `json:"log_index"`
}

// NewMsgStakeUpdate represents stake update
func NewMsgStakeUpdate(from hmTypes.HeimdallAddress, id uint64, txhash hmTypes.HeimdallHash, logIndex uint64) MsgStakeUpdate {
	return MsgStakeUpdate{
		From:     from,
		ID:       hmTypes.NewValidatorID(id),
		TxHash:   txhash,
		LogIndex: logIndex,
	}
}

func (msg MsgStakeUpdate) Type() string {
	return "validator-stake-update"
}

func (msg MsgStakeUpdate) Route() string {
	return RouterKey
}

func (msg MsgStakeUpdate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.HeimdallAddressToAccAddress(msg.From)}
}

func (msg MsgStakeUpdate) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgStakeUpdate) ValidateBasic() sdk.Error {
	if msg.ID <= 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	}

	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	}

	return nil
}

//
// validator update
//
var _ sdk.Msg = &MsgSignerUpdate{}

// MsgSignerUpdate signer update struct
// TODO add old signer sig check
type MsgSignerUpdate struct {
	From            hmTypes.HeimdallAddress `json:"from"`
	ID              hmTypes.ValidatorID     `json:"id"`
	NewSignerPubKey hmTypes.PubKey          `json:"pubKey"`
	TxHash          hmTypes.HeimdallHash    `json:"tx_hash"`
	LogIndex        uint64                  `json:"log_index"`
}

func NewMsgSignerUpdate(
	from hmTypes.HeimdallAddress,
	id uint64,
	pubKey hmTypes.PubKey,
	txhash hmTypes.HeimdallHash,
	logIndex uint64,
) MsgSignerUpdate {
	return MsgSignerUpdate{
		From:            from,
		ID:              hmTypes.NewValidatorID(id),
		NewSignerPubKey: pubKey,
		TxHash:          txhash,
		LogIndex:        logIndex,
	}
}

func (msg MsgSignerUpdate) Type() string {
	return "signer-update"
}

func (msg MsgSignerUpdate) Route() string {
	return RouterKey
}

func (msg MsgSignerUpdate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.HeimdallAddressToAccAddress(msg.From)}
}

func (msg MsgSignerUpdate) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgSignerUpdate) ValidateBasic() sdk.Error {
	if msg.ID <= 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	}

	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	}

	if bytes.Equal(msg.NewSignerPubKey.Bytes(), helper.ZeroPubKey.Bytes()) {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid pub key %v", msg.NewSignerPubKey.String())
	}

	return nil
}

//
// validator exit
//

var _ sdk.Msg = &MsgValidatorExit{}

type MsgValidatorExit struct {
	From     hmTypes.HeimdallAddress `json:"from"`
	ID       hmTypes.ValidatorID     `json:"id"`
	TxHash   hmTypes.HeimdallHash    `json:"tx_hash"`
	LogIndex uint64                  `json:"log_index"`
}

func NewMsgValidatorExit(from hmTypes.HeimdallAddress, id uint64, txhash hmTypes.HeimdallHash, logIndex uint64) MsgValidatorExit {
	return MsgValidatorExit{
		From:     from,
		ID:       hmTypes.NewValidatorID(id),
		TxHash:   txhash,
		LogIndex: logIndex,
	}
}

func (msg MsgValidatorExit) Type() string {
	return "validator-exit"
}

func (msg MsgValidatorExit) Route() string {
	return RouterKey
}

func (msg MsgValidatorExit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.HeimdallAddressToAccAddress(msg.From)}
}

func (msg MsgValidatorExit) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgValidatorExit) ValidateBasic() sdk.Error {
	if msg.ID <= 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid validator ID %v", msg.ID)
	}

	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	}

	return nil
}

// //
// Delegator Bond
//

var _ sdk.Msg = &MsgDelegatorBond{}

type MsgDelegatorBond struct {
	From     types.HeimdallAddress `json:"from"`
	ID       types.DelegatorID     `json:"id"`
	TxHash   types.HeimdallHash    `json:"tx_hash"`
	LogIndex uint64                `json:"log_index"`
}

// NewMsgDelegatorBond creates new delegator-bond
func NewMsgDelegatorBond(
	from hmTypes.HeimdallAddress,
	id hmTypes.DelegatorID,
	txhash hmTypes.HeimdallHash,
	logIndex uint64,
) MsgDelegatorBond {

	return MsgDelegatorBond{
		From:     from,
		ID:       id,
		TxHash:   txhash,
		LogIndex: logIndex,
	}
}

func (msg MsgDelegatorBond) Type() string {
	return "delegator-bond"
}

func (msg MsgDelegatorBond) Route() string {
	return RouterKey
}

func (msg MsgDelegatorBond) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.HeimdallAddressToAccAddress(msg.From)}
}

func (msg MsgDelegatorBond) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgDelegatorBond) ValidateBasic() sdk.Error {
	if msg.ID <= 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid delegator ID %v", msg.ID)
	}

	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	}

	return nil
}

// //
// Delegator unbond
//

var _ sdk.Msg = &MsgDelegatorUnBond{}

type MsgDelegatorUnBond struct {
	From     types.HeimdallAddress `json:"from"`
	ID       types.DelegatorID     `json:"id"`
	TxHash   types.HeimdallHash    `json:"tx_hash"`
	LogIndex uint64                `json:"log_index"`
}

// NewMsgDelegatorUnBond creates new delegator-unbond
func NewMsgDelegatorUnBond(
	from hmTypes.HeimdallAddress,
	id hmTypes.DelegatorID,
	txhash hmTypes.HeimdallHash,
	logIndex uint64,
) MsgDelegatorUnBond {

	return MsgDelegatorUnBond{
		From:     from,
		ID:       id,
		TxHash:   txhash,
		LogIndex: logIndex,
	}
}

func (msg MsgDelegatorUnBond) Type() string {
	return "delegator-unbond"
}

func (msg MsgDelegatorUnBond) Route() string {
	return RouterKey
}

func (msg MsgDelegatorUnBond) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.HeimdallAddressToAccAddress(msg.From)}
}

func (msg MsgDelegatorUnBond) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgDelegatorUnBond) ValidateBasic() sdk.Error {
	if msg.ID <= 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid delegator ID %v", msg.ID)
	}

	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	}

	return nil
}

// //
// Delegator rebond
//

var _ sdk.Msg = &MsgDelegatorReBond{}

type MsgDelegatorReBond struct {
	From     types.HeimdallAddress `json:"from"`
	ID       types.DelegatorID     `json:"id"`
	TxHash   types.HeimdallHash    `json:"tx_hash"`
	LogIndex uint64                `json:"log_index"`
}

// NewMsgDelegatorReBond creates new delegator-rebond
func NewMsgDelegatorReBond(
	from hmTypes.HeimdallAddress,
	id hmTypes.DelegatorID,
	txhash hmTypes.HeimdallHash,
	logIndex uint64,
) MsgDelegatorReBond {

	return MsgDelegatorReBond{
		From:     from,
		ID:       id,
		TxHash:   txhash,
		LogIndex: logIndex,
	}
}

func (msg MsgDelegatorReBond) Type() string {
	return "delegator-rebond"
}

func (msg MsgDelegatorReBond) Route() string {
	return RouterKey
}

func (msg MsgDelegatorReBond) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{types.HeimdallAddressToAccAddress(msg.From)}
}

func (msg MsgDelegatorReBond) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgDelegatorReBond) ValidateBasic() sdk.Error {
	if msg.ID <= 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid delegator ID %v", msg.ID)
	}

	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	}

	return nil
}

var _ sdk.Msg = &MsgDelStakeUpdate{}

// MsgDelStakeUpdate represents stake update
type MsgDelStakeUpdate struct {
	From     hmTypes.HeimdallAddress `json:"from"`
	ID       types.DelegatorID       `json:"id"`
	TxHash   hmTypes.HeimdallHash    `json:"tx_hash"`
	LogIndex uint64                  `json:"log_index"`
}

// NewMsgStakeUpdate represents stake update
func NewMsgDelStakeUpdate(from hmTypes.HeimdallAddress, id hmTypes.DelegatorID, txhash hmTypes.HeimdallHash, logIndex uint64) MsgDelStakeUpdate {
	return MsgDelStakeUpdate{
		From:     from,
		ID:       id,
		TxHash:   txhash,
		LogIndex: logIndex,
	}
}

func (msg MsgDelStakeUpdate) Type() string {
	return "delegator-stake-update"
}

func (msg MsgDelStakeUpdate) Route() string {
	return RouterKey
}

func (msg MsgDelStakeUpdate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{hmTypes.HeimdallAddressToAccAddress(msg.From)}
}

func (msg MsgDelStakeUpdate) GetSignBytes() []byte {
	b, err := cdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgDelStakeUpdate) ValidateBasic() sdk.Error {
	if msg.ID <= 0 {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid delegator ID %v", msg.ID)
	}

	if msg.From.Empty() {
		return hmCommon.ErrInvalidMsg(hmCommon.DefaultCodespace, "Invalid proposer %v", msg.From.String())
	}

	return nil
}
