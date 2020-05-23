package topup

import (
	"bytes"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/auth"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/common"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

// NewSideTxHandler returns a side handler for "topup" type messages.
func NewSideTxHandler(k Keeper, contractCaller helper.IContractCaller) hmTypes.SideTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgTopup:
			return SideHandleMsgTopup(ctx, k, msg, contractCaller)
		default:
			return abci.ResponseDeliverSideTx{
				Code: uint32(sdk.CodeUnknownRequest),
			}
		}
	}
}

// NewPostTxHandler returns a side handler for "bank" type messages.
func NewPostTxHandler(k Keeper, contractCaller helper.IContractCaller) hmTypes.PostTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgTopup:
			return PostHandleMsgTopup(ctx, k, msg, sideTxResult)
		default:
			return sdk.ErrUnknownRequest("Unrecognized topup msg type").Result()
		}
	}
}

// SideHandleMsgTopup handles MsgTopup message for external call
func SideHandleMsgTopup(ctx sdk.Context, k Keeper, msg types.MsgTopup, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {

	k.Logger(ctx).Debug("✅ Validating External call for topup msg",
		"txHash", hmTypes.BytesToHeimdallHash(msg.TxHash.Bytes()),
		"logIndex", uint64(msg.LogIndex),
		"blockNumber", msg.BlockNumber,
	)

	// chainManager params
	params := k.chainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash(), params.MainchainTxConfirmations)
	if err != nil || receipt == nil {
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeWaitFrConfirmation)
	}

	// get event log for topup
	eventLog, err := contractCaller.DecodeValidatorTopupFeesEvent(chainParams.StakingInfoAddress.EthAddress(), receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeErrDecodeEvent)
	}

	if receipt.BlockNumber.Uint64() != msg.BlockNumber {
		k.Logger(ctx).Error("BlockNumber in message doesn't match blocknumber in receipt", "MsgBlockNumber", msg.BlockNumber, "ReceiptBlockNumber", receipt.BlockNumber.Uint64)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	if !bytes.Equal(eventLog.User.Bytes(), msg.User.Bytes()) {
		k.Logger(ctx).Error(
			"User Address from event does not match with Msg user",
			"EventUser", eventLog.User.String(),
			"MsgUser", msg.User.String(),
		)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	if eventLog.Fee.Cmp(msg.Fee.BigInt()) != 0 {
		k.Logger(ctx).Error("Fee in message doesn't match Fee in event logs", "MsgFee", msg.Fee, "FeeFromEvent", eventLog.Fee)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	k.Logger(ctx).Debug("✅ Succesfully validated External call for topup msg")
	result.Result = abci.SideTxResultType_Yes
	return
}

func PostHandleMsgTopup(ctx sdk.Context, k Keeper, msg types.MsgTopup, sideTxResult abci.SideTxResultType) sdk.Result {

	// Skip handler if topup is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		k.Logger(ctx).Debug("Skipping new topup since side-tx didn't get yes votes")
		return common.ErrSideTxValidation(k.Codespace()).Result()
	}

	// check if incoming tx is older
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	if k.HasTopupSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("Persisting topup state", "sideTxResult", sideTxResult)

	// use event log user
	user := msg.User

	// create topup amount
	topupAmount := sdk.Coins{sdk.Coin{Denom: authTypes.FeeToken, Amount: msg.Fee}}

	// increase coins in account
	if _, err := k.bk.AddCoins(ctx, user, topupAmount); err != nil {
		k.Logger(ctx).Error("Error while adding coins to user", "user", user, "topupAmount", topupAmount, "error", err)
		return err.Result()
	}

	// transfer fees to sender (proposer)
	if err := k.bk.SendCoins(ctx, user, msg.FromAddress, auth.DefaultFeeWantedPerTx); err != nil {
		return err.Result()
	}

	k.Logger(ctx).Debug("Persisted topup state for", "user", user, "topupAmount", topupAmount.String())

	// save topup
	k.SetTopupSequence(ctx, sequence.String())

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTopup,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                  // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),             // result
			sdk.NewAttribute(types.AttributeKeySender, msg.FromAddress.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.User.String()),
			sdk.NewAttribute(types.AttributeKeyTopupAmount, msg.Fee.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
