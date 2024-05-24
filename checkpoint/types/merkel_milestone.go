package types

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// ValidateMilestone validates the structure of the milestone
func ValidateMilestone(start uint64, end uint64, rootHash hmTypes.HeimdallHash, milestoneID string, contractCaller helper.IContractCaller, milestoneLength uint64, confirmations uint64) (bool, error) {
	msgMilestoneLength := int64(end) - int64(start) + 1

	// Check for the minimum length of the milestone
	if msgMilestoneLength < int64(milestoneLength) {
		return false, errors.New(fmt.Sprint("Invalid milestone, difference in start and end block is less than milestone length", "Milestone Length:", milestoneLength))
	}

	// Check if blocks+confirmations  exist locally
	if !contractCaller.CheckIfBlocksExist(end + confirmations) {
		return false, errors.New(fmt.Sprint("End block number with confirmation is not available in the Bor chain", "EndBlock", end, "Confirmation", confirmations))
	}

	// Get the vote on hash of milestone from Bor
	vote, err := contractCaller.GetVoteOnHash(start, end, milestoneLength, rootHash.String(), milestoneID)
	if err != nil {
		return false, err
	}

	// validate that milestoneID is composed by `UUID - HexAddressOfTheProposer`
	splitMilestoneID := strings.Split(strings.TrimSpace(milestoneID), " - ")
	if len(splitMilestoneID) != 2 {
		return false, errors.New(fmt.Sprint("invalid milestoneID, it should be composed by `UUID - HexAddressOfTheProposer`", "milestoneID", milestoneID))
	}

	_, err = uuid.Parse(splitMilestoneID[0])
	if err != nil {
		return false, errors.New(fmt.Sprint("invalid milestoneID, the UUID is not correct", "milestoneID", milestoneID))
	}

	if !common.IsHexAddress(splitMilestoneID[1]) {
		return false, errors.New(fmt.Sprint("invalid milestoneID, the proposer address is not correct", "milestoneID", milestoneID))
	}

	return vote, nil
}
