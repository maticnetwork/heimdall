package types

import (
	"github.com/maticnetwork/bor/rlp"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

// SortAndRLPEncodeSlashInfos  - RLP encoded slashing infos
func SortAndRLPEncodeSlashInfos(slashingInfos []*hmTypes.ValidatorSlashingInfo) ([]byte, error) {

	// Sort the slashingInfos by ID
	slashingInfos = hmTypes.SortValidatorSlashingInfoByID(slashingInfos)

	// Encode slashInfos
	encodedSlashInfos, err := rlp.EncodeToBytes(slashingInfos)

	return encodedSlashInfos, err
}

func RLPDecodeSlashInfos(encodedSlashInfo []byte) ([]*hmTypes.ValidatorSlashingInfo, error) {
	var slashingInfoList []*hmTypes.ValidatorSlashingInfo
	err := rlp.DecodeBytes(encodedSlashInfo, &slashingInfoList)
	return slashingInfoList, err

}

func RLPDeocdeTickVoteBytes(tickMsgVoteBytes []byte) (tmTypes.CanonicalRLPVote, error) {

	var vote tmTypes.CanonicalRLPVote
	err := rlp.DecodeBytes(tickMsgVoteBytes, &vote)
	return vote, err
}

func RLPDeocdeStdTxBytes(stdTxBytes []byte) (authTypes.StdTx, error) {
	var stdTx authTypes.StdTx
	err := rlp.DecodeBytes(stdTxBytes, &stdTx)
	return stdTx, err
}
