package types

import (
	"bytes"
	"errors"

	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// ValidateMilestone - Validates if milestone rootHash matches or not
func ValidateMilestone(start uint64, end uint64, rootHash hmTypes.HeimdallHash, contractCaller helper.IContractCaller, sprintLength uint64) (bool, error) {

	if start+sprintLength-1 != end {
		return false, errors.New("Invalid milestone, difference in start and end block is not equal to sprint length")
	}

	// Check if blocks exist locally
	if !contractCaller.CheckIfBlocksExist(end) {
		return false, errors.New("blocks not found locally")
	}

	// Compare RootHash
	root, err := contractCaller.GetRootHash(start, end, sprintLength)
	if err != nil {
		return false, err
	}

	if bytes.Equal(root, rootHash.Bytes()) {
		return true, nil
	}

	return false, nil
}

func convertTo32(input []byte) (output [32]byte, err error) {
	l := len(input)
	if l > 32 || l == 0 {
		return
	}

	copy(output[32-l:], input[:])

	return
}

func appendBytes32(data ...[]byte) []byte {
	var result []byte

	for _, v := range data {
		paddedV, err := convertTo32(v)
		if err == nil {
			result = append(result, paddedV[:]...)
		}
	}

	return result
}
