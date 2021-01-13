package cli

// Proposal flags
const (
	FlagTitle        = "title"
	FlagDescription  = "description"
	FlagProposalType = "type"
	FlagDeposit      = "deposit"
	FlagVoter        = "voter"
	FlagDepositor    = "depositor"
	FlagStatus       = "status"
	FlagNumLimit     = "limit"
	FlagProposal     = "proposal"
	FlagValidatorID  = "validator-id"
)

// ProposalFlags defines the core required fields of a proposal. It is used to
// verify that these values are not provided in conjunction with a JSON proposal
// file.
var ProposalFlags = []string{
	FlagTitle,
	FlagDescription,
	FlagProposalType,
	FlagDeposit,
}
