package kubectl

import "github.com/spf13/cobra"

var describeCmd = &cobra.Command{
	Use:   "describe <resource>",
	Short: "descirbe resources",
	Long:  "descirbe resources",
	Args:  cobra.MinimumNArgs(1),
	Run:   describe,
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
