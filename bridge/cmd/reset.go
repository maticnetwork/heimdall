package cmd

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// resetCmd represents the start command
var resetCmd = &cobra.Command{
	Use:   "unsafe-reset-all",
	Short: "Reset bridge server data",
	Run: func(cmd *cobra.Command, args []string) {
		dbLocation := viper.GetString(bridgeDBFlag)
		if dir, err := ioutil.ReadDir(dbLocation); err != nil {
			// fmt.Println(err)
		} else {
			for _, d := range dir {
				os.RemoveAll(path.Join([]string{dbLocation, d.Name()}...))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
