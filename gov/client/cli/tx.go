package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	govutils "github.com/maticnetwork/heimdall/gov/client/utils"
	"github.com/maticnetwork/heimdall/gov/types"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var logger = helper.Logger.With("module", "gov/client/cli")

// Proposal flags
const (
	FlagTitle        = "title"
	FlagDescription  = "description"
	flagProposalType = "type"
	FlagDeposit      = "deposit"
	flagVoter        = "voter"
	flagDepositor    = "depositor"
	flagStatus       = "status"
	flagNumLimit     = "limit"
	FlagProposal     = "proposal"
	FlagValidatorID  = "validator-id"
)

type proposal struct {
	Title       string
	Description string
	Type        string
	Deposit     string
}

// ProposalFlags defines the core required fields of a proposal. It is used to
// verify that these values are not provided in conjunction with a JSON proposal
// file.
var ProposalFlags = []string{
	FlagTitle,
	FlagDescription,
	flagProposalType,
	FlagDeposit,
}

// GetTxCmd returns the transaction commands for this module
// governance ModuleClient is slightly different from other ModuleClients in that
// it contains a slice of "proposal" child commands. These commands are respective
// to proposal type handlers that are implemented in other modules but are mounted
// under the governance CLI (eg. parameter change proposals).
func GetTxCmd(storeKey string, cdc *codec.Codec, pcmds []*cobra.Command) *cobra.Command {
	govTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Governance transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmdSubmitProp := GetCmdSubmitProposal(cdc)
	for _, pcmd := range pcmds {
		cmdSubmitProp.AddCommand(client.PostCommands(pcmd)[0])
	}

	govTxCmd.AddCommand(client.PostCommands(
		GetCmdDeposit(cdc),
		GetCmdVote(cdc),
		cmdSubmitProp,
	)...)

	return govTxCmd
}

// GetCmdSubmitProposal implements submitting a proposal transaction command.
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-proposal",
		Short: "Submit a proposal along with an initial deposit",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a proposal along with an initial deposit.
Proposal title, description, type and deposit can be given directly or through a proposal JSON file.

Example:
$ %s tx gov submit-proposal --proposal="path/to/proposal.json" --from mykey

Where proposal.json contains:

{
  "title": "Test Proposal",
  "description": "My awesome proposal",
  "type": "Text",
  "deposit": "10test"
}

Which is equivalent to:

$ %s tx gov submit-proposal --title="Test Proposal" --description="My awesome proposal" --type="Text" --deposit="10test" --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := parseSubmitProposalFlags()
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoins(proposal.Deposit)
			if err != nil {
				return err
			}

			validatorID := viper.GetInt64(FlagValidatorID)
			if validatorID == 0 {
				return fmt.Errorf("Valid validator ID required")
			}

			content := types.ContentFromProposalType(proposal.Title, proposal.Description, proposal.Type)
			from := helper.GetFromAddress(cliCtx)

			msg := types.NewMsgSubmitProposal(content, amount, from, hmTypes.ValidatorID(validatorID))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(FlagTitle, "", "Title of proposal")
	cmd.Flags().String(FlagDescription, "", "Description of proposal")
	cmd.Flags().String(flagProposalType, "", "Type of proposal, types: text/parameter_change/software_upgrade")
	cmd.Flags().String(FlagDeposit, "", "Deposit of proposal")
	cmd.Flags().String(FlagProposal, "", "Proposal file path (if this path is given, other proposal flags are ignored)")
	cmd.Flags().Int(FlagValidatorID, 0, "--validator-id=<validator ID here>")
	if err := cmd.MarkFlagRequired(FlagValidatorID); err != nil {
		logger.Error("GetCmdSubmitProposal | MarkFlagRequired | FlagValidatorID", "Error", err)
	}

	return cmd
}

// GetCmdDeposit implements depositing tokens for an active proposal.
func GetCmdDeposit(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [proposal-id] [deposit]",
		Args:  cobra.ExactArgs(2),
		Short: "Deposit tokens for an active proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a deposit for an active proposal. You can
find the proposal-id by running "%s query gov proposals".

Example:
$ %s tx gov deposit 1 10stake --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid uint, please input a valid proposal-id", args[0])
			}

			// Get depositor address
			from := helper.GetFromAddress(cliCtx)

			// Get amount of coins
			amount, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			validatorID := viper.GetInt64(FlagValidatorID)
			if validatorID == 0 {
				return fmt.Errorf("Valid validator ID required")
			}

			msg := types.NewMsgDeposit(from, proposalID, amount, hmTypes.ValidatorID(validatorID))
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Int(FlagValidatorID, 0, "--validator-id=<validator ID here>")
	if err := cmd.MarkFlagRequired(FlagValidatorID); err != nil {
		logger.Error("GetCmdDeposit | MarkFlagRequired | FlagValidatorID", "Error", err)
	}

	return cmd
}

// GetCmdVote implements creating a new vote command.
func GetCmdVote(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [proposal-id] [option]",
		Args:  cobra.ExactArgs(2),
		Short: "Vote for an active proposal, options: yes/no/no_with_veto/abstain",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a vote for an active proposal. You can
find the proposal-id by running "%s query gov proposals".


Example:
$ %s tx gov vote 1 yes --from mykey
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// Get voting address
			from := helper.GetFromAddress(cliCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			validatorID := viper.GetInt64(FlagValidatorID)
			if validatorID == 0 {
				return fmt.Errorf("Valid validator ID required")
			}

			// Find out which vote option user chose
			byteVoteOption, err := types.VoteOptionFromString(govutils.NormalizeVoteOption(args[1]))
			if err != nil {
				return err
			}

			// Build vote message and run basic validation
			msg := types.NewMsgVote(from, proposalID, byteVoteOption, hmTypes.ValidatorID(validatorID))
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Int(FlagValidatorID, 0, "--validator-id=<validator ID here>")
	if err := cmd.MarkFlagRequired(FlagValidatorID); err != nil {
		logger.Error("GetCmdVote | MarkFlagRequired | FlagValidatorID", "Error", err)
	}

	return cmd
}

// DONTCOVER
