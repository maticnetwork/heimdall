package version

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"gopkg.in/yaml.v3"
)

const flagLong = "long"

func init() {
	Cmd.Flags().Bool(flagLong, false, "Print long version information")
}

// Cmd prints out the application's version information passed via build flags.
var Cmd = &cobra.Command{
	Use:   "version",
	Short: "Print the app version",
	RunE: func(_ *cobra.Command, _ []string) error {
		verInfo := NewInfo()

		if !viper.GetBool(flagLong) {
			fmt.Println(verInfo.Version)
			return nil
		}

		var bz []byte
		var err error

		switch viper.GetString(cli.OutputFlag) {
		case "json":
			bz, err = jsoniter.ConfigFastest.Marshal(verInfo)
		default:
			bz, err = yaml.Marshal(&verInfo)
		}

		if err != nil {
			return err
		}

		_, err = fmt.Println(string(bz))
		return err
	},
}
