package cli

import (
	"fmt"
	"strings"

	"github.com/maticnetwork/heimdall/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"

	govTypes "github.com/maticnetwork/heimdall/gov/types"
	paramscutils "github.com/maticnetwork/heimdall/params/client/utils"
	"github.com/maticnetwork/heimdall/params/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var logger = helper.Logger.With("module", "params/client/cli")

// GetCmdSubmitProposal implements a command handler for submitting a parameter
// change proposal transaction.
func GetCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
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
  "deposit": [
    {
      "denom": "matic",
      "amount": "1000000000000000000" 
    }
  ]
}
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			proposal, err := paramscutils.ParseParamChangeProposalJSON(cdc, args[0])
			if err != nil {
				return err
			}

			validatorID := viper.GetUint64(FlagValidatorID)
			if validatorID == 0 {
				return fmt.Errorf("Valid validator ID required")
			}

			from := helper.GetFromAddress(cliCtx)
			content := types.NewParameterChangeProposal(proposal.Title, proposal.Description, proposal.Changes.ToParamChanges())

			// create submit proposal
			msg := govTypes.NewMsgSubmitProposal(content, proposal.Deposit, from, hmTypes.NewValidatorID(validatorID))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Int(FlagValidatorID, 0, "--validator-id=<validator ID here>")
	if err := cmd.MarkFlagRequired(FlagValidatorID); err != nil {
		logger.Error("GetCmdSubmitProposal | MarkFlagRequired | FlagValidatorID", "Error", err)
	}

	return cmd
}
