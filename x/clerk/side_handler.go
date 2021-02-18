package clerk

import (
	"bytes"
	"math/big"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"

	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/x/clerk/keeper"
	"github.com/maticnetwork/heimdall/x/clerk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

// NewSideTxHandler returns a side handler for "topup" type messages.
func NewSideTxHandler(k keeper.Keeper, contractCaller helper.IContractCaller) hmTypes.SideTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgEventRecordRequest:
			return SideHandleMsgEventRecord(ctx, k, *msg, contractCaller)
		default:
			return abci.ResponseDeliverSideTx{
				Code: uint32(6), // TODO should be changed like `sdk.CodeUnknownRequest`
			}
		}
	}
}

// NewPostTxHandler returns a side handler for "bank" type messages.
func NewPostTxHandler(k keeper.Keeper, contractCaller helper.IContractCaller) hmTypes.PostTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgEventRecordRequest:
			return PostHandleMsgEventRecord(ctx, k, *msg, sideTxResult)
		default:
			return nil, sdkerrors.ErrUnknownRequest
		}
	}
}

func SideHandleMsgEventRecord(
	ctx sdk.Context,
	k keeper.Keeper,
	msg types.MsgEventRecordRequest,
	contractCaller helper.IContractCaller,
) (result abci.ResponseDeliverSideTx) {

	k.Logger(ctx).Debug("✅ Validating External call for clerk msg",
		"txHash", msg.TxHash,
		"logIndex", uint64(msg.LogIndex),
		"blockNumber", msg.BlockNumber,
	)

	// chainManager params
	params := k.ChainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// get confirmed tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(hmCommonTypes.HexToHeimdallHash(msg.TxHash).EthHash(), params.MainchainTxConfirmations)
	if receipt == nil || err != nil {
		return hmCommon.ErrorSideTx(hmCommon.ErrWaitForConfirmation)
	}

	// get event log for topup
	stakingSenderAddress, _ := sdk.AccAddressFromHex(chainParams.StateSenderAddress)
	eventLog, err := contractCaller.DecodeStateSyncedEvent(stakingSenderAddress, receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrorSideTx(hmCommon.ErrWaitForConfirmation)
	}

	if receipt.BlockNumber.Uint64() != msg.BlockNumber {
		k.Logger(ctx).Error(
			"BlockNumber in message doesn't match blocknumber in receipt", "MsgBlockNumber",
			msg.BlockNumber,
			"ReceiptBlockNumber",
			receipt.BlockNumber.Uint64(),
		)
		return hmCommon.ErrorSideTx(hmCommon.ErrInvalidMsg)
	}

	// check if message and event log matches
	if eventLog.Id.Uint64() != msg.Id {
		k.Logger(ctx).Error("ID in message doesn't match with id in log", "msgId", msg.Id, "stateIdFromTx", eventLog.Id)
		return hmCommon.ErrorSideTx(hmCommon.ErrInvalidMsg)
	}

	// TODO - check this
	if !strings.EqualFold(eventLog.ContractAddress.String(), msg.ContractAddress) {
		k.Logger(ctx).Error(
			"ContractAddress from event does not match with Msg ContractAddress",
			"EventContractAddress", eventLog.ContractAddress.String(),
			"MsgContractAddress", msg.ContractAddress,
		)
		return hmCommon.ErrorSideTx(hmCommon.ErrInvalidMsg)
	}

	if !bytes.Equal(eventLog.Data, msg.Data) {
		k.Logger(ctx).Error(
			"Data from event does not match with Msg Data",
			"EventData", hmTypes.BytesToHexBytes(eventLog.Data),
			"MsgData", hmTypes.BytesToHexBytes(msg.Data),
		)
		return hmCommon.ErrorSideTx(hmCommon.ErrInvalidMsg)
	}

	result.Result = tmprototypes.SideTxResultType_YES
	return
}

func PostHandleMsgEventRecord(
	ctx sdk.Context,
	k keeper.Keeper,
	msg types.MsgEventRecordRequest,
	sideTxResult tmprototypes.SideTxResultType,
) (*sdk.Result, error) {

	// Skip handler if clerk is not approved
	if sideTxResult != tmprototypes.SideTxResultType_YES {
		k.Logger(ctx).Debug("Skipping new clerk since side-tx didn't get yes votes")
		return nil, hmCommon.ErrSideTxValidation
	}

	// check for replay
	if k.HasEventRecord(ctx, msg.Id) {
		k.Logger(ctx).Debug("Skipping new clerk record as it's already processed")
		return nil, hmCommon.ErrOldTx
	}

	k.Logger(ctx).Debug("Persisting clerk state", "sideTxResult", sideTxResult)

	// sequence id
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// create event record
	txHash := hmCommonTypes.HexToHeimdallHash(msg.TxHash)
	contractAddress, err := sdk.AccAddressFromHex(msg.ContractAddress)
	if err != nil {
		return nil, err
	}
	record := types.NewEventRecord(
		txHash,
		msg.LogIndex,
		msg.Id,
		contractAddress,
		msg.Data,
		msg.ChainId,
		ctx.BlockTime(),
	)

	// save event into state
	if err := k.SetEventRecord(ctx, record); err != nil {
		k.Logger(ctx).Error("Unable to update event record", "error", err, "id", msg.Id)
		return nil, hmCommon.ErrEventUpdate
	}

	// save record sequence
	k.SetRecordSequence(ctx, sequence.String())

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRecord,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                        // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                      // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmCommonTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(types.AttributeKeyRecordTxLogIndex, strconv.FormatUint(msg.LogIndex, 10)),
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()), // result
			sdk.NewAttribute(types.AttributeKeyRecordID, strconv.FormatUint(msg.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyRecordContract, msg.ContractAddress),
		),
	})

	return &sdk.Result{}, nil
}
