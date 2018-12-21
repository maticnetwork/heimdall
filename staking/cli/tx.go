package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/types"
	"encoding/hex"
	"encoding/json"
)

// send validator join transaction
func GetValidatorJoinTx(cdc *codec.Codec) *cobra.Command  {
	cmd:=&cobra.Command{
		Use:   "validator-join",
		Short: "Join Heimdall as a validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			validatorStr:= viper.GetString(FlagValidatorAddress)
			pubkeyStr:=viper.GetString(FlagSignerPubkey)
			startEpoch := viper.GetInt64(FlagStartEpoch)
			endEpoch:= viper.GetInt64(FlagEndEpoch)
			amountStr := viper.GetString(FlagAmount)

			validatorAddr := common.HexToAddress(validatorStr)

			pubkeyBytes,err:=hex.DecodeString(pubkeyStr)
			if err!=nil{
				return err
			}
			pubkey:=types.NewPubKey(pubkeyBytes)

			msg := staking.NewMsgValidatorJoin(validatorAddr,pubkey,uint64(startEpoch),uint64(endEpoch),json.Number(amountStr))

			return helper.CreateAndSendTx(msg,cliCtx)
		},
	}
	return cmd
}

// send validator exit transaction
func GetValidatorExitTx(cdc *codec.Codec) *cobra.Command  {
	cmd:=&cobra.Command{
		Use:   "validator-exit",
		Short: "Exit heimdall as a valdiator ",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			validatorStr:= viper.GetString(FlagValidatorAddress)
			validatorAddr := common.HexToAddress(validatorStr)

			msg := staking.NewMsgValidatorExit(validatorAddr)

			return helper.CreateAndSendTx(msg,cliCtx)
		},
	}
	return cmd
}

// send validator update transaction
func GetValidatorUpdateTx(cdc *codec.Codec) *cobra.Command  {
	cmd:=&cobra.Command{
		Use:   "signer-update",
		Short: "Update signer for a validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			validatorStr:= viper.GetString(FlagValidatorAddress)
			pubkeyStr:=viper.GetString(FlagSignerPubkey)
			amountStr := viper.GetString(FlagAmount)

			validatorAddr := common.HexToAddress(validatorStr)

			pubkeyBytes,err:=hex.DecodeString(pubkeyStr)
			if err!=nil{
				return err
			}
			pubkey:=types.NewPubKey(pubkeyBytes)

			msg := staking.NewMsgValidatorUpdate(validatorAddr,pubkey,json.Number(amountStr))

			return helper.CreateAndSendTx(msg,cliCtx)
		},
	}
	return cmd
}