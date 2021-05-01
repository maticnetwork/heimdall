package topup

import (
	"bytes"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/common"
	hmCommon "github.com/maticnetwork/heimdall/common"

	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/topup/keeper"
	"github.com/maticnetwork/heimdall/x/topup/types"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

// NewSideTxHandler returns a side handler for "topup" type messages.
func NewSideTxHandler(k keeper.Keeper, contractCaller helper.IContractCaller) hmTypes.SideTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgTopup:
			return SideHandleMsgTopup(ctx, k, *msg, contractCaller)
		default:
			return abci.ResponseDeliverSideTx{
				Code: sdkerrors.ErrUnknownRequest.ABCICode(),
			}
		}
	}
}

// NewPostTxHandler returns a side handler for "bank" type messages.
func NewPostTxHandler(k keeper.Keeper, contractCaller helper.IContractCaller) hmTypes.PostTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgTopup:
			return PostHandleMsgTopup(ctx, k, *msg, sideTxResult)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// SideHandleMsgTopup handles MsgTopup message for external call
func SideHandleMsgTopup(ctx sdk.Context, k keeper.Keeper, msg types.MsgTopup, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {

	k.Logger(ctx).Debug("✅ Validating External call for topup msg",
		"txHash", hmCommonTypes.BytesToHeimdallHash([]byte(msg.TxHash)),
		"logIndex", uint64(msg.LogIndex),
		"blockNumber", msg.BlockNumber,
	)

	// chainManager params
	params := k.ChainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(hmCommonTypes.HexToHeimdallHash(msg.TxHash).EthHash(), params.MainchainTxConfirmations)
	if err != nil || receipt == nil {
		return hmCommon.ErrorSideTx(common.ErrWaitForConfirmation)
	}

	// get event log for topup
	//var stakingAddress [20]byte
	//copy(stakingAddress[:], chainParams.StakingInfoAddress)

	accountAddr, _ := sdk.AccAddressFromHex(chainParams.StakingInfoAddress)
	eventLog, err := contractCaller.DecodeValidatorTopupFeesEvent(accountAddr, receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrorSideTx(common.ErrDecodeEvent)
	}

	if receipt.BlockNumber.Uint64() != msg.BlockNumber {
		k.Logger(ctx).Error("BlockNumber in message doesn't match blocknumber in receipt", "MsgBlockNumber", msg.BlockNumber, "ReceiptBlockNumber", receipt.BlockNumber.Uint64)
		return hmCommon.ErrorSideTx(common.ErrInvalidMsg)
	}

	accountAddr1, _ := sdk.AccAddressFromHex(msg.User)
	if !bytes.Equal(eventLog.User.Bytes(), accountAddr1.Bytes()) {
		k.Logger(ctx).Error(
			"User Address from event does not match with Msg user",
			"EventUser", eventLog.User,
			"MsgUser", msg.User,
		)
		return hmCommon.ErrorSideTx(common.ErrInvalidMsg)
	}

	if eventLog.Fee.Cmp(msg.Fee.BigInt()) != 0 {
		k.Logger(ctx).Error("Fee in message doesn't match Fee in event logs", "MsgFee", msg.Fee, "FeeFromEvent", eventLog.Fee)
		return hmCommon.ErrorSideTx(common.ErrInvalidMsg)
	}

	k.Logger(ctx).Debug("✅ Succesfully validated External call for topup msg")
	result.Result = tmprototypes.SideTxResultType_YES
	return
}

func PostHandleMsgTopup(ctx sdk.Context, k keeper.Keeper, msg types.MsgTopup, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {

	// Skip handler if topup is not approved
	if sideTxResult != tmprototypes.SideTxResultType_YES {
		k.Logger(ctx).Debug("Skipping new topup since side-tx didn't get yes votes")
		return nil, hmCommon.ErrSideTxValidation
	}

	// check if incoming tx is older
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	if k.HasTopupSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return nil, hmCommon.ErrOldTx
	}

	k.Logger(ctx).Debug("Persisting topup state", "sideTxResult", sideTxResult)

	// use event log user
	user := msg.User

	// create topup amount
	topupAmount := sdk.Coins{sdk.Coin{Denom: types.FeeToken, Amount: *msg.Fee}}

	// increase coins in account
	userAddr, _ := sdk.AccAddressFromHex(user)
	if err := k.Bk.AddCoins(ctx, userAddr, topupAmount); err != nil {
		k.Logger(ctx).Error("Error while adding coins to user", "user", user, "topupAmount", topupAmount, "error", err)
		return nil, err
	}

	// transfer fees to sender (proposer)
	fromAddr, _ := sdk.AccAddressFromHex(msg.FromAddress)
	if err := k.Bk.SendCoins(ctx, fromAddr, userAddr, topupAmount); err != nil {
		return nil, err
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
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                        // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                      // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmCommonTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),                   // result
			sdk.NewAttribute(types.AttributeKeySender, fromAddr.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, userAddr.String()),
			sdk.NewAttribute(types.AttributeKeyTopupAmount, msg.Fee.String()),
		),
	})

	return &sdk.Result{
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}
