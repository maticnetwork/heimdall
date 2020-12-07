package keeper

import (
	"context"
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/staking/types"
)

type msgServer struct {
	Keeper
	contractCaller helper.IContractCaller
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, contractCaller helper.IContractCaller) types.MsgServer {
	return &msgServer{Keeper: keeper, contractCaller: contractCaller}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) ValidatorJoin(goCtx context.Context, msg *types.MsgValidatorJoin) (*types.MsgValidatorJoinResponse, error) {
	// k.Logger(ctx).Info("Handling new validator join", "msg", msg)

	ctx := sdk.UnwrapSDKContext(goCtx)
	// params := k.chainKeeper.GetParams(ctx)
	// chainParams := params.ChainParams

	// // get main tx receipt
	// receipt, err := k.contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash(), params.TxConfirmationTime)
	// if err != nil || receipt == nil {
	// 	return nil, hmCommon.ErrWaitForConfirmation
	// }

	// // decode validator join event
	// eventLog, err := k.contractCaller.DecodeValidatorJoinEvent(chainParams.StakingInfoAddress.EthAddress(), receipt, msg.LogIndex)
	// if err != nil || eventLog == nil {
	// 	return nil, hmCommon.ErrInvalidMsg
	// }

	// Generate PubKey from Pubkey in message and signer
	// pubkey := msg.SignerPubKey
	// signer := pubkey.Address()

	// // check signer pubkey in message corresponds
	// if !bytes.Equal(pubkey.Bytes()[1:], eventLog.SignerPubkey) {
	// 	k.Logger(ctx).Error(
	// 		"Signer Pubkey does not match",
	// 		"msgValidator", pubkey.String(),
	// 		"mainchainValidator", hmTypes.BytesToHexBytes(eventLog.SignerPubkey),
	// 	)
	// 	return nil, hmCommon.ErrValSignerPubKeyMismatch
	// }

	// // check signer corresponding to pubkey matches signer from event
	// if !bytes.Equal(signer.Bytes(), eventLog.Signer.Bytes()) {
	// 	k.Logger(ctx).Error(
	// 		"Signer Address from Pubkey does not match",
	// 		"Validator", signer.String(),
	// 		"mainchainValidator", eventLog.Signer.Hex(),
	// 	)
	// 	return nil, hmCommon.ErrValSignerMismatch
	// }

	// // check msg id
	// if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
	// 	k.Logger(ctx).Error("ID in message doesn't match with id in log", "msgId", msg.ID, "validatorIdFromTx", eventLog.ValidatorId)
	// 	return nil, hmCommon.ErrInvalidMsg
	// }

	// Check if validator has been validator before
	if _, ok := k.GetSignerFromValidatorID(ctx, msg.ID); ok {
		k.Logger(ctx).Error("Validator has been validator before, cannot join with same ID", "validatorId", msg.ID)
		return nil, hmCommon.ErrValidatorAlreadyJoined
	}

	// // get validator by signer
	// checkVal, err := k.GetValidatorInfo(ctx, signer.Bytes())
	// // if err == nil || bytes.Equal(checkVal.Signer.Bytes(), signer.Bytes()) {
	// // 	return nil, hmCommon.ErrValidatorAlreadyJoined
	// // }

	// get voting power from amount
	_, err := helper.GetPowerFromAmount(msg.Amount.BigInt())
	if err != nil {
		return nil, hmCommon.ErrInvalidMsg
	}

	// sequence id
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// // create new validator
	// newValidator := hmTypes.Validator{
	// 	ID:          msg.ID,
	// 	StartEpoch:  eventLog.ActivationEpoch.Uint64(),
	// 	EndEpoch:    0,
	// 	VotingPower: votingPower.Int64(),
	// 	PubKey:      pubkey,
	// 	Signer:      hmTypes.BytesToHeimdallAddress(signer.Bytes()),
	// 	LastUpdated: "",
	// }

	// //sequence id
	// sequence := new(big.Int).Mul(receipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	// sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// // check if incoming tx is older
	// if k.HasStakingSequence(ctx, sequence.String()) {
	// 	k.Logger(ctx).Error("Older invalid tx found")
	// 	return nil, hmCommon.ErrOldTx
	// }

	// // update last updated
	// newValidator.LastUpdated = sequence.String()

	// // add validator to store
	// k.Logger(ctx).Debug("Adding new validator to state", "validator", newValidator.String())
	// err = k.AddValidator(ctx, newValidator)
	// if err != nil {
	// 	k.Logger(ctx).Error("Unable to add validator to state", "error", err, "validator", newValidator.String())
	// 	return nil, hmCommon.ErrValidatorSave
	// }

	// save staking sequence
	// k.SetStakingSequence(ctx, sequence.String())

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return nil, hmCommon.ErrOldTx
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeValidatorJoin,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(msg.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeySigner, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return &types.MsgValidatorJoinResponse{}, nil

}

func (k msgServer) StakeUpdate(goCtx context.Context, msg *types.MsgStakeUpdate) (*types.MsgStakeUpdateResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.MsgStakeUpdateResponse{}, nil
}

func (k msgServer) SignerUpdate(goCtx context.Context, msg *types.MsgSignerUpdate) (*types.MsgSignerUpdateResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.MsgSignerUpdateResponse{}, nil

}

func (k msgServer) ValidatorExit(goCtx context.Context, msg *types.MsgValidatorExit) (*types.MsgValidatorExitResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.MsgValidatorExitResponse{}, nil
}
