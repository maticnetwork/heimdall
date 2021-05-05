package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/gov/types"
	paramscutils "github.com/maticnetwork/heimdall/x/params/client/utils"
	paramproposal "github.com/maticnetwork/heimdall/x/params/types/proposal"
)

// NewSubmitParamChangeProposalTxCmd returns a CLI command handler for creating
// a parameter change proposal governance transaction.
func NewSubmitParamChangeProposalTxCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "param-change [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a parameter change proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a parameter proposal along with an initial deposit.
The proposal details must be supplied via a JSON file. For values that contains
objects, only non-empty fields will be updated.

IMPORTANT: Currently parameter changes are evaluated but not validated, so it is
very important that any "value" change is valid (ie. correct type and within bounds)
for its respective parameter, eg. "MaxValidators" should be an integer and not a decimal.

Proper vetting of a parameter change proposal should prevent this from happening
(no deposits should occur during the governance process), but it should be noted
regardless.

Example:
$ %s tx gov submit-proposal param-change <path/to/proposal.json> --from=<key_or_address>

Where proposal.json contains:

{
  "title": "Staking Param Change",
  "description": "Update max validators",
  "changes": [
    {
      "subspace": "staking",
      "key": "MaxValidators",
      "value": 105
    }
  ],
  "deposit": "1000stake"
}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			proposal, err := paramscutils.ParseParamChangeProposalJSON(clientCtx.LegacyAmino, args[0])
			if err != nil {
				return err
			}

			content := paramproposal.NewParameterChangeProposal(
				proposal.Title, proposal.Description, proposal.Changes.ToParamChanges(),
			)

			amount, err := sdk.ParseCoinsNormalized(proposal.Deposit)
			if err != nil {
				return err
			}

			validatorID, err := cmd.Flags().GetInt(FlagValidatorID)
			if err != nil {
				return err
			}
			if validatorID == 0 {
				return fmt.Errorf("Valid validator ID required")
			}

			msg, err := types.NewMsgSubmitProposal(content, amount, helper.GetFromAddress(clientCtx), hmTypes.ValidatorID(validatorID))
			if err != nil {
				return fmt.Errorf("invalid message: %w", err)
			}

			if err = msg.ValidateBasic(); err != nil {
				return fmt.Errorf("message validation failed: %w", err)
			}

			return helper.GenerateOrBroadcastTxCli(clientCtx, cmd.Flags(), msg)
		},
	}
}
