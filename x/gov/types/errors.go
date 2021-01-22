package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/gov module sentinel errors
var (
	ErrUnknownProposal         = sdkerrors.Register(ModuleName, 10002, "unknown proposal")
	ErrInactiveProposal        = sdkerrors.Register(ModuleName, 10003, "inactive proposal")
	ErrAlreadyActiveProposal   = sdkerrors.Register(ModuleName, 10004, "proposal already active")
	ErrInvalidProposalContent  = sdkerrors.Register(ModuleName, 10005, "invalid proposal content")
	ErrInvalidProposalType     = sdkerrors.Register(ModuleName, 10006, "invalid proposal type")
	ErrInvalidVote             = sdkerrors.Register(ModuleName, 10007, "invalid vote option")
	ErrInvalidGenesis          = sdkerrors.Register(ModuleName, 10008, "invalid genesis state")
	ErrNoProposalHandlerExists = sdkerrors.Register(ModuleName, 10009, "no handler exists for proposal type")
	ErrAlreadyFinishedProposal = sdkerrors.Register(ModuleName, 10010, "proposal has already passed its voting period")
)
