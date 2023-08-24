package cmd

import (
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// resetCmd represents the start command
var resetCmd = &cobra.Command{
	Use:   "unsafe-reset-all",
	Short: "Reset bridge server data",
	RunE: func(cmd *cobra.Command, args []string) error {
		dbLocation := viper.GetString(bridgeDBFlag)
		dir, err := os.ReadDir(dbLocation)
		if err != nil {
			return err
		}

		for _, d := range dir {
			os.RemoveAll(path.Join([]string{dbLocation, d.Name()}...))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
