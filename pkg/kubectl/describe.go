package kubectl

import "github.com/spf13/cobra"

var describeCmd = &cobra.Command{
	Use:   "describe <resource> | (<resource> <resource-name>)",
	Short: "descirbe resources",
	Long:  "descirbe resources",
	Args:  cobra.RangeArgs(1, 2),
	Run:   describe,
}
