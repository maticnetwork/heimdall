package cli

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	amino "github.com/tendermint/go-amino"

	"github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/helper"
)

var logger = helper.Logger.With("module", "auth/client/cli")

// GetSignCommand returns the transaction sign command.
func GetSignCommand(codec *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sign [file]",
		Short: "Sign transactions generated offline",
		Long: `Sign transactions created with the --generate-only flag.
It will read a transaction from [file], sign it, and print its JSON encoding.

If the flag --signature-only flag is set, it will output a JSON representation
of the generated signature only.

If the flag --validate-signatures is set, then the command would check whether all required
signers have signed the transactions, whether the signatures were collected in the right
order, and if the signature is valid over the given transaction. If the --offline
flag is also set, signature validation over the transaction will be not be
performed as that will require RPC communication with a full node.

The --offline flag makes sure that the client will not reach out to full node.
As a result, the account and sequence number queries will not be performed and
it is required to set such parameters manually. Note, invalid values will cause
the transaction to fail.
`,
		PreRun: preSignCmd,
		RunE:   makeSignCmd(codec),
		Args:   cobra.ExactArgs(1),
	}

	cmd.Flags().Bool(flagSigOnly, false, "Print only the generated signature, then exit")
	cmd.Flags().Bool(flagOffline, false, "Offline mode; Do not query a full node")
	cmd.Flags().String(flagOutfile, "", "The document will be written to the given file instead of STDOUT")

	cmd = client.PostCommands(cmd)[0]
	// cmd.MarkFlagRequired(client.FlagFrom)

	return cmd
}

func preSignCmd(cmd *cobra.Command, _ []string) {
	// Conditionally mark the account and sequence numbers required as no RPC
	// query will be done.
	if viper.GetBool(flagOffline) {
		if err := cmd.MarkFlagRequired(client.FlagAccountNumber); err != nil {
			logger.Error("preSignCmd | MarkFlagRequired | FlagAccountNumber", "Error", err)
		}
		if err := cmd.MarkFlagRequired(client.FlagSequence); err != nil {
			logger.Error("preSignCmd | MarkFlagRequired | FlagSequence", "Error", err)
		}
	}
}

func makeSignCmd(cdc *amino.Codec) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) (err error) {
		cliCtx := context.NewCLIContext().WithCodec(cdc)
		stdTx, err := helper.ReadStdTxFromFile(cliCtx.Codec, args[0])
		if err != nil {
			return err
		}

		offline := viper.GetBool(flagOffline)

		// if --signature-only is on, then override --append
		var newTx types.StdTx
		generateSignatureOnly := viper.GetBool(flagSigOnly)

		appendSig := viper.GetBool(flagAppend) && !generateSignatureOnly
		newTx, err = helper.SignStdTx(cliCtx, stdTx, appendSig, offline)

		if err != nil {
			return err
		}

		var json []byte

		switch generateSignatureOnly {
		case true:
			switch cliCtx.Indent {
			case true:
				json, err = cdc.MarshalJSONIndent(newTx.Signature, "", "  ")

			default:
				json, err = cdc.MarshalJSON(newTx.Signature)
			}

		default:
			switch cliCtx.Indent {
			case true:
				json, err = cdc.MarshalJSONIndent(newTx, "", "  ")

			default:
				json, err = cdc.MarshalJSON(newTx)
			}
		}

		if err != nil {
			return err
		}

		if viper.GetString(flagOutfile) == "" {
			fmt.Printf("%s\n", json)
			return
		}

		fp, err := os.OpenFile(
			viper.GetString(flagOutfile), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644,
		)
		if err != nil {
			return err
		}

		defer fp.Close()
		fmt.Fprintf(fp, "%s\n", json)

		return
	}
}
