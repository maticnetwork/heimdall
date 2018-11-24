package helper

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
	hmtypes "github.com/maticnetwork/heimdall/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/maticnetwork/heimdall/checkpoint"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

type validatorPretty struct {
	Address cmn.HexBytes `json:"address"`
	Power   int64        `json:"power"`
}

func ValidatorsToString(vs []abci.Validator) string {
	s := make([]validatorPretty, len(vs))
	for i, v := range vs {
		s[i] = validatorPretty{
			Address: v.Address,
			Power:   v.Power,
		}
	}
	b, err := json.Marshal(s)
	if err != nil {
		panic(err.Error())
	}
	return string(b)
}

func UpdateValidators(currentSet *types.ValidatorSet, abciUpdates []abci.ValidatorUpdate) error {
	updates, err := types.PB2TM.ValidatorUpdates(abciUpdates)
	if err != nil {
		return err
	}

	// these are tendermint types now
	for _, valUpdate := range updates {
		if valUpdate.VotingPower < 0 {
			return fmt.Errorf("Voting power can't be negative %v", valUpdate)
		}

		address := valUpdate.Address
		_, val := currentSet.GetByAddress(address)
		if valUpdate.VotingPower == 0 {
			// remove val
			_, removed := currentSet.Remove(address)
			if !removed {
				return fmt.Errorf("Failed to remove validator %X", address)
			}
		} else if val == nil {
			// add val
			added := currentSet.Add(valUpdate)
			if !added {
				return fmt.Errorf("Failed to add new validator %v", valUpdate)
			}
		} else {
			// update val
			updated := currentSet.Update(valUpdate)
			if !updated {
				return fmt.Errorf("Failed to update validator %X to %v", address, valUpdate)
			}
		}
	}
	return nil
}

// convert string to Pubkey
func StringToPubkey(pubkeyStr string) (crypto.PubKey, error) {

	var pubkeyBytes secp256k1.PubKeySecp256k1
	_pubkey, err := hex.DecodeString(pubkeyStr)
	if err != nil {
		Logger.Error("Decoding of pubkey(string) to pubkey failed", "Error", err, "PubkeyString", pubkeyStr)
		return nil, err
	}
	// copy
	copy(pubkeyBytes[:], _pubkey)

	return pubkeyBytes, nil
}

func CreateTxBytes(msg sdk.Msg) ([]byte, error) {
	tx := hmtypes.NewBaseTx(msg)

	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		Logger.Error("Error generating TX Bytes", "error", err)

		return []byte(""), err
	}
	return txBytes, nil
}

func SendTendermintRequest(cliCtx context.CLIContext, txBytes []byte) (*ctypes.ResultBroadcastTxCommit, error) {
	Logger.Info("Broadcasting tx bytes to Tendermint", "txBytes", hex.EncodeToString(txBytes))
	return cliCtx.BroadcastTx(txBytes)
}

func getSigs(votes []tmtypes.Vote) (sigs []byte) {

	// loop votes and append to sig to sigs
	for _, vote := range votes {
		sigs = append(sigs[:], vote.Signature[:]...)
	}
	return
}

func getVoteBytes(votes []tmtypes.Vote, ctx sdk.Context) []byte {
	// sign bytes for vote
	return votes[0].SignBytes(ctx.ChainID())
}

func getExtraData(_checkpoint checkpoint.CheckpointBlockHeader, ctx sdk.Context) []byte {
	Logger.Debug("Creating extra data", "startBlock", _checkpoint.StartBlock, "endBlock", _checkpoint.EndBlock, "roothash", _checkpoint.RootHash)

	// craft a message
	msg := checkpoint.NewMsgCheckpointBlock(_checkpoint.Proposer, _checkpoint.StartBlock, _checkpoint.EndBlock, _checkpoint.RootHash)

	// decoding transaction
	tx := hmtypes.NewBaseTx(msg)
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		Logger.Error("Error decoding transaction data", "error", err)
	}

	return txBytes
}

func PrepareAndSendCheckpoint(ctx sdk.Context, checkpointKeeper checkpoint.Keeper) {
	// fetch votes from block header
	var votes []tmtypes.Vote
	err := json.Unmarshal(ctx.BlockHeader().Votes, &votes)
	if err != nil {
		Logger.Error("Error while unmarshalling vote", "error", err)
	}

	// get sigs from votes
	sigs := getSigs(votes)

	// Getting latest checkpoint data from store using height as key and unmarshall
	_checkpoint, err := checkpointKeeper.GetCheckpointFromBuffer(ctx)
	if err != nil {
		Logger.Error("Unable to unmarshall checkpoint while fetching from buffer while preparing checkpoint tx for rootchain", "error", err, "height", ctx.BlockHeight())
		panic(err)
	} else {
		// Get extra data
		extraData := getExtraData(_checkpoint, ctx)

		//fetch current child block from rootchain contract
		lastblock, err := CurrentChildBlock()
		if err != nil {
			Logger.Error("Could not fetch last block from mainchain", "Error", err)
			panic(err)
		}

		if lastblock == _checkpoint.StartBlock {
			Logger.Info("Sending Valid Checkpoint ...")
			//SendCheckpoint(getVoteBytes(votes, ctx), sigs, extraData)
		} else {
			Logger.Error("Start block does not match", "lastBlock", lastblock, "startBlock", _checkpoint.StartBlock)
			// TODO panic ?
		}
	}
}
