package kubectl

import (
	"MiniK8S/config"
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show kubectl version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("kubectl version:", config.Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
