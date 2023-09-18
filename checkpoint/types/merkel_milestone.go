package types

import (
	"errors"
	"fmt"

	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// ValidateMilestone - Validates if milestone rootHash matches or not
func ValidateMilestone(start uint64, end uint64, rootHash hmTypes.HeimdallHash, milestoneID string, contractCaller helper.IContractCaller, milestoneLength uint64, confirmations uint64) (bool, error) {
	msgMilestoneLength := int64(end) - int64(start) + 1

	//Check for the minimum length of the milestone
	if msgMilestoneLength < int64(milestoneLength) {
		return false, errors.New(fmt.Sprint("Invalid milestone, difference in start and end block is less than milestone length", "Milestone Length:", milestoneLength))
	}

	fmt.Println("#########Entered Here 2######")
	// Check if blocks+confirmations  exist locally
	if !contractCaller.CheckIfBlocksExist(end + confirmations) {
		return false, errors.New(fmt.Sprint("End block number with confirmation is not availbale in the Bor chain", "EndBlock", end, "Confirmation", confirmations))
	}

	fmt.Println("#########Entered Here 4 4######")
	//Get the vote on hash of milestone from Bor
	vote, err := contractCaller.GetVoteOnHash(start, end, milestoneLength, rootHash.String(), milestoneID)

	fmt.Println("#########Entered Here End######")
	if err != nil {
		return false, err
	}

	return vote, nil
}
